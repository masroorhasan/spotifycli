package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zmb3/spotify"
)

var (
	track string
	// album    string
	// artist   string
	// playlist string
)

func newSearchCmd() *cobra.Command {
	searchCmd := &cobra.Command{
		Use:   "search",
		Short: "search tracks, albums, artists, playlists by name",
		RunE: func(cmd *cobra.Command, args []string) error {
			return search(cmd, args)
		},
	}
	searchCmd.Flags().StringVar(&track, "t", "", "Name of track to search for.")
	// searchCmd.Flags().StringVar(&album, "al", "", "Name of album to search for.")
	// searchCmd.Flags().StringVar(&artist, "ar", "", "Name of artist to search for.")
	// searchCmd.Flags().StringVar(&playlist, "p", "", "Name of playlist to search for.")
	return searchCmd
}

func search(cmd *cobra.Command, args []string) error {
	results, err := client.Search(track, spotify.SearchTypeTrack)
	if err != nil {
		return err
	}
	if results.Tracks != nil {
		fmt.Println("Tracks:")
		for _, item := range results.Tracks.Tracks {
			// TODO: pretty print
			fmt.Printf("%s\t%s\t%s\t%d\n", item.ID, item.Name, item.Album.Name, item.Popularity)
		}
	}
	return nil
}
