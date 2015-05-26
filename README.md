# go-mpd-parser [![Circle CI](https://circleci.com/gh/streamrail/go-mpd-parser.svg?style=svg)](https://circleci.com/gh/streamrail/go-mpd-parser)
Parse MPD files.

This is a rewrite of Google's Shaka player MPD parser from JavaScript to Go.

## Usage:

```go

// Transforms XML to Go struct
mpd, _ := ParseMpd(MPD_URL)

mpdProcessor := NewMpdProcessor()

// Construct manifest from Mpd struct
mpdProcessor.Process(mpd)

// Inspect mpdProcessor.ManifestInfo

```