package mpd

import (
	"github.com/moovweb/gokogiri/xml"
)

type Representation struct {
	/** @type {?string} */
	Id string

	/**
	 * Never seen on the Representation itself, but inherited from AdapationSet
	 * for convenience.
	 * @see AdaptationSet.lang
	 * @type {?string}
	 */
	Lang string

	/**
	 * Bandwidth required, in bits per second, to assure uninterrupted playback,
	 * assuming that |minBufferTime| seconds of video are in buffer before
	 * playback begins.
	 * @type {?number}
	 */
	Bandwidth uint32

	/** @type {?number} */
	Width int

	/** @type {?number} */
	Height int

	/** @type {?string} */
	MimeType string

	/** @type {?string} */
	Codecs string

	/** @type {*BaseUrl} */
	BaseUrl *BaseUrl

	/** @type {SegmentBase} */
	SegmentBase *SegmentBase

	/** @type {SegmentList} */
	SegmentList *SegmentList

	/** @type {SegmentTemplate} */
	SegmentTemplate *SegmentTemplate

	/** @type {!Array.<ContentProtection>} */
	ContentProtections []*ContentProtection

	/** @type {boolean} */
	Main bool
}

/**
 * Parses a "Representation" tag.
 * @param {!AdaptationSet} parent The parent AdaptationSet.
 * @param {!Node} elem The Representation XML element.
 */
func (representation *Representation) Parse(parent Node, elem xml.Node) {
	var err error
	var p *AdaptationSet = (parent).(*AdaptationSet)

	// Parse attributes.
	representation.Id, _ = parseAttrAsString(elem, "id")
	representation.Bandwidth, _ = parseAttrAsUnsignedInt(elem, "bandwidth")
	if representation.Width, err = parseAttrAsPositiveInt(elem, "width"); err != nil {
		representation.Width = p.Width
	}

	if representation.Height, err = parseAttrAsPositiveInt(elem, "height"); err != nil {
		representation.Height = p.Height
	}

	if representation.MimeType, err = parseAttrAsString(elem, "mimeType"); err != nil {
		representation.MimeType = p.MimeType
	}

	if representation.Codecs, err = parseAttrAsString(elem, "codecs"); err != nil {
		representation.Codecs = p.Codecs
	}

	// Never seen on this element itself, but inherited for convenience.
	representation.Lang = p.Lang

	// Parse simple child elements.
	ok := false
	if representation.BaseUrl, ok = parseChild(representation, elem, BaseUrl_TAG_NAME).(*BaseUrl); ok == false {
		representation.BaseUrl = p.BaseUrl
	}

	// representation.ContentProtections = parseChildren(representation, elem, ContentProtection_TAG_NAME)

	// Parse hierarchical children.
	if p.SegmentBase != nil {
		representation.SegmentBase = mergeChild(representation, elem, p.SegmentBase, SegmentBase_TAG_NAME).(*SegmentBase)
	} else {
		if representation.SegmentBase, ok = parseChild(representation, elem, SegmentBase_TAG_NAME).(*SegmentBase); ok == false {
			representation.SegmentBase = nil
		}
	}

	if p.SegmentList != nil {
		representation.SegmentList = mergeChild(representation, elem, p.SegmentList, SegmentList_TAG_NAME).(*SegmentList)
	} else {
		if representation.SegmentList, ok = parseChild(representation, elem, SegmentList_TAG_NAME).(*SegmentList); ok == false {
			representation.SegmentList = nil
		}
	}

	if p.SegmentTemplate != nil {
		representation.SegmentTemplate = mergeChild(representation, elem, p.SegmentTemplate, SegmentTemplate_TAG_NAME).(*SegmentTemplate)
	} else {
		if representation.SegmentTemplate, ok = parseChild(representation, elem, SegmentTemplate_TAG_NAME).(*SegmentTemplate); ok == false {
			representation.SegmentTemplate = nil
		}
	}

	// if len(representation.ContentProtections) == 0 {
	// 	representation.ContentProtections = p.ContentProtections
	// }
}

func NewRepresentation() Node {
	return &Representation{}
}
