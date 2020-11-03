package cmd

import (
	"fmt"
	"log"
	"net/http"

	"github.com/agext/uuid"
	"github.com/spf13/cobra"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

var (
	ch = make(chan *oauth2.Token)
)

type authenticationHandler struct {
	auth  spotify.Authenticator
	state string
}

func newLoginCmd() *cobra.Command {
	loginCmd := &cobra.Command{
		Use:   "login",
		Short: "Login to authenticate Spotify account",
		RunE: func(cmd *cobra.Command, args []string) error {
			return authorize(cmd, args)
		},
	}
	return loginCmd
}

func newLogoutCmd() *cobra.Command {
	logoutCmd := &cobra.Command{
		Use:   "logout",
		Short: "Logout from Spotify account",
		RunE: func(cmd *cobra.Command, args []string) error {
			return deleteToken()
		},
	}
	return logoutCmd
}

func authorize(cmd *cobra.Command, args []string) error {
	// use uuid as state
	state := string(uuid.New().Hex())

	// setup server for callback
	http.Handle("/callback", &authenticationHandler{auth: auth, state: state})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for: ", r.URL.String())
	})
	go http.ListenAndServe(":8080", nil)

	// User authentication process
	fmt.Println("authorize")
	url := auth.AuthURL(state)
	fmt.Println("Please log in to Spotify by visiting the following page in your browser:", url)

	// persist token
	token := <-ch
	return persistToken(token)
}

func (handler *authenticationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	token, err := handler.auth.Token(handler.state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != handler.state {
		http.NotFound(w, r)
		log.Fatalf("state miss match: %s != %s\n", st, handler.state)
	}
	ch <- token
}
