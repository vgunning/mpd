package main

import (
// "github.com/moovweb/gokogiri/xml"
)

type ContentProtection struct {
	/**
	 * @type {?string}
	 * @expose
	 */
	SchemeIdUri string

	/**
	 * @type {?string}
	 * @expose
	 */
	Value string

	/**
	 * @type {!Array.<!Node>}
	 * @expose
	 */
	Children []Node

	/**
	 * @type {CencPssh}
	 * @expose
	 */
	Pssh CencPssh
}

/**
 * Parses a "ContentProtection" tag.
 * @param {*} parent The parent object.
 * @param {!Node} elem The ContentProtection XML element.
 */
// func (contentProtection *ContentProtection) Parse(parent Node, elem xml.Node) {

// 	// Parse attributes.
// 	contentProtection.SchemeIdUri = parseAttrAsString(elem, "schemeIdUri")
// 	contentProtection.Value = parseAttrAsString(elem, "value")

// 	// Parse simple child elements.
// 	// TODO: Parse CencPssh
// 	// contentProtection.Pssh = parseChild(contentProtection, elem, CencPssh_TAG_NAME)

// 	// NOTE: A given ContentProtection tag could contain anything, and a scheme
// 	// could be application-specific.  Therefore we must capture whatever it
// 	// contains, and let the application choose a scheme and map it to a key
// 	// system.

// 	// TODO: Address line below.
// 	// contentProtection.Children = Array.prototype.slice.call(elem.childNodes);
// }
