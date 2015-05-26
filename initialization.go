package mpd

import (
	"github.com/moovweb/gokogiri/xml"
)

type Initialization struct {
	/** @type {String} */
	Url string

	/** @type {Range} */
	Range *Range
}

/**
 * Parses an "Initialization" tag.
 * @param {!SegmentBase|!SegmentList} parent
 *     The parent SegmentBase or parent SegmentList.
 * @param {!Node} elem The Initialization XML element.
 */
func (initialization *Initialization) Parse(parent Node, elem xml.Node) {

	// Parse attributes.
	initialization.Url, _ = parseAttrAsString(elem, "sourceURL")

	initialization.Range, _ = parseAttrAsRange(elem, "range")
}

/**
 * Creates a deep copy of this Initialization.
 * @return {!Initialization}
 */
func (initialization *Initialization) Clone() Node {

	clone := &Initialization{Url: initialization.Url,
		Range: initialization.Range.Clone(),
	}

	return clone
}

func NewInitialization() Node {
	return &Initialization{}
}
