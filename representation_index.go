package mpd

import (
	"github.com/moovweb/gokogiri/xml"
)

type RepresentationIndex struct {
	/** @type {string} */
	Url string

	/**
	 * Inherits the value of SegmentBase.indexRange if not specified.
	 * @type {Range}
	 */
	Range *Range
}

/**
 * Parses a "RepresentationIndex" tag.
 * @param {!SegmentBase} parent The parent SegmentBase.
 * @param {!Node} elem The RepresentationIndex XML element.
 */
func (representationIndex *RepresentationIndex) Parse(parent Node, elem xml.Node) {
	var err error
	var p *SegmentBase = (parent).(*SegmentBase)
	// Parse attributes.
	representationIndex.Url, _ = parseAttrAsString(elem, "sourceURL")

	if representationIndex.Range, err = parseAttrAsRange(elem, "range"); err != nil {
		representationIndex.Range = p.IndexRange.Clone()
	}
}

/**
 * Creates a deep copy of this RepresentationIndex.
 * @return {!RepresentationIndex}
 */
func (representationIndex *RepresentationIndex) Clone() Node {

	clone := &RepresentationIndex{Url: representationIndex.Url,
		Range: representationIndex.Range.Clone(),
	}

	return clone
}

func NewRepresentationIndex() Node {
	return &RepresentationIndex{}
}
