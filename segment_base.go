package mpd

import (
	"github.com/moovweb/gokogiri/xml"
)

type SegmentBase struct {
	/**
	 * This not an actual XML attribute of SegmentBase. It is inherited from the
	 * SegmentBase's parent Representation.
	 * @type {*BaseUrl}
	 */
	BaseUrl *BaseUrl

	/** @type {?number} */
	Timescale int

	/** @type {?number} */
	PresentationTimeOffset int

	/** @type {Range} */
	IndexRange *Range

	/** @type {RepresentationIndex} */
	RepresentationIndex *RepresentationIndex

	/** @type {Initialization} */
	Initialization *Initialization
}

/**
 * Parses a "SegmentBase" tag.
 * @param {*} parent The parent object.
 * @param {!Node} elem The SegmentBase XML element.
 */
func (segmentBase *SegmentBase) Parse(parent Node, elem xml.Node) {

	switch p := parent.(type) {
	case *AdaptationSet:
	case *Period:
	case *Representation:
		segmentBase.BaseUrl = p.BaseUrl
	}

	var err error

	// When parsing attributes and child elements fallback to |this| to provide
	// default values. If |this| is a new SegmentBase then |this| will have
	// default values from its constructor, and if |this| was cloned from a
	// higher level SegmentBase then |this| will have values from that
	// SegmentBase.
	segmentBase.Timescale, _ = parseAttrAsPositiveInt(elem, "timescale")

	if segmentBase.PresentationTimeOffset, err = parseAttrAsNonNegativeInt(elem, "presentationTimeOffset"); err != nil {
		segmentBase.PresentationTimeOffset = -1
	}

	// Parse attributes.
	segmentBase.IndexRange, _ = parseAttrAsRange(elem, "indexRange")

	// Parse simple child elements.
	ok := false
	if segmentBase.RepresentationIndex, ok = parseChild(segmentBase, elem, RepresentationIndex_TAG_NAME).(*RepresentationIndex); ok == false {
		segmentBase.RepresentationIndex = nil
	}

	if segmentBase.Initialization, ok = parseChild(segmentBase, elem, Initialization_TAG_NAME).(*Initialization); ok == false {
		segmentBase.Initialization = nil
	}
}

/**
 * Creates a deep copy of this SegmentBase.
 * @return {!SegmentBase}
 */
func (segmentBase SegmentBase) Clone() Node {
	var representationIndex *RepresentationIndex = nil
	var initialization *Initialization = nil

	if segmentBase.RepresentationIndex != nil {
		representationIndex = segmentBase.RepresentationIndex.Clone().(*RepresentationIndex)
	}

	if segmentBase.Initialization != nil {
		initialization = segmentBase.Initialization.Clone().(*Initialization)
	}

	clone := &SegmentBase{
		BaseUrl:                segmentBase.BaseUrl,
		Timescale:              segmentBase.Timescale,
		PresentationTimeOffset: segmentBase.PresentationTimeOffset,
		IndexRange:             segmentBase.IndexRange.Clone(),
		RepresentationIndex:    representationIndex,
		Initialization:         initialization,
	}

	return clone
}

func NewSegmentBase() Node {
	return &SegmentBase{}
}
