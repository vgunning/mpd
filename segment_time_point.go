package mpd

import (
	"github.com/moovweb/gokogiri/xml"
)

type SegmentTimePoint struct {
	/**
	 * The start time of the media segment, in seconds, relative to the beginning
	 * of the Period.
	 * @type {?number}
	 */
	StartTime int

	/** @type {?number} */
	Duration int

	/** @type {?number} */
	Repeat int
}

/**
 * Parses an "S" tag.
 * @param {!SegmentTimeline} parent The parent SegmentTimeline.
 * @param {!Node} elem The SegmentTimePoint XML element.
 */
func (segmentTimePoint *SegmentTimePoint) Parse(parent Node, elem xml.Node) {
	var err error
	// Parse attributes.
	if segmentTimePoint.StartTime, err = parseAttrAsNonNegativeInt(elem, "t"); err != nil {
		segmentTimePoint.StartTime = -1
	}

	if segmentTimePoint.Duration, err = parseAttrAsNonNegativeInt(elem, "d"); err != nil {
		segmentTimePoint.Duration = -1
	}
	if segmentTimePoint.Repeat, err = parseAttrAsNonNegativeInt(elem, "r"); err != nil {
		segmentTimePoint.Repeat = -1
	}
}

/**
 * Creates a deep copy of this SegmentTimePoint.
 * @return {!SegmentTimePoint}
 */
func (segmentTimePoint *SegmentTimePoint) Clone() Node {
	clone := &SegmentTimePoint{
		StartTime: segmentTimePoint.StartTime,
		Duration:  segmentTimePoint.Duration,
		Repeat:    segmentTimePoint.Repeat,
	}

	return clone
}

func NewSegmentTimePoint() Node {
	return &SegmentTimePoint{}
}
