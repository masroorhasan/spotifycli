package cmd

import (
	"fmt"
	"reflect"

	"github.com/spf13/cobra"
	"github.com/zmb3/spotify"
)

var (
	addtoPlaylistName string
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

var (
	addTrackID             string
	addTrackToPlaylistName string
)

func newAddTrackByIDToPlaylistCmd() *cobra.Command {
	addCmd := &cobra.Command{
		Use:   "add --tid [TRACK_ID] --p [PLAYLIST]",
		Short: "Add track to playlist",
		RunE: func(cmd *cobra.Command, args []string) error {
			return addTrackByIDToPlaylist(cmd, args)
		},
	}
	addCmd.Flags().StringVar(&addTrackID, "tid", "", "Id of track to add to playlist.")
	addCmd.Flags().StringVar(&addTrackToPlaylistName, "p", "", "Name of playlist to add track to.")
	return addCmd
}

var (
	newPlaylistName string
)

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

var (
	delPlaylistName string
)

func newDeletePlaylistCmd() *cobra.Command {
	deleteCmd := &cobra.Command{
		Use:   "del --p [PLAYLIST]",
		Short: "Delete a playlist",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	deleteCmd.Flags().StringVar(&delPlaylistName, "p", "", "Name of playlist to delete.")
	return deleteCmd
}

func addto(cmd *cobra.Command, args []string) error {
	// current user
	user, err := client.CurrentUser()
	if err != nil {
		return err
	}
	fmt.Println("User: ", user.DisplayName)

	// get current playing song
	playing, err := client.PlayerCurrentlyPlaying()
	if err != nil {
		return err
	}
	fmt.Println("Track: ", playing.Item.Name)

	// get my playlists
	pl, err := getPlaylistByName(addtoPlaylistName)
	if err != nil {
		return err
	}

	// add track to playlist
	snapshotID, err := client.AddTracksToPlaylist(user.ID, pl.ID, playing.Item.ID)
	if err != nil {
		return err
	}
	fmt.Printf("Added track to playlist %s, snapshotID: %s", pl.Name, snapshotID)
	return nil
}

func newPlaylist(cmd *cobra.Command, args []string) error {
	// current user
	user, err := client.CurrentUser()
	if err != nil {
		return err
	}
	fmt.Println("User: ", user.DisplayName)

	// create new playlist
	playlist, err := client.CreatePlaylistForUser(user.ID, newPlaylistName, true)
	if err != nil {
		return err
	}
	fmt.Println("Created public playlist: ", playlist.Name)
	return nil
}

func deletePlaylist(cmd *cobra.Command, args []string) error {
	// current user
	user, err := client.CurrentUser()
	if err != nil {
		return err
	}
	fmt.Println("User: ", user.DisplayName)

	// get the playlist
	pl, err := getPlaylistByName(delPlaylistName)
	if err != nil {
		return err
	}

	// unfollow and return
	// TODO: delete != unfollow?
	return client.UnfollowPlaylist(spotify.ID(user.ID), pl.ID)
}

func addTrackByIDToPlaylist(cmd *cobra.Command, args []string) error {
	// current user
	user, err := client.CurrentUser()
	if err != nil {
		return err
	}
	fmt.Println("User: ", user.DisplayName)

	// get the track (check for existence)
	tr, err := client.GetTrack(spotify.ID(addTrackID))
	if err != nil {
		return err
	}
	fmt.Println("Track: ", tr.Name)

	// get the playlist by name
	pl, err := getPlaylistByName(addTrackToPlaylistName)
	if err != nil {
		return err
	}

	// add track to playlist
	snapshotID, err := client.AddTracksToPlaylist(user.ID, pl.ID, tr.ID)
	if err != nil {
		return err
	}
	fmt.Printf("Added track to playlist %s, snapshotID: %s", pl.Name, snapshotID)
	return nil
}

func getPlaylistByName(playlistName string) (spotify.SimplePlaylist, error) {
	// get current user's playlists
	playlists, err := client.CurrentUsersPlaylists()
	if err != nil {
		return spotify.SimplePlaylist{}, err
	}

	// match by name
	var matchPlaylist spotify.SimplePlaylist
	for _, p := range playlists.Playlists {
		if playlistName == p.Name {
			matchPlaylist = p
			break
		}
	}

	// check if found and return
	if reflect.DeepEqual(matchPlaylist, spotify.SimplePlaylist{}) {
		return spotify.SimplePlaylist{}, fmt.Errorf("playlist not found: %s", playlistName)
	}
	return matchPlaylist, nil
}
