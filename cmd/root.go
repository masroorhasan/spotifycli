package cmd

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

const (
	tokenFile   = ".sptok"
	redirectURI = "http://localhost:8080/callback"
)

var (
	auth   spotify.Authenticator
	client spotify.Client
)

// NewRootCmd gets the root cmd.
func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:               "spotifycli",
		Short:             "A command line interface to manage Spotify playlists.",
		PersistentPreRun:  prerun,
		PersistentPostRun: postrun,
	}
	// auth ops
	rootCmd.AddCommand(newLoginCmd())
	rootCmd.AddCommand(newLogoutCmd())

	// search ops
	rootCmd.AddCommand(newSearchCmd())

	// playlist ops
	rootCmd.AddCommand(newCurrentTrackCmd())
	rootCmd.AddCommand(newCreatePlaylistCmd())
	rootCmd.AddCommand(newDeletePlaylistCmd())
	rootCmd.AddCommand(newAddtoPlaylistCmd())
	rootCmd.AddCommand(newAddTrackByIDToPlaylistCmd())
	rootCmd.AddCommand(newAddTrackByNameToPlaylistCmd())
	rootCmd.AddCommand(newRemoveTrackFromPlaylistCmd())
	rootCmd.AddCommand(newListPlaylistTracksCmd())
	return rootCmd
}

func prerun(cmd *cobra.Command, args []string) {
	// initialize authenticator
	auth = spotify.NewAuthenticator(
		redirectURI,
		spotify.ScopeUserReadPrivate,
		spotify.ScopeUserReadCurrentlyPlaying,
		spotify.ScopePlaylistModifyPrivate,
		spotify.ScopePlaylistModifyPublic)
	auth.SetAuthInfo(os.Getenv("SPOTIFY_ID"), os.Getenv("SPOTIFY_SECRET"))

	// exit early
	if cmd.Use == "login" || cmd.Use == "logout" {
		return
	}

	// get token
	token, err := getToken()
	if err != nil {
		if err := authorize(cmd, args); err != nil {
			log.Fatal(err)
		}
	}
	client = auth.NewClient(token)
}

func postrun(cmd *cobra.Command, args []string) {
	// exit early
	if cmd.Use == "login" || cmd.Use == "logout" {
		return
	}

	// refresh token
	currTok, err := client.Token()
	if err != nil {
		log.Fatal(err)
	}
	token, err := getToken()
	if err != nil {
		log.Fatal(err)
	}
	if token != currTok {
		if err := persistToken(token); err != nil {
			log.Fatal(err)
		}
	}
}

func persistToken(token *oauth2.Token) error {
	u, err := user.Current()
	if err != nil {
		return err
	}

	// detect if file exists
	tokPath := filepath.Join(u.HomeDir, tokenFile)
	_, err = os.Stat(tokPath)

	// create and write
	if os.IsNotExist(err) {
		file, err := os.Create(tokPath)
		if err != nil {
			return err
		}
		defer file.Close()

		b, err := json.Marshal(token)
		if err != nil {
			return err
		}

		// write
		_, err = file.Write(b)
		if err != nil {
			return err
		}
		return file.Sync()
	}
	return nil
}

func getToken() (*oauth2.Token, error) {
	u, err := user.Current()
	if err != nil {
		return nil, err
	}

	// detect if file exists
	tokPath := filepath.Join(u.HomeDir, tokenFile)
	_, err = os.Stat(tokPath)

	// exit early
	if os.IsNotExist(err) {
		return nil, err
	}

	// read data
	data, err := ioutil.ReadFile(tokPath)
	if err != nil {
		return nil, err
	}

	// unmarshal data
	var token oauth2.Token
	if err := json.Unmarshal(data, &token); err != nil {
		return nil, err
	}
	return &token, nil
}

func deleteToken() error {
	u, err := user.Current()
	if err != nil {
		return err
	}

	// detect if file exists
	tokPath := filepath.Join(u.HomeDir, tokenFile)
	_, err = os.Stat(tokPath)
	if os.IsNotExist(err) {
		return err
	}
	return os.Remove(tokPath)
}
