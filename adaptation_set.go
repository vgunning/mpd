package main

import (
	"github.com/moovweb/gokogiri/xml"
	"strings"
)

type AdaptationSet struct {
	/** @type {?string} */
	Id string

	/**
	 * The language.
	 * @type {?string}
	 * @see IETF RFC 5646
	 * @see ISO 639
	 */
	Lang string

	/**
	 * Should be 'video' or 'audio', not a MIME type.
	 * If not specified, will be inferred from the MIME type.
	 * @type {?string}
	 */
	ContentType string

	/** @type {?number} */
	Width int

	/** @type {?number} */
	Height int

	/**
	 * If not specified, will be inferred from the first representation.
	 * @type {?string}
	 */
	MimeType string

	/** @type {?string} */
	Codecs string

	/** @type {!bool}*/
	Main bool

	/** @type {*BaseUrl} */
	BaseUrl *BaseUrl

	/** @type {SegmentBase} */
	SegmentBase *SegmentBase

	/** @type {SegmentList} */
	SegmentList *SegmentList

	/** @type {SegmentTemplate} */
	SegmentTemplate *SegmentTemplate

	/** @type {!Array.<!ContentProtection>} */
	// ContentProtections []ContentProtection

	/** @type {!Array.<!Representation>} */
	Representations []*Representation
}

/**
 * Parses an "AdaptationSet" tag.
 * @param {!Period} parent The parent Period.
 * @param {!Node} elem The AdaptationSet XML element.
 */
func (adaptationSet *AdaptationSet) Parse(parent Node, elem xml.Node) {
	var err error
	var p *Period = parent.(*Period)
	var contentComponent *ContentComponent
	var role *Role

	ok := false

	// Parse children which provide properties of the AdaptationSet.
	if contentComponent, ok = parseChild(adaptationSet, elem, ContentComponent_TAG_NAME).(*ContentComponent); ok == false {
		contentComponent = nil
	}

	if role, ok = parseChild(adaptationSet, elem, Role_TAG_NAME).(*Role); ok == false {
		role = nil
	}

	// Parse attributes.
	adaptationSet.Id, _ = parseAttrAsString(elem, "id")
	adaptationSet.Width, _ = parseAttrAsPositiveInt(elem, "width")
	adaptationSet.Height, _ = parseAttrAsPositiveInt(elem, "height")
	adaptationSet.MimeType, _ = parseAttrAsString(elem, "mimeType")
	adaptationSet.Codecs, _ = parseAttrAsString(elem, "codecs")

	if adaptationSet.Lang, err = parseAttrAsString(elem, "lang"); err != nil {
		if contentComponent != nil {
			adaptationSet.Lang = contentComponent.Lang
		} else {
			adaptationSet.Lang = ""
		}
	}

	if adaptationSet.ContentType, err = parseAttrAsString(elem, "contentType"); err != nil {
		if contentComponent != nil {
			adaptationSet.ContentType = contentComponent.ContentType
		} else {
			adaptationSet.ContentType = ""
		}
	}

	adaptationSet.Main = (role != nil && role.Value == "main")

	// Normalize the language tag.
	// TODO: normalize language.
	// if (this.lang) this.lang = shaka.util.LanguageUtils.normalize(this.lang);

	// Parse simple child elements.
	if adaptationSet.BaseUrl, ok = parseChild(adaptationSet, elem, BaseUrl_TAG_NAME).(*BaseUrl); ok == false {
		adaptationSet.BaseUrl = p.BaseUrl
	}

	// adaptationSet.ContentProtections = parseChildren(adaptationSet, elem, ContentProtection_TAG_NAME)

	if len(adaptationSet.ContentType) == 0 && len(adaptationSet.MimeType) != 0 {
		// Infer contentType from mimeType. This must be done before parsing any
		// child Representations, as Representation inherits contentType.
		adaptationSet.ContentType = strings.Split(adaptationSet.MimeType, "/")[0]
	}

	// Parse hierarchical children.
	if p.SegmentBase != nil {
		adaptationSet.SegmentBase = mergeChild(adaptationSet, elem, p.SegmentBase, SegmentBase_TAG_NAME).(*SegmentBase)
	} else {
		if adaptationSet.SegmentBase, ok = parseChild(adaptationSet, elem, SegmentBase_TAG_NAME).(*SegmentBase); ok == false {
			adaptationSet.SegmentBase = nil
		}
	}

	if p.SegmentList != nil {
		adaptationSet.SegmentList = mergeChild(adaptationSet, elem, p.SegmentList, SegmentList_TAG_NAME).(*SegmentList)
	} else {
		if adaptationSet.SegmentList, ok = parseChild(adaptationSet, elem, SegmentList_TAG_NAME).(*SegmentList); ok == false {
			adaptationSet.SegmentList = nil
		}
	}

	if p.SegmentTemplate != nil {
		adaptationSet.SegmentTemplate = mergeChild(adaptationSet, elem, p.SegmentTemplate, SegmentTemplate_TAG_NAME).(*SegmentTemplate)
	} else {
		if adaptationSet.SegmentTemplate, ok = parseChild(adaptationSet, elem, SegmentTemplate_TAG_NAME).(*SegmentTemplate); ok == false {
			adaptationSet.SegmentTemplate = nil
		}
	}

	children := parseChildren(adaptationSet, elem, Representation_TAG_NAME)
	adaptationSet.Representations = make([]*Representation, len(children))
	for i, child := range children {
		adaptationSet.Representations[i] = child.(*Representation)
	}

	if len(adaptationSet.MimeType) == 0 && len(adaptationSet.Representations) != 0 {
		// Infer mimeType from children.  MpdProcessor will deal with the case
		// where Representations have inconsistent mimeTypes.
		adaptationSet.MimeType = adaptationSet.Representations[0].MimeType

		if len(adaptationSet.ContentType) == 0 && len(adaptationSet.MimeType) != 0 {
			adaptationSet.ContentType = strings.Split(adaptationSet.MimeType, "/")[0]
		}
	}
}

func NewAdaptationSet() Node {
	return &AdaptationSet{}
}
