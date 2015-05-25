package main

import (
	"github.com/moovweb/gokogiri/xml"
)

type BaseUrl struct {
	/** @type {?string} */
	Url string
}

func NewBaseUrl() Node {
	return &BaseUrl{}
}

/**
 * Parses a "BaseURL" tag.
 * @param {*} parent The parent object.
 * @param {!Node} elem The BaseURL XML element.
 */
func (baseUrl *BaseUrl) Parse(parent Node, elem xml.Node) {
	baseUrl.Url, _ = getContents(elem)
}
