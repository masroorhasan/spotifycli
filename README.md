# Spotifycli

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
  rm          Remove track from playlist
  search      search tracks, albums, artists, playlists by name

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
---------------------------  ----------------------------------------------------  ----------------------------------------------------  -------------------------  ---------------
                        ID                                                  Name                                                 Album                     Artist       Popularity 
---------------------------  ----------------------------------------------------  ----------------------------------------------------  -------------------------  ---------------
    3NixJL6maKw7OxnfDOb4o9                       One Step Closer - Live In Texas                                         Live In Texas                Linkin Park               45 

    6CsrXvZhQh4I28NWMo38Cm       One Step Closer - Live LP Underground Tour 2003                                           Reanimation                Linkin Park               35 

    6CycnO02O5PIh36NNWaAyq           One Step Closer - Live from Frankfurt, 2008                   Hybrid Theory Live Around The World                Linkin Park               42 

    5KrQ4ememR3tu0CvcghC37               One Step Closer - Live At Milton Keynes             Road To Revolution: Live At Milton Keynes                Linkin Park               37 

    0mk8seWpbNcsp5ivXLw6wt                                       One Step Closer                             Live at Billy Bob's Texas                 Wade Bowen               18 

    2lvv0YF2yoxYlbhPjChRyQ                                       One Step Closer                        Live At The Greek Theater 1982        The Doobie Brothers               13 

    0npJWKQARLaaN91chFIR2d                                       One Step Closer                                                  Live             Brandon Rhyder                6 

    645zt2vQKrQzZixz3vnoZ8                                         River of Live                                       One Step Closer            One Step Closer                0 

    6WVlXciwY2Lt67cpyEF8mA                                       One Step Closer                    One Step Closer (Live & Unplugged)                  9 Red Sun                0 

    6TOvmS4KNO3p4Ny3pO1vVt                                       One Step Closer                                   In His Hands (Live)       The Heirs Of Harmony                4 

    6Sa8qEXJ3fllLhSlJ3gq48                                One Step Closer - Live                                Fantasia Live In Tokyo                       Asia                2 

    32SctXYBQcOHHjXOmTmMAV                                        Live for Today                                       One Step Closer                Gus Hergert                0 

    1B1NmYCNVBWwEZPNRHGpy6                        One Step Closer to Home - Live                                  Greatest Hits (Live)                  The Alarm                3 

    6j5I5ZHFXoMTPwkbQIJpaD                                One Step Closer - Live                                          High Voltage                       Asia                2 

    0I9BuK2HV9DnNQILQ6hQsE                     One Step Closer - Live Remastered       Electric Folklore (Live 1987-1988) [Remastered]                  The Alarm                3 

    02dvw7IeKQ0TC16381pDVT                                  One Moment of Heaven                    One Step Closer (Live & Unplugged)                  9 Red Sun                0 

    0HdKZHUNVFbxtXcNyWGPI6                               One Step Closer To Home                                    Greatest Hits Live                  The Alarm                2 

    72oQJ6QKxdbQdg3UCj2fZc                                      Be Your Own Hero                    One Step Closer (Live & Unplugged)                  9 Red Sun                0 

    2mx1XEs0AzftbjkKtSBUxn                                               You Are                    One Step Closer (Live & Unplugged)                  9 Red Sun                0 

    0ql4T68sh2RmClUiOGZZQU                                     1000 Open Windows                    One Step Closer (Live & Unplugged)                  9 Red Sun                0 
---------------------------  ----------------------------------------------------  ----------------------------------------------------  -------------------------  ---------------
```

## WIP
* Package for release.