package main

import (
	"fmt"
)

func main() {

	fmt.Println("start")

	if mpd, err := ParseMpd("example2.xml", "streamrail.com/"); err != nil {
		fmt.Println(err)
	} else {
		mpdProcessor := NewMpdProcessor()
		mpdProcessor.Process(*mpd)
	}

	fmt.Println("all done.")
}
