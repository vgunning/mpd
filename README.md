# mpd
Parse MPD files.

This is a rewrite of [Google's Shaka player][] MPD parser from JavaScript to Go.

[Google's Shaka player]: https://github.com/google/shaka-player

## Usage:
 
```go

// Transforms XML to Go struct
mpd, _ := ParseMpd("http://sdk.streamrail.com/pepsi/cdn/0.0.1/601486e52319059b8790c13f7477d2036d042768/dash/manifest.mpd")

// Print parsed mpd
PrintMPD(mpd, 0)

mpdProcessor := NewMpdProcessor()

// Construct manifest from Mpd struct
mpdProcessor.Process(mpd)

// Inspect mpdProcessor.ManifestInfo

```

The snipet above parse given mpd (which you can watch [here][])
[here]: http://play.streamrail.com/#/vjs

PrintMPD prooduces the following output:

```go

main.Mpd
	Id =
	Type = static
	BaseUrl = &{http://sdk.streamrail.com/pepsi/cdn/0.0.1/601486e52319059b8790c13f7477d2036d042768/dash/}
	main.BaseUrl
		Url = http://sdk.streamrail.com/pepsi/cdn/0.0.1/601486e52319059b8790c13f7477d2036d042768/dash/
	MediaPresentationDuration = 126
	MinBufferTime = 5
	MinUpdatePeriod = 0
	AvailabilityStartTime = -1
	TimeShiftBufferDepth = 0
	SuggestedPresentationDelay = 1
	Periods = [0xc208064060]
	main.Period
		Id =
		Start = -1
		Duration = -1
		BaseUrl = &{http://sdk.streamrail.com/pepsi/cdn/0.0.1/601486e52319059b8790c13f7477d2036d042768/dash/}
		main.BaseUrl
			Url = http://sdk.streamrail.com/pepsi/cdn/0.0.1/601486e52319059b8790c13f7477d2036d042768/dash/
		SegmentBase = <nil>
		SegmentList = <nil>
		SegmentTemplate = <nil>
		AdaptationSets = [0xc208066280 0xc2080663c0]
		main.AdaptationSet
			Id =
			Lang = und
			ContentType = audio
			Width = 0
			Height = 0
			MimeType = audio/mp4
			Codecs =
			Main = false
			BaseUrl = &{http://sdk.streamrail.com/pepsi/cdn/0.0.1/601486e52319059b8790c13f7477d2036d042768/dash/}
			main.BaseUrl
				Url = http://sdk.streamrail.com/pepsi/cdn/0.0.1/601486e52319059b8790c13f7477d2036d042768/dash/
			SegmentBase = <nil>
			SegmentList = <nil>
			SegmentTemplate = &{1000 -1 1968 1 $RepresentationID$/audio/und/seg-$Number$.m4f  $RepresentationID$/audio/und/init.mp4 <nil>}
			main.SegmentTemplate
				Timescale = 1000
				PresentationTimeOffset = -1
				SegmentDuration = 1968
				StartNumber = 1
				MediaUrlTemplate = $RepresentationID$/audio/und/seg-$Number$.m4f
				IndexUrlTemplate =
				InitializationUrlTemplate = $RepresentationID$/audio/und/init.mp4
				Timeline = <nil>
			Representations = [0xc208066320]
			main.Representation
				Id = 700k
				Lang = und
				Bandwidth = 114244
				Width = 0
				Height = 0
				MimeType = audio/mp4
				Codecs = mp4a.40.2
				BaseUrl = &{http://sdk.streamrail.com/pepsi/cdn/0.0.1/601486e52319059b8790c13f7477d2036d042768/dash/}
				main.BaseUrl
					Url = http://sdk.streamrail.com/pepsi/cdn/0.0.1/601486e52319059b8790c13f7477d2036d042768/dash/
				SegmentBase = <nil>
				SegmentList = <nil>
				SegmentTemplate = &{1000 -1 1968 1 $RepresentationID$/audio/und/seg-$Number$.m4f  $RepresentationID$/audio/und/init.mp4 <nil>}
				main.SegmentTemplate
					Timescale = 1000
					PresentationTimeOffset = -1
					SegmentDuration = 1968
					StartNumber = 1
					MediaUrlTemplate = $RepresentationID$/audio/und/seg-$Number$.m4f
					IndexUrlTemplate =
					InitializationUrlTemplate = $RepresentationID$/audio/und/init.mp4
					Timeline = <nil>
				ContentProtections = []
				Main = false
		main.AdaptationSet
			Id =
			Lang =
			ContentType = video
			Width = 0
			Height = 0
			MimeType = video/mp4
			Codecs =
			Main = false
			BaseUrl = &{http://sdk.streamrail.com/pepsi/cdn/0.0.1/601486e52319059b8790c13f7477d2036d042768/dash/}
			main.BaseUrl
				Url = http://sdk.streamrail.com/pepsi/cdn/0.0.1/601486e52319059b8790c13f7477d2036d042768/dash/
			SegmentBase = <nil>
			SegmentList = <nil>
			SegmentTemplate = &{1000 -1 1968 1 $RepresentationID$/video/1/seg-$Number$.m4f  $RepresentationID$/video/1/init.mp4 <nil>}
			main.SegmentTemplate
				Timescale = 1000
				PresentationTimeOffset = -1
				SegmentDuration = 1968
				StartNumber = 1
				MediaUrlTemplate = $RepresentationID$/video/1/seg-$Number$.m4f
				IndexUrlTemplate =
				InitializationUrlTemplate = $RepresentationID$/video/1/init.mp4
				Timeline = <nil>
			Representations = [0xc208066460 0xc208066500 0xc2080665a0]
			main.Representation
				Id = 700k
				Lang =
				Bandwidth = 948337
				Width = 1920
				Height = 1080
				MimeType = video/mp4
				Codecs = avc1.42c028
				BaseUrl = &{http://sdk.streamrail.com/pepsi/cdn/0.0.1/601486e52319059b8790c13f7477d2036d042768/dash/}
				main.BaseUrl
					Url = http://sdk.streamrail.com/pepsi/cdn/0.0.1/601486e52319059b8790c13f7477d2036d042768/dash/
				SegmentBase = <nil>
				SegmentList = <nil>
				SegmentTemplate = &{1000 -1 1968 1 $RepresentationID$/video/1/seg-$Number$.m4f  $RepresentationID$/video/1/init.mp4 <nil>}
				main.SegmentTemplate
					Timescale = 1000
					PresentationTimeOffset = -1
					SegmentDuration = 1968
					StartNumber = 1
					MediaUrlTemplate = $RepresentationID$/video/1/seg-$Number$.m4f
					IndexUrlTemplate =
					InitializationUrlTemplate = $RepresentationID$/video/1/init.mp4
					Timeline = <nil>
				ContentProtections = []
				Main = false
			main.Representation
				Id = 1200k
				Lang =
				Bandwidth = 1467491
				Width = 1920
				Height = 1080
				MimeType = video/mp4
				Codecs = avc1.42c028
				BaseUrl = &{http://sdk.streamrail.com/pepsi/cdn/0.0.1/601486e52319059b8790c13f7477d2036d042768/dash/}
				main.BaseUrl
					Url = http://sdk.streamrail.com/pepsi/cdn/0.0.1/601486e52319059b8790c13f7477d2036d042768/dash/
				SegmentBase = <nil>
				SegmentList = <nil>
				SegmentTemplate = &{1000 -1 1968 1 $RepresentationID$/video/1/seg-$Number$.m4f  $RepresentationID$/video/1/init.mp4 <nil>}
				main.SegmentTemplate
					Timescale = 1000
					PresentationTimeOffset = -1
					SegmentDuration = 1968
					StartNumber = 1
					MediaUrlTemplate = $RepresentationID$/video/1/seg-$Number$.m4f
					IndexUrlTemplate =
					InitializationUrlTemplate = $RepresentationID$/video/1/init.mp4
					Timeline = <nil>
				ContentProtections = []
				Main = false
			main.Representation
				Id = 4508k
				Lang =
				Bandwidth = 6963128
				Width = 1920
				Height = 1080
				MimeType = video/mp4
				Codecs = avc1.42c028
				BaseUrl = &{http://sdk.streamrail.com/pepsi/cdn/0.0.1/601486e52319059b8790c13f7477d2036d042768/dash/}
				main.BaseUrl
					Url = http://sdk.streamrail.com/pepsi/cdn/0.0.1/601486e52319059b8790c13f7477d2036d042768/dash/
				SegmentBase = <nil>
				SegmentList = <nil>
				SegmentTemplate = &{1000 -1 1968 1 $RepresentationID$/video/1/seg-$Number$.m4f  $RepresentationID$/video/1/init.mp4 <nil>}
				main.SegmentTemplate
					Timescale = 1000
					PresentationTimeOffset = -1
					SegmentDuration = 1968
					StartNumber = 1
					MediaUrlTemplate = $RepresentationID$/video/1/seg-$Number$.m4f
					IndexUrlTemplate =
					InitializationUrlTemplate = $RepresentationID$/video/1/init.mp4
					Timeline = <nil>
				ContentProtections = []
				Main = false

```

For more examples see mpd_processor_test.go