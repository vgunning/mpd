package main

import (
	"errors"
	"fmt"
	"github.com/moovweb/gokogiri"
	"github.com/moovweb/gokogiri/xml"
	"io/ioutil"
	"reflect"
	"regexp"
	"strconv"
	"time"
)

const (
	DEFAULT_MIN_BUFFER_TIME_              = 5
	DEFAULT_SUGGESTED_PRESENTATION_DELAY_ = 1
)

// MPD tag names --------------------------------------------------------------
const (

	/**
	 * @expose all TAG_NAME properties so that they do not get stripped during
	 *     advanced compilation.
	 */
	Mpd_TAG_NAME = "MPD"

	Period_TAG_NAME = "Period"

	AdaptationSet_TAG_NAME = "AdaptationSet"

	Role_TAG_NAME = "Role"

	ContentComponent_TAG_NAME = "ContentComponent"

	Representation_TAG_NAME = "Representation"

	ContentProtection_TAG_NAME = "ContentProtection"

	CencPssh_TAG_NAME = "cenc:pssh"

	BaseUrl_TAG_NAME = "BaseURL"

	SegmentBase_TAG_NAME = "SegmentBase"

	RepresentationIndex_TAG_NAME = "RepresentationIndex"

	Initialization_TAG_NAME = "Initialization"

	SegmentList_TAG_NAME = "SegmentList"

	SegmentUrl_TAG_NAME = "SegmentURL"

	SegmentTemplate_TAG_NAME = "SegmentTemplate"

	SegmentTimeline_TAG_NAME = "SegmentTimeline"

	SegmentTimePoint_TAG_NAME = "S"
)

type Node interface {
	Parse(parent Node, elem xml.Node)
}

type Cloneable interface {
	Clone() Node
}

func ParseMpd(source, url string) (*Mpd, error) {
	initTypeRegistry()

	// read xml file
	content, err := ioutil.ReadFile(source)
	if err != nil {
		fmt.Printf("failed to open example file, error: %s\n", err)
		return nil, err
	}

	// parse file
	doc, err := gokogiri.ParseXml(content)

	if err != nil {
		fmt.Printf("failed to parse xml file, error: %s\n", err)
		return nil, err
	}

	// important -- don't forget to free the resources when you're done!
	defer doc.Free()

	// Construct a virtual parent for the MPD to use in resolving relative URLs.
	parent := FakeNode{BaseUrl: &BaseUrl{Url: url}}

	fmt.Println("Start parsing")
	ok := false
	root, ok := parseChild(parent, doc, Mpd_TAG_NAME).(*Mpd)

	if ok {
		PrintMPD(root, 0)
		return root, nil
	} else {
		return nil, errors.New("failed to parse mpd")
	}
}

func PrintMPD(root Node, ident int) {

	// Check for zero value
	v := reflect.ValueOf(root)
	if v.Interface() == reflect.Zero(v.Type()).Interface() {
		return
	}

	s := reflect.ValueOf(root).Elem()
	typeOfT := s.Type()

	printTabs(ident)
	fmt.Println(typeOfT)

	// Scan type fields
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)

		printTabs(ident + 1)
		// fmt.Printf("%d: %s %s = %v\n", i, typeOfT.Field(i).Name, f.Type(), f.Interface())
		fmt.Printf("%s = %v\n", typeOfT.Field(i).Name, f.Interface())

		ok := false
		var n Node

		// Is this a slice?
		switch f.Kind() {
		case reflect.Slice:
			for i := 0; i < f.Len(); i++ {
				e := f.Index(i)
				if n, ok = e.Interface().(Node); ok == true {
					PrintMPD(n, ident+1)
				}
			}
		}

		if n, ok = f.Interface().(Node); ok == true {
			PrintMPD(n, ident+1)
		}
	}
}

func printTabs(ammount int) {
	for i := 0; i < ammount; i++ {
		fmt.Printf("\t")
	}
}

type constructor func() Node

// var typeRegistry = make(map[string]reflect.Type)
var typeRegistry = make(map[string]constructor)

func initTypeRegistry() {
	typeRegistry[Mpd_TAG_NAME] = NewMpd

	typeRegistry[Period_TAG_NAME] = NewPeriod

	typeRegistry[AdaptationSet_TAG_NAME] = NewAdaptationSet

	typeRegistry[Role_TAG_NAME] = NewRole

	typeRegistry[ContentComponent_TAG_NAME] = NewContentComponent

	typeRegistry[Representation_TAG_NAME] = NewRepresentation

	// typeRegistry[ContentProtection_TAG_NAME] = NewContentProtection

	// typeRegistry[CencPssh_TAG_NAME] = NewCencPssh

	typeRegistry[BaseUrl_TAG_NAME] = NewBaseUrl

	typeRegistry[SegmentBase_TAG_NAME] = NewSegmentBase

	typeRegistry[RepresentationIndex_TAG_NAME] = NewRepresentationIndex

	typeRegistry[Initialization_TAG_NAME] = NewInitialization

	typeRegistry[SegmentList_TAG_NAME] = NewSegmentList

	typeRegistry[SegmentUrl_TAG_NAME] = NewSegmentUrl

	typeRegistry[SegmentTemplate_TAG_NAME] = NewSegmentTemplate

	typeRegistry[SegmentTimeline_TAG_NAME] = NewSegmentTimeline

	typeRegistry[SegmentTimePoint_TAG_NAME] = NewSegmentTimePoint
}

func createInstance(name string) Node {
	if _, ok := typeRegistry[name]; !ok {
		fmt.Println("name missing from typeRegistry")
		return nil
	}

	v := typeRegistry[name]()
	return v
}

/**
 * A Range.
 * @param {number} begin The beginning of the range.
 * @param {number} end The end of the range.
 */
type Range struct {
	/** @const {number} */
	Begin int
	/** @const {number} */
	End int
}

func newRange(begin, end int) *Range {
	return &Range{
		Begin: begin,
		End:   end,
	}
}

/**
 * Creates a deep copy of this Range.
 * @return {Range}
 */
func (r Range) Clone() *Range {
	return newRange(r.Begin, r.End)
}

/**
 * @param {T} obj
 * @return {T} A clone of |obj| if |obj| is non-null; otherwise, return null.
 * @private
 * @template T
 */
func clone(obj Cloneable) Node {
	return obj.Clone()
}

/**
 * Gets the text contents of a node.
 * @param {!Node} elem The XML element.
 * @return {?string} The text contents, or null if there are none.
 * @private
 */
func getContents(elem xml.Node) (string, error) {
	if elem.FirstChild().NodeType() != xml.XML_TEXT_NODE {
		return "", errors.New("wrong node type")
	} else {
		return elem.FirstChild().Content(), nil
	}
}

/**
 * Parses a child XML element and merges it into an existing MPD node object.
 * @param {*} parent The parent MPD node object.
 * @param {!Node} elem The parent XML element.
 * @param {!T} original The existing MPD node object.
 * @param {!String} original's tag name.
 * @return {!T} The merged MPD node object. If a child XML element cannot be
 *     parsed (see parseChild_) then the merged MPD node object is identical
 *     to |original|, although it is not the same object.
 * @template T
 * @private
 */
func mergeChild(parent Node, elem xml.Node, original Cloneable, originalTagName string) Node {
	// fmt.Printf("mergeChild parent: %s, elem: %s, original: %s, originalTagName: %s\r\n", parent, elem, original, originalTagName)
	merged := original.Clone()

	if childElement, err := findChild(elem, originalTagName); err == nil {
		merged.Parse(parent, childElement)
	}

	return merged
}

/**
 * Parses a child XML element.
 * @param {*} parent The parent MPD node object.
 * @param {!Node} elem The parent XML element.
 * @param {function(new:T)} constructor The constructor of the parsed
 *     child XML element. The constructor must define the attribute "TAG_NAME".
 * @return {T} The parsed child XML element on success, or null if a child
 *     XML element does not exist with the given tag name OR if there exists
 *     more than one child XML element with the given tag name OR if the child
 *     XML element could not be parsed.
 * @template T
 * @private
 */
func parseChild(parent Node, elem xml.Node, name string) Node {
	var parsedChild Node
	var childElement xml.Node
	var err error

	if childElement, err = findChild(elem, name); err != nil {
		fmt.Println("Element missing")
		return parsedChild
	}

	parsedChild = createInstance(name)
	parsedChild.Parse(parent, childElement)

	return parsedChild
}

func findChild(elem xml.Node, name string) (xml.Node, error) {
	var childElement xml.Node
	found := false

	for child := elem.FirstChild(); child != nil; child = child.NextSibling() {
		if child.Name() != name {
			continue
		}

		if found == true {
			return childElement, errors.New("more than one child with given tag name exists")
		}

		found = true
		childElement = child
	}

	if found == true {
		return childElement, nil
	} else {
		return childElement, errors.New("child with given tag name is missing")
	}
}

/**
 * Parses an array of child XML elements.
 * @param {*} parent The parsed parent object.
 * @param {!Node} elem The parent XML element.
 * @param {function(new:T)} constructor The constructor of each parsed child
 *     XML element. The constructor must define the attribute "TAG_NAME".
 * @return {!Array.<!T>} The parsed child XML elements.
 * @template T
 * @private
 */
func parseChildren(parent Node, elem xml.Node, name string) []Node {
	var parsedChildren []Node

	for childNode := elem.FirstChild(); childNode != nil; childNode = childNode.NextSibling() {
		if childNode.Name() != name {
			continue
		}

		if parsedChild := createInstance(name); parsedChild != nil {
			parsedChild.Parse(parent, childNode)
			parsedChildren = append(parsedChildren, parsedChild)
		}
	}

	return parsedChildren
}

/**
 * Parses an XML date string.
 * @param {string} dateString
 * @return {?number}
 * @private
 */
func parseAttAsDate(elem xml.Node, name string) (int64, error) {
	attribute := elem.Attribute(name)
	if attribute == nil {
		return 0, errors.New("missing attribute")
	}

	layout := "2006-01-02T15:04:05.000Z"
	if t, err := time.Parse(layout, attribute.Value()); err != nil {
		return 0, err
	} else {
		return t.Unix(), nil
	}
}

/**
 * Parses an XML duration string.
 * Negative values are not supported. Years and months are treated as exactly
 * 365 and 30 days respectively.
 * @param {string} durationString The duration string, e.g., "PT1H3M43.2S",
 *     which means 1 hour, 3 minutes, and 43.2 seconds.
 * @return {?number} The parsed duration in seconds, or null if the duration
 *     string could not be parsed.
 * @see http://www.datypic.com/sc/xsd/t-xsd_duration.html
 * @private
 */
func parseAttrAsDuration(elem xml.Node, name string) (int, error) {
	attribute := elem.Attribute(name)
	if attribute == nil {
		return 0, errors.New("missing attribute")
	}

	re := regexp.MustCompile("^P(?:([0-9]*)Y)?(?:([0-9]*)M)?(?:([0-9]*)D)?(?:T(?:([0-9]*)H)?(?:([0-9]*)M)?(?:([0-9.]*)S)?)?$")
	matches := re.FindStringSubmatch(attribute.Value())

	if matches == nil {
		return 0, errors.New("attribute is not a duration")
	}

	duration := 0

	// Assume a year always has 365 days.
	if len(matches[1]) > 0 {
		if years, err := parseNonNegativeInt(matches[1]); err != nil {
			return 0, err
		} else {
			duration += (60 * 60 * 24 * 365) * years
		}
	}

	// Assume a month is 30 days.
	if len(matches[2]) > 0 {
		if months, err := parseNonNegativeInt(matches[2]); err != nil {
			return 0, err
		} else {
			duration += (60 * 60 * 24 * 30) * months
		}
	}

	if len(matches[3]) > 0 {
		if days, err := parseNonNegativeInt(matches[3]); err != nil {
			return 0, err
		} else {
			duration += (60 * 60 * 24) * days
		}
	}

	if len(matches[4]) > 0 {
		if hours, err := parseNonNegativeInt(matches[4]); err != nil {
			return 0, err
		} else {
			duration += (60 * 60) * hours
		}
	}

	if len(matches[5]) > 0 {
		if minutes, err := parseNonNegativeInt(matches[5]); err != nil {
			return 0, err
		} else {
			duration += 60 * minutes
		}
	}

	if len(matches[6]) > 0 {
		if seconds, err := strconv.ParseInt(matches[6], 10, 32); err != nil {
			return 0, err
		} else {
			duration += int(seconds)
		}
	}

	return duration, nil
}

/**
 * Parses a range string.
 * @param {string} rangeString The range string, e.g., "101-9213"
 * @return {Range} The parsed range, or null if the range string
 *     could not be parsed.
 * @private
 */
func parseAttrAsRange(elem xml.Node, name string) (*Range, error) {
	var err error
	attribute := elem.Attribute(name)
	begin := 0
	end := 0

	if attribute == nil {
		return nil, errors.New("missing attribute")
	}

	valueStr := attribute.Value()
	re := regexp.MustCompile("([0-9]+)-([0-9]+)")
	matches := re.FindStringSubmatch(valueStr)

	if matches == nil {
		return nil, errors.New("attribute is not a range")
	}

	if begin, err = parseNonNegativeInt(matches[1]); err != nil {
		return nil, errors.New("invalid range begin value")
	}

	if end, err = parseNonNegativeInt(matches[2]); err != nil {
		return nil, errors.New("invalid range end value")
	}

	return newRange(begin, end), nil
}

/**
 * Parses a positive integer.
 * @param {string} intString The integer string.
 * @return {?number} The parsed positive integer on success; otherwise,
 *     return null.
 * @private
 */
func parseAttrAsPositiveInt(elem xml.Node, name string) (int, error) {
	attribute := elem.Attribute(name)
	if attribute == nil {
		return 0, errors.New("missing attribute")
	}
	return parsePositiveInt(attribute.Value())
}

func parsePositiveInt(value string) (int, error) {
	if value, err := strconv.ParseInt(value, 10, 32); err != nil {
		return 0, err
	} else {
		if value <= 0 {
			return 0, errors.New("attribute is not positive")
		} else {
			return int(value), nil
		}
	}
}

/**
 * Parses a non-negative integer.
 * @param {string} intString The integer string.
 * @return {?number} The parsed non-negative integer on success; otherwise,
 *     return null.
 * @private
 */
func parseAttrAsNonNegativeInt(elem xml.Node, name string) (int, error) {
	attribute := elem.Attribute(name)
	if attribute == nil {
		return 0, errors.New("missing attribute")
	}
	return parseNonNegativeInt(attribute.Value())
}

func parseNonNegativeInt(value string) (int, error) {
	if value, err := strconv.ParseInt(value, 10, 32); err != nil {
		return 0, err
	} else {
		if value < 0 {
			return 0, errors.New("attribute is negative")
		} else {
			return int(value), nil
		}
	}
}

/**
 * A misnomer.  Does no parsing, just returns the input string as-is.
 * @param {string} inputString The inputString.
 * @return {?string} The "parsed" string.  The type is specified as nullable
 *     only to fit into the parseAttr_() template, but null will never be
 *     returned.
 * @private
 */
func parseAttrAsString(elem xml.Node, name string) (string, error) {
	attribute := elem.Attribute(name)
	if attribute == nil {
		return "", errors.New("missing attribute")
	}
	return attribute.Value(), nil
}
