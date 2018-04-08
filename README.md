# Spotifycli

A command line interface to manage Spotify playlists.

## Install

To use `spotifycli` you have to register the application on Spotify's developer platform. Sign up or login available [here](https://beta.developer.spotify.com/dashboard/login). Set the following environment variables with the client Id and secret. 

```
export SPOTIFY_ID=xxx
export SPOTIFY_SECRET=xxx
```

## Usage

List of available commands:
```
$ ./spotifycli -h
A command line interface to manage Spotify playlists.

Usage:
  spotifycli [command]

Available Commands:
  add         Add track to playlist
  addto       Add currently playing track to playlist
  del         Delete a playlist
  help        Help about any command
  login       Login to authenticate Spotify account
  logout      Logout from Spotify account
  new         Create new playlist
  search      search tracks, albums, artists, playlists by name

Flags:
  -h, --help   help for spotifycli

Use "spotifycli [command] --help" for more information about a command.
```

## WIP
* Complex playlist ops.
* Package for release.