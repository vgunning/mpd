package mpd

import (
	"fmt"
)

func assert(exp bool) {
	if !exp {
		fmt.Println("Assert evaluated to False")
	}
}
