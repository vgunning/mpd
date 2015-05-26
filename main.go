package main

import (
	"fmt"
)

func main() {

	fmt.Println("start")

	if mpd, err := ParseMpd("https://sdk.streamrail.com/pepsi/cdn/0.0.1/827deadb082df0496457950ae31326eceec2e505/dash/manifest.mpd"); err != nil {
		fmt.Println(err)
	} else {
		mpdProcessor := NewMpdProcessor()
		mpdProcessor.Process(mpd)
	}

	fmt.Println("all done.")
}
