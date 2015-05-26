package mpd

import (
	"github.com/moovweb/gokogiri/xml"
)

type Period struct {
	/** @type {?string} */
	Id string

	/**
	 * The start time of the Period, in seconds, with respect to the media
	 * presentation timeline. Note that the Period becomes/became available at
	 * Mpd.availabilityStartTime + Period.start.
	 * @type {?number}
	 */
	Start int

	/**
	 * The duration in seconds.
	 * @type {?number}
	 */
	Duration int

	/** @type {*BaseUrl} */
	BaseUrl *BaseUrl

	/** @type {SegmentBase} */
	SegmentBase *SegmentBase

	/** @type {SegmentList} */
	SegmentList *SegmentList

	/** @type {SegmentTemplate} */
	SegmentTemplate *SegmentTemplate

	/** @type {!Array.<!AdaptationSet>} */
	AdaptationSets []*AdaptationSet
}

/**
 * Parses a "Period" tag.
 * @param {!Mpd} parent The parent Mpd.
 * @param {!Node} elem The Period XML element.
 */
func (period *Period) Parse(parent Node, elem xml.Node) {
	var p *Mpd = parent.(*Mpd)
	var err error

	// Parse attributes.
	period.Id, _ = parseAttrAsString(elem, "id")

	if period.Start, err = parseAttrAsDuration(elem, "start"); err != nil {
		period.Start = -1
	}

	if period.Duration, err = parseAttrAsDuration(elem, "duration"); err != nil {
		period.Duration = -1
	}

	// Parse simple child elements.
	ok := false
	if period.BaseUrl, ok = parseChild(period, elem, BaseUrl_TAG_NAME).(*BaseUrl); ok == false {
		period.BaseUrl = p.BaseUrl
	}

	// Parse hierarchical children.
	if period.SegmentBase, ok = parseChild(period, elem, SegmentBase_TAG_NAME).(*SegmentBase); ok == false {
		period.SegmentBase = nil
	}

	if period.SegmentList, ok = parseChild(period, elem, SegmentList_TAG_NAME).(*SegmentList); ok == false {
		period.SegmentList = nil
	}

	if period.SegmentTemplate, ok = parseChild(period, elem, SegmentTemplate_TAG_NAME).(*SegmentTemplate); ok == false {
		period.SegmentTemplate = nil
	}

	children := parseChildren(period, elem, AdaptationSet_TAG_NAME)
	period.AdaptationSets = make([]*AdaptationSet, len(children))
	for i, child := range children {
		period.AdaptationSets[i] = child.(*AdaptationSet)
	}
}

func NewPeriod() Node {
	return &Period{}
}
