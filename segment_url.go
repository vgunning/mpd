package main

import (
	"github.com/moovweb/gokogiri/xml"
)

type SegmentUrl struct {
	/** @type {string} */
	MediaUrl string

	/** @type {Range} */
	MediaRange *Range
}

/**
 * Parses a "SegmentUrl" tag.
 * @param {!SegmentList} parent The parent SegmentList.
 * @param {!Node} elem The SegmentUrl XML element.
 */
func (segmentUrl *SegmentUrl) Parse(parent Node, elem xml.Node) {

	// Parse attributes.
	segmentUrl.MediaUrl, _ = parseAttrAsString(elem, "media")

	segmentUrl.MediaRange, _ = parseAttrAsRange(elem, "mediaRange")
}

/**
 * Creates a deep copy of this SegmentUrl.
 * @return {!SegmentUrl}
 */
func (segmentUrl *SegmentUrl) Clone() Node {

	clone := &SegmentUrl{MediaUrl: segmentUrl.MediaUrl,
		MediaRange: segmentUrl.MediaRange.Clone(),
	}

	return clone
}

func NewSegmentUrl() Node {
	return &SegmentUrl{}
}
