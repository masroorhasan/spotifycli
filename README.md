# Spotifycli

[![CircleCI](https://circleci.com/gh/masroorhasan/spotifycli/tree/master.svg?style=svg)](https://circleci.com/gh/masroorhasan/spotifycli/tree/master)

A command line interface to manage Spotify playlists.

## Install

To use `spotifycli` you have to register the application on Spotify's developer platform. Sign up or login available [here](https://beta.developer.spotify.com/dashboard/login). Set the following environment variables with the client Id and secret.

```
export SPOTIFY_ID=xxx
export SPOTIFY_SECRET=xxx
```

## Usage

### Commands
List of available commands:
```
$ ./spotifycli --help
A command line interface to manage Spotify playlists.

Usage:
  spotifycli [command]

Available Commands:
  add         Add track by name to playlist
  aid         Add track by ID to playlist
  ato         Add currently playing track to playlist
  del         Delete a playlist
  help        Help about any command
  list        List tracks in playlist
  login       Login to authenticate Spotify account
  logout      Logout from Spotify account
  new         Create new playlist
  now         Displays the currently playing track
  playlists   Show all playlists
  rm          Remove track from playlist
  search      search tracks, albums, artists, playlists by name
  show        Display information about a track by ID

Flags:
  -h, --help   help for spotifycli

Use "spotifycli [command] --help" for more information about a command.
```

### Search
Search using query terms on top of tracks (`tr`), albums (`al`), artists (`ar`) or playlists (`pl`) by name.

```
./spotifycli search --help
search tracks, albums, artists, playlists by name

Usage:
  spotifycli search --t [SEARCH_TYPE] --q [SEARCH_QUERY] [flags]

Flags:
  -h, --help       help for search
      --q string   The search query term.
      --t string   The search type (tr, al, ar, pl).
```

Sample search for type `tr` (track).
```
./spotifycli search --t "tr" --q "one step closer - live"
```
