package main

import (
	"fmt"
)

func main() {

	fmt.Println("start")

	if mpd, err := ParseMpd("http://sdk.streamrail.com/pepsi/cdn/0.0.1/601486e52319059b8790c13f7477d2036d042768/dash/manifest.mpd"); err != nil {
		fmt.Println(err)
	} else {
		PrintMPD(mpd, 0)
		mpdProcessor := NewMpdProcessor()
		mpdProcessor.Process(mpd)
	}

	fmt.Println("all done.")
}
