package mpd

import (
	"github.com/moovweb/gokogiri/xml"
)

type SegmentTimeline struct {
	/** @type {!Array.<!SegmentTimePoint>} */
	TimePoints []*SegmentTimePoint
}

/**
 * Parses a "SegmentTimeline" tag.
 * @param {!SegmentTemplate} parent The parent SegmentTemplate.
 * @param {!Node} elem The SegmentTimeline XML element.
 */
func (segmentTimeline *SegmentTimeline) Parse(parent Node, elem xml.Node) {
	children := parseChildren(segmentTimeline, elem, SegmentTimePoint_TAG_NAME)
	segmentTimeline.TimePoints = make([]*SegmentTimePoint, len(children))

	for i, child := range children {
		segmentTimeline.TimePoints[i] = (child).(*SegmentTimePoint)
	}
}

/**
 * Creates a deep copy of this SegmentTimeline.
 * @return {!SegmentTimeline}
 */
func (segmentTimeline *SegmentTimeline) Clone() Node {
	clone := &SegmentTimeline{
		TimePoints: make([]*SegmentTimePoint, 0),
	}

	for _, timePoint := range segmentTimeline.TimePoints {
		clone.TimePoints = append(clone.TimePoints, timePoint.Clone().(*SegmentTimePoint))
	}

	return clone
}

func NewSegmentTimeline() Node {
	return &SegmentTimeline{}
}
