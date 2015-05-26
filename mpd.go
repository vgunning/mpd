package mpd

import (
	"github.com/moovweb/gokogiri/xml"
)

type Mpd struct {
	/** @type {?string} */
	Id string

	/** @type {string} */
	Type string

	/** @type {*BaseUrl} */
	BaseUrl *BaseUrl

	/**
	 * The entire stream's duration, in seconds.
	 * @type {?number}
	 */
	MediaPresentationDuration int

	/**
	 * The quantity of media, in terms of seconds, that should be buffered before
	 * playback begins, to ensure uninterrupted playback.
	 * @type {?number}
	 */
	MinBufferTime int

	/**
	 * The interval, in seconds, to poll the media server for an updated
	 * MPD, or null if updates are not required.
	 * @type {?number}
	 */
	MinUpdatePeriod int

	/**
	 * The wall-clock time, in seconds, that the media content specified within
	 * the MPD started/will start to stream.
	 * @type {?number}
	 */
	AvailabilityStartTime int64

	/**
	 * The duration, in seconds, that the media server retains live media
	 * content, excluding the current segment and the previous segment, which are
	 * always available. For example, if this value is 60 then only media content
	 * up to 60 seconds from the beginning of the previous segment may be
	 * requested from the media server.
	 * @type {?number}
	 */
	TimeShiftBufferDepth int

	/**
	 * The duration, in seconds, that the media server takes to make live media
	 * content available. For example, if this value is 30 then only media
	 * content at least 30 seconds in the past may be requested from the media
	 * server.
	 * @type {?number}
	 */

	//DEFAULT_SUGGESTED_PRESENTATION_DELAY_;
	SuggestedPresentationDelay int

	/** @type {!Array.<!Period>} */
	Periods []*Period
}

func NewMpd() Node {
	return &Mpd{}
}

// MPD parsing functions ------------------------------------------------------

/**
 * Parses an "MPD" tag.
 * @param {!Object} parent A virtual parent tag containing a BaseURL which
 *     refers to the MPD resource itself.
 * @param {!Node} elem The MPD XML element.
 */
func (mpd *Mpd) Parse(parent Node, elem xml.Node) {
	var err error
	p := parent.(FakeNode)

	// Parse attributes.
	if mpd.Id, err = parseAttrAsString(elem, "id"); err != nil {

	}

	if mpd.Type, err = parseAttrAsString(elem, "type"); err != nil {
		mpd.Type = "static"
	}

	if mpd.MediaPresentationDuration, err = parseAttrAsDuration(elem, "mediaPresentationDuration"); err != nil {
		mpd.MediaPresentationDuration = -1
	}

	if mpd.MinBufferTime, err = parseAttrAsDuration(elem, "minBufferTime"); err != nil {
		mpd.MinBufferTime = DEFAULT_MIN_BUFFER_TIME_
	}

	if mpd.MinUpdatePeriod, err = parseAttrAsDuration(elem, "minimumUpdatePeriod"); err != nil {
		mpd.MinUpdatePeriod = 0
	}

	if mpd.AvailabilityStartTime, err = parseAttAsDate(elem, "availabilityStartTime"); err != nil {
		mpd.AvailabilityStartTime = -1
	}

	if mpd.TimeShiftBufferDepth, err = parseAttrAsDuration(elem, "timeShiftBufferDepth"); err != nil {
		mpd.TimeShiftBufferDepth = 0
	}

	if mpd.SuggestedPresentationDelay, err = parseAttrAsDuration(elem, "suggestedPresentationDelay"); err != nil {
		mpd.SuggestedPresentationDelay = DEFAULT_SUGGESTED_PRESENTATION_DELAY_
	}

	// Parse simple child elements.
	ok := false
	if mpd.BaseUrl, ok = parseChild(mpd, elem, BaseUrl_TAG_NAME).(*BaseUrl); ok == false {
		mpd.BaseUrl = p.BaseUrl
	}

	// Parse hierarchical children.
	children := parseChildren(mpd, elem, Period_TAG_NAME)
	mpd.Periods = make([]*Period, len(children))
	for i, child := range children {
		mpd.Periods[i] = child.(*Period)
	}
}
