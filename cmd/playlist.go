package cmd

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/zmb3/spotify"
)

var (
	addtoPlaylistName string
)

var (
	addTrackID             string
	addTrackToPlaylistName string
)

var (
	rmTrackName             string
	rmTrackFromPlaylistName string
)

var (
	newPlaylistName string
)

var (
	delPlaylistName string
)

var (
	listPlaylistTracksName string
)

func newAddtoPlaylistCmd() *cobra.Command {
	addtoCmd := &cobra.Command{
		Use:   "addto --p [PLAYLIST_NAME]",
		Short: "Add currently playing track to playlist",
		RunE: func(cmd *cobra.Command, args []string) error {
			return addto(cmd, args)
		},
	}
	addtoCmd.Flags().StringVar(&addtoPlaylistName, "p", "", "Add current track to specified playlist.")
	return addtoCmd
}

func newAddTrackByIDToPlaylistCmd() *cobra.Command {
	addCmd := &cobra.Command{
		Use:   "add --tid [TRACK_ID] --p [PLAYLIST_NAME]",
		Short: "Add track to playlist",
		RunE: func(cmd *cobra.Command, args []string) error {
			return addTrackByIDToPlaylist(cmd, args)
		},
	}
	addCmd.Flags().StringVar(&addTrackID, "tid", "", "Id of track to add to playlist.")
	addCmd.Flags().StringVar(&addTrackToPlaylistName, "p", "", "Name of playlist to add track to.")
	return addCmd
}

func newRemoveTrackFromPlaylistCmd() *cobra.Command {
	rmCmd := &cobra.Command{
		Use:   "rm --t [TRACK_NAME] --p [PLAYLIST_NAME]",
		Short: "Remove track from playlist",
		RunE: func(cmd *cobra.Command, args []string) error {
			return rmTrackByNameFromPlaylist(cmd, args)
		},
	}
	rmCmd.Flags().StringVar(&rmTrackName, "t", "", "Name of track to remove.")
	rmCmd.Flags().StringVar(&rmTrackFromPlaylistName, "p", "", "Name of playlist to remove track from.")
	return rmCmd
}

func newCreatePlaylistCmd() *cobra.Command {
	newCmd := &cobra.Command{
		Use:   "new --p [PLAYLIST_NAME]",
		Short: "Create new playlist",
		RunE: func(cmd *cobra.Command, args []string) error {
			return newPlaylist(cmd, args)
		},
	}
	newCmd.Flags().StringVar(&newPlaylistName, "p", "", "Name of new playlist.")
	return newCmd
}

func newDeletePlaylistCmd() *cobra.Command {
	deleteCmd := &cobra.Command{
		Use:   "del --p [PLAYLIST_NAME]",
		Short: "Delete a playlist",
		RunE: func(cmd *cobra.Command, args []string) error {
			return deletePlaylist(cmd, args)
		},
	}
	deleteCmd.Flags().StringVar(&delPlaylistName, "p", "", "Name of playlist to delete.")
	return deleteCmd
}

func newListPlaylistTracksCmd() *cobra.Command {
	listCmd := &cobra.Command{
		Use:   "list --p [PLAYLIST_NAME]",
		Short: "List tracks in playlist",
		RunE: func(cmd *cobra.Command, args []string) error {
			return listTracksFromPlaylist(cmd, args)
		},
	}
	listCmd.Flags().StringVar(&listPlaylistTracksName, "p", "", "Name of playlist to list tracks from.")
	return listCmd
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

func rmTrackByNameFromPlaylist(cmd *cobra.Command, args []string) error {
	// current user
	user, err := client.CurrentUser()
	if err != nil {
		return err
	}
	fmt.Println("User: ", user.DisplayName)

	// get the playlist by name
	pl, err := getPlaylistByName(rmTrackFromPlaylistName)
	if err != nil {
		return err
	}

	// get track in playlist and validate existence
	var matchedTrack spotify.SimpleTrack
	ptracks, err := client.GetPlaylistTracks(user.ID, pl.ID)
	for _, t := range ptracks.Tracks {
		if rmTrackName == t.Track.SimpleTrack.Name {
			matchedTrack = t.Track.SimpleTrack
			break
		}
	}
	if reflect.DeepEqual(matchedTrack, spotify.SimpleTrack{}) {
		return fmt.Errorf("track %s not found in playlist %s", rmTrackName, rmTrackFromPlaylistName)
	}

	// remove track from playlist
	snapshotID, err := client.RemoveTracksFromPlaylist(user.ID, pl.ID, matchedTrack.ID)
	if err != nil {
		return err
	}
	fmt.Printf("Removed track %s from playlist %s, snapshotID %s", rmTrackName, rmTrackFromPlaylistName, snapshotID)
	return nil
}

func listTracksFromPlaylist(cmd *cobra.Command, args []string) error {
	// current user
	user, err := client.CurrentUser()
	if err != nil {
		return err
	}
	fmt.Println("User: ", user.DisplayName)

	pl, err := getPlaylistByName(listPlaylistTracksName)
	if err != nil {
		return err
	}

	// get tracks from playlist
	tracks, err := client.GetPlaylistTracks(user.ID, pl.ID)
	if err != nil {
		return err
	}

	// format resulting data
	var data [][]interface{}
	if tracks.Tracks != nil {
		for _, item := range tracks.Tracks {
			track := []string{
				string(item.Track.ID),
				item.Track.Name,
				item.Track.Album.Name,
				item.Track.Artists[0].Name,
				strconv.Itoa(item.Track.Popularity)}
			row := make([]interface{}, len(track))
			for i, d := range track {
				row[i] = d
			}
			data = append(data, row)
		}
	}

	// pretty print track results
	printSimple([]string{"ID", "Name", "Album", "Artist", "Popularity"}, data)
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
