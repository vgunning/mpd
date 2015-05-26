package mpd

import (
	"github.com/moovweb/gokogiri/xml"
)

type FakeNode struct {
	/** @type {*BaseUrl} */
	BaseUrl *BaseUrl
}

func (fakeNode FakeNode) Parse(parent Node, elem xml.Node) {

}
