package mpd

import (
	"github.com/moovweb/gokogiri/xml"
)

type Role struct {
	/** @type {?string} */
	Value string
}

/**
 * Parses a "Role" tag.
 * @param {!AdaptationSet} parent The parent AdaptationSet.
 * @param {!Node} elem The Role XML element.
 */
func (role *Role) Parse(parent Node, elem xml.Node) {
	// Parse attributes.
	role.Value, _ = parseAttrAsString(elem, "value")
}

func NewRole() Node {
	return &Role{}
}
