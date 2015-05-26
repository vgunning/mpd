package mpd

import (
	"github.com/moovweb/gokogiri/xml"
)

type SegmentList struct {
	/**
	 * This not an actual XML attribute of SegmentList. It is inherited from the
	 * SegmentList's parent Representation.
	 * @type {BaseUrl}
	 */
	BaseUrl *BaseUrl

	/** @type {?number} */
	Timescale int

	/** @type {?number} */
	PresentationTimeOffset int

	/**
	 * Each segment's duration. This value is never zero.
	 * @type {?number}
	 */
	SegmentDuration int

	/**
	 * The segment number origin. This value is never zero.
	 * @type {?number}
	 */
	StartNumber int

	/** @type {Initialization} */
	Initialization *Initialization

	/** @type {!Array.<SegmentUrl>} */
	SegmentUrls []*SegmentUrl
}

/**
 * Parses a "SegmentList" tag.
 * @param {*} parent The parent object.
 * @param {!Node} elem The SegmentList XML element.
 */
func (segmentList *SegmentList) Parse(parent Node, elem xml.Node) {
	var err error

	switch p := parent.(type) {
	case *AdaptationSet:
	case *Period:
	case *Representation:
		segmentList.BaseUrl = p.BaseUrl
	}

	// Parse attributes.
	segmentList.Timescale, _ = parseAttrAsPositiveInt(elem, "timescale")

	segmentList.PresentationTimeOffset, _ = parseAttrAsNonNegativeInt(elem, "presentationTimeOffset")

	if segmentList.SegmentDuration, err = parseAttrAsPositiveInt(elem, "duration"); err != nil {
		segmentList.SegmentDuration = -1
	}

	if segmentList.StartNumber, err = parseAttrAsPositiveInt(elem, "startNumber"); err != nil {
		segmentList.StartNumber = 1
	}

	// Parse simple children
	segmentList.Initialization = parseChild(segmentList, elem, Initialization_TAG_NAME).(*Initialization)

	children := parseChildren(segmentList, elem, SegmentUrl_TAG_NAME)
	segmentList.SegmentUrls = make([]*SegmentUrl, len(children))
	for i, child := range children {
		segmentList.SegmentUrls[i] = child.(*SegmentUrl)
	}
}

/**
 * Creates a deep copy of this SegmentList.
 * @return {!SegmentList}
 */
func (segmentList SegmentList) Clone() Node {
	clone := &SegmentList{}

	clone.BaseUrl = segmentList.BaseUrl
	clone.Timescale = segmentList.Timescale
	clone.PresentationTimeOffset = segmentList.PresentationTimeOffset
	clone.SegmentDuration = segmentList.SegmentDuration
	clone.StartNumber = segmentList.StartNumber
	clone.Initialization = segmentList.Initialization.Clone().(*Initialization)

	for _, segmentUrl := range segmentList.SegmentUrls {
		clone.SegmentUrls = append(clone.SegmentUrls, segmentUrl.Clone().(*SegmentUrl))
	}

	return clone
}

func NewSegmentList() Node {
	return &SegmentList{}
}
