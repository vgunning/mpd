package main

import (
	"github.com/moovweb/gokogiri/xml"
)

type ContentComponent struct {
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
	 * @type {?string}
	 */
	ContentType string
}

/**
 * Parses a "ContentComponent" tag.
 * @param {!AdaptationSet} parent The parent AdaptationSet.
 * @param {!Node} elem The ContentComponent XML element.
 */
func (contentComponent *ContentComponent) Parse(parent Node, elem xml.Node) {

	// Parse attributes.
	contentComponent.Id, _ = parseAttrAsString(elem, "id")
	contentComponent.Lang, _ = parseAttrAsString(elem, "lang")
	contentComponent.ContentType, _ = parseAttrAsString(elem, "contentType")

	// Normalize the language tag.
	// TODO: Normalize the language tag.
	// if (this.lang) this.lang = shaka.util.LanguageUtils.normalize(this.lang);
}

func NewContentComponent() Node {
	return &ContentComponent{}
}
