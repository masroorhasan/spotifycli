package cmd

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/bndr/gotabulate"
	"github.com/spf13/cobra"
	"github.com/zmb3/spotify"
)

var (
	searchQuery string
	searchType  string
)

func newSearchCmd() *cobra.Command {
	searchCmd := &cobra.Command{
		Use:   "search --t [SEARCH_TYPE] --q [SEARCH_QUERY]",
		Short: "search tracks, albums, artists, playlists by name",
		RunE: func(cmd *cobra.Command, args []string) error {
			return search(cmd, args)
		},
	}
	searchCmd.Flags().StringVar(&searchType, "t", "", "The search type (tr, al, ar, pl).")
	searchCmd.Flags().StringVar(&searchQuery, "q", "", "The search query.")
	return searchCmd
}

func search(cmd *cobra.Command, args []string) error {
	switch searchType {
	case "tr":
		return searchTracks(searchQuery)
	case "al":
		return searchAlbums(searchQuery)
	case "ar":
		return searchArtists(searchQuery)
	case "pl":
		return searchPlaylists(searchQuery)
	default:
		return errors.New("Not supported")
	}
}

func searchTracks(query string) error {
	results, err := client.Search(query, spotify.SearchTypeTrack)
	if err != nil {
		return err
	}

	var data [][]interface{}
	if results.Tracks != nil {
		for _, item := range results.Tracks.Tracks {
			track := []string{string(item.ID), item.Name, item.Album.Name, item.Artists[0].Name, strconv.Itoa(item.Popularity)}
			row := make([]interface{}, len(track))
			for i, d := range track {
				row[i] = d
			}
			data = append(data, row)
		}
	}
	printSimple([]string{"ID", "Name", "Album", "Artist", "Popularity"}, data)
	return nil
}

func searchAlbums(query string) error {
	results, err := client.Search(query, spotify.SearchTypeAlbum)
	if err != nil {
		return err
	}

	var data [][]interface{}
	if results.Albums != nil {
		for _, item := range results.Albums.Albums {
			album := []string{string(item.ID), item.Name, item.Artists[0].Name, item.AlbumType, item.Endpoint}
			row := make([]interface{}, len(album))
			for i, d := range album {
				row[i] = d
			}
			data = append(data, row)
		}
	}
	printSimple([]string{"ID", "Name", "Artist", "Type", "Endpoint"}, data)
	return nil
}

func searchArtists(query string) error {
	results, err := client.Search(query, spotify.SearchTypeArtist)
	if err != nil {
		return nil
	}

	var data [][]interface{}
	if results.Artists != nil {
		for _, item := range results.Artists.Artists {
			artist := []string{string(item.ID), item.Name, strings.Join(item.Genres, ","), strconv.Itoa(int(item.Followers.Count)), item.Endpoint}
			row := make([]interface{}, len(artist))
			for i, d := range artist {
				row[i] = d
			}
			data = append(data, row)
		}
	}
	printSimple([]string{"ID", "Name", "Genres", "Followers", "Endpoint"}, data)
	return nil
}

func searchPlaylists(query string) error {
	results, err := client.Search(query, spotify.SearchTypePlaylist)
	if err != nil {
		return nil
	}

	var data [][]interface{}
	if results.Playlists != nil {
		for _, item := range results.Playlists.Playlists {
			playlist := []string{string(item.ID), item.Name, item.Owner.DisplayName, strconv.Itoa(int(item.Tracks.Total)), item.Endpoint}
			row := make([]interface{}, len(playlist))
			for i, d := range playlist {
				row[i] = d
			}
			data = append(data, row)
		}
	}
	printSimple([]string{"ID", "Name", "Owner", "Total Tracks", "Endpoint"}, data)
	return nil
}

func printSimple(headers []string, data [][]interface{}) {
	tabulate := gotabulate.Create(data)
	tabulate.SetHeaders(headers)
	fmt.Println(tabulate.Render("simple"))
}
