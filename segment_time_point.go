package mpd

import (
	"fmt"

	"github.com/moovweb/gokogiri/xml"
)

type SegmentTimePoint struct {
	/**
	 * The start time of the media segment, in seconds, relative to the beginning
	 * of the Period.
	 * @type {?number}
	 */
	StartTime uint64

	/** @type {?number} */
	Duration uint64

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
	if segmentTimePoint.StartTime, err = parseAttrAsUnsignedLong(elem, "t"); err != nil {
		segmentTimePoint.StartTime = ^uint64(0)
		fmt.Printf("%s\r\n", err.Error())
	}

	if segmentTimePoint.Duration, err = parseAttrAsUnsignedLong(elem, "d"); err != nil {
		segmentTimePoint.Duration = ^uint64(0)
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
