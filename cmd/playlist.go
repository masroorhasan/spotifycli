package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zmb3/spotify"
)

var (
	addtoPlaylistName string
	newPlaylistName   string
)

func newAddtoPlaylistCmd() *cobra.Command {
	addtoCmd := &cobra.Command{
		Use:   "addto --p [PLAYLIST]",
		Short: "Add currently playing track to playlist",
		RunE: func(cmd *cobra.Command, args []string) error {
			return addto(cmd, args)
		},
	}
	addtoCmd.Flags().StringVar(&addtoPlaylistName, "p", "", "Add current track to specified playlist.")
	return addtoCmd
}

func newCreatePlaylistCmd() *cobra.Command {
	newCmd := &cobra.Command{
		Use:   "new --p [PLAYLIST]",
		Short: "Create new playlist",
		RunE: func(cmd *cobra.Command, args []string) error {
			return newPlaylist(cmd, args)
		},
	}
	newCmd.Flags().StringVar(&newPlaylistName, "p", "", "Name of new playlist.")
	return newCmd
}

func addto(cmd *cobra.Command, args []string) error {
	// current user
	user, err := client.CurrentUser()
	if err != nil {
		return err
	}
	fmt.Println("Current user: ", user.DisplayName)

	// get current playing song
	playing, err := client.PlayerCurrentlyPlaying()
	if err != nil {
		return err
	}
	fmt.Println("Currently playing: ", playing.Item.Name)

	// get my playlists
	playlists, err := client.CurrentUsersPlaylists()
	if err != nil {
		return err
	}
	var matchPlaylist spotify.SimplePlaylist
	for _, p := range playlists.Playlists {
		if addtoPlaylistName == p.Name {
			matchPlaylist = p
			break
		}
	}
	fmt.Println("Matched playlist: ", matchPlaylist.Name)

	// add track to playlist
	snapshotID, err := client.AddTracksToPlaylist(user.ID, matchPlaylist.ID, playing.Item.ID)
	if err != nil {
		return err
	}
	fmt.Println("Added track to playlist, snapshotID: ", snapshotID)
	return nil
}

func newPlaylist(cmd *cobra.Command, args []string) error {
	// current user
	user, err := client.CurrentUser()
	if err != nil {
		return err
	}
	fmt.Println("Current user: ", user.DisplayName)

	// create new playlist
	playlist, err := client.CreatePlaylistForUser(user.ID, newPlaylistName, true)
	if err != nil {
		return err
	}
	fmt.Println("Created public playlist: ", playlist.Name)
	return nil
}
