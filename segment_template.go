package main

import (
	"github.com/moovweb/gokogiri/xml"
)

type SegmentTemplate struct {
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

	/** @type {?string} */
	MediaUrlTemplate string

	/** @type {?string} */
	IndexUrlTemplate string

	/** @type {?string} */
	InitializationUrlTemplate string

	/** @type {SegmentTimeline} */
	Timeline *SegmentTimeline
}

/**
 * Parses a "SegmentTemplate" tag.
 * @param {*} parent The parent object.
 * @param {!Node} elem The SegmentTemplate XML element.
 */
func (segmentTemplate *SegmentTemplate) Parse(parent Node, elem xml.Node) {
	var err error

	// Parse attributes.
	segmentTemplate.Timescale, _ = parseAttrAsPositiveInt(elem, "timescale")

	if segmentTemplate.PresentationTimeOffset, err = parseAttrAsNonNegativeInt(elem, "presentationTimeOffset"); err != nil {
		segmentTemplate.PresentationTimeOffset = -1
	}

	if segmentTemplate.SegmentDuration, err = parseAttrAsPositiveInt(elem, "duration"); err != nil {
		segmentTemplate.SegmentDuration = -1
	}

	if segmentTemplate.StartNumber, err = parseAttrAsPositiveInt(elem, "startNumber"); err != nil {
		segmentTemplate.StartNumber = 1
	}

	if segmentTemplate.MediaUrlTemplate, err = parseAttrAsString(elem, "media"); err != nil {
		segmentTemplate.MediaUrlTemplate = ""
	}

	if segmentTemplate.IndexUrlTemplate, err = parseAttrAsString(elem, "index"); err != nil {
		segmentTemplate.IndexUrlTemplate = ""
	}

	if segmentTemplate.InitializationUrlTemplate, err = parseAttrAsString(elem, "initialization"); err != nil {
		segmentTemplate.InitializationUrlTemplate = ""
	}

	// Parse hierarchical children.
	ok := false
	if segmentTemplate.Timeline, ok = parseChild(segmentTemplate, elem, SegmentTimeline_TAG_NAME).(*SegmentTimeline); ok == false {
		segmentTemplate.Timeline = nil
	}
}

func (segmentTemplate SegmentTemplate) Clone() Node {

	var timeLine *SegmentTimeline = nil
	if segmentTemplate.Timeline != nil {
		timeLine = segmentTemplate.Timeline.Clone().(*SegmentTimeline)
	}

	return &SegmentTemplate{
		Timescale:                 segmentTemplate.Timescale,
		PresentationTimeOffset:    segmentTemplate.PresentationTimeOffset,
		SegmentDuration:           segmentTemplate.SegmentDuration,
		StartNumber:               segmentTemplate.StartNumber,
		MediaUrlTemplate:          segmentTemplate.MediaUrlTemplate,
		IndexUrlTemplate:          segmentTemplate.IndexUrlTemplate,
		InitializationUrlTemplate: segmentTemplate.InitializationUrlTemplate,
		Timeline:                  timeLine,
	}
}

func NewSegmentTemplate() Node {
	return &SegmentTemplate{}
}
