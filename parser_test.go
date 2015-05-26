package main

import (
	"github.com/moovweb/gokogiri/xml"
	"testing"
)

func TestParseAttAsDate(t *testing.T) {
	var sinceEpoch int64
	var err error

	xmlText := "<Root birthday=\"1984-10-21T05:00:00.000Z\"></Root>"
	doc, err := xml.Parse([]byte(xmlText), xml.DefaultEncodingBytes, nil, xml.DefaultParseOption, xml.DefaultEncodingBytes)

	if err != nil {
		t.Error(err)
	}

	if sinceEpoch, err = parseAttAsDate(doc.FirstChild(), "birthday"); err != nil {
		t.Error(err)
	}
	if sinceEpoch != 467182800 {
		t.Errorf("expecting time from epoch to be %d, got: %d", 467182800, sinceEpoch)
	}

	if sinceEpoch, err = parseAttAsDate(doc.FirstChild(), "MIA"); err == nil {
		t.Error("expecting to receive an error, got nil")
	}
	if sinceEpoch != 0 {
		t.Errorf("expecting sinceEpoch to be 0 got: %d", sinceEpoch)
	}
}

func TestParseAttrAsDuration(t *testing.T) {
	var seconds int
	var err error

	xmlText := "<Root duration=\"PT2H12M52S\"></Root>"
	doc, err := xml.Parse([]byte(xmlText), xml.DefaultEncodingBytes, nil, xml.DefaultParseOption, xml.DefaultEncodingBytes)

	if err != nil {
		t.Error(err)
	}

	if seconds, err = parseAttrAsDuration(doc.FirstChild(), "duration"); err != nil {
		t.Error(err)
	}
	if seconds != (2*60*60)+(12*60)+52 {
		t.Errorf("expecting seconds to be %d, got: %d", (12*60)+52, seconds)
	}

	if seconds, err = parseAttrAsDuration(doc.FirstChild(), "MIA"); err == nil {
		t.Error("expecting to receive an error, got nil")
	}
	if seconds != 0 {
		t.Errorf("expecting seconds to be 0 got: %d", seconds)
	}
}

func TestParseAttrAsRange(t *testing.T) {
	var r *Range
	var err error

	xmlText := "<Root acceleration=\"0-1000\"></Root>"
	doc, err := xml.Parse([]byte(xmlText), xml.DefaultEncodingBytes, nil, xml.DefaultParseOption, xml.DefaultEncodingBytes)

	if err != nil {
		t.Error(err)
	}

	if r, err = parseAttrAsRange(doc.FirstChild(), "acceleration"); err != nil {
		t.Error(err)
	}
	if r.Begin != 0 {
		t.Errorf("expecting range to start at 0, got %d", r.Begin)
	}
	if r.End != 1000 {
		t.Errorf("expecting range to End at 1000, got %d", r.End)
	}

	if r, err = parseAttrAsRange(doc.FirstChild(), "MIA"); err == nil {
		t.Error("expecting to receive an error, got nil")
	}
	if r != nil {
		t.Errorf("expecting range to be nil got: %s", r)
	}
}

func TestparseAttrAsPositiveInt(t *testing.T) {
	var num int
	var err error

	xmlText := "<Root meaning=\"42\" freez=\"−173\" void=\"0\"></Root>"
	doc, err := xml.Parse([]byte(xmlText), xml.DefaultEncodingBytes, nil, xml.DefaultParseOption, xml.DefaultEncodingBytes)

	if err != nil {
		t.Error(err)
	}

	if num, err = parseAttrAsPositiveInt(doc.FirstChild(), "meaning"); err != nil {
		t.Error(err)
	}
	if num != 42 {
		t.Errorf("expecting num to be 42, got: %d", num)
	}

	if num, err = parseAttrAsPositiveInt(doc.FirstChild(), "void"); err == nil {
		t.Error("expecting to receive an error, got nil")
	}
	if num != 0 {
		t.Errorf("expecting num to be 0, got: %d", num)
	}

	if num, err = parseAttrAsPositiveInt(doc.FirstChild(), "freez"); err == nil {
		t.Error("expecting to receive an error, got nil")
	}
	if num != 0 {
		t.Errorf("expecting num to be 0, got: %d", num)
	}
}

func TestParseAttrAsNonNegativeInt(t *testing.T) {
	var num int
	var err error

	xmlText := "<Root meaning=\"42\" freez=\"−173\" void=\"0\"></Root>"
	doc, err := xml.Parse([]byte(xmlText), xml.DefaultEncodingBytes, nil, xml.DefaultParseOption, xml.DefaultEncodingBytes)

	if err != nil {
		t.Error(err)
	}

	if num, err = parseAttrAsNonNegativeInt(doc.FirstChild(), "meaning"); err != nil {
		t.Error(err)
	}
	if num != 42 {
		t.Errorf("expecting num to be 42, got: %d", num)
	}

	if num, err = parseAttrAsNonNegativeInt(doc.FirstChild(), "void"); err != nil {
		t.Error(err)
	}
	if num != 0 {
		t.Errorf("expecting num to be 0, got: %d", num)
	}

	if num, err = parseAttrAsNonNegativeInt(doc.FirstChild(), "freez"); err == nil {
		t.Error("expecting to receive an error, got nil")
	}
	if num != 0 {
		t.Errorf("expecting num to be 0, got: %d", num)
	}
}

func TestParseAttrAsString(t *testing.T) {
	var str string
	var err error

	xmlText := "<Root monkey=\"business\"></Root>"
	doc, err := xml.Parse([]byte(xmlText), xml.DefaultEncodingBytes, nil, xml.DefaultParseOption, xml.DefaultEncodingBytes)

	if err != nil {
		t.Error(err)
	}

	if str, err = parseAttrAsString(doc.FirstChild(), "monkey"); err != nil {
		t.Error(err)
	}
	if str != "business" {
		t.Errorf("expecting string to be 'business', got: %s", str)
	}

	if str, err = parseAttrAsString(doc.FirstChild(), "monkey-bar"); err == nil {
		t.Errorf("error should have returned for missing attribure, got nil")
	}
	if str != "" {
		t.Error("string should return empty for missing attribute")
	}
}

func TestParser(t *testing.T) {
	var root *Mpd
	var err error

	if root, err = ParseMpd("http://sdk.streamrail.com/pepsi/cdn/0.0.1/3a5dd80efc3a867e55c69996c7f22051f6c3b94d/dash/manifest.mpd"); err != nil {
		t.Error(err)
	}

	// warning Long validation ahead

	if root == nil {
		t.Errorf("MPD root node is nil")
	}

	if root.Type != "static" {
		t.Error("expecting mpd type to be static, got %s", root.Type)
	}

	if len(root.Periods) != 1 {
		t.Errorf("expecting mpd to have one period, got %d", len(root.Periods))
	}

	if len(root.Periods[0].AdaptationSets) != 2 {
		t.Error("expecting mpd to have two adaptation sets, got %d", len(root.Periods[0].AdaptationSets))
	}

	if root.Periods[0].AdaptationSets[0].ContentType != "audio" {
		t.Errorf("expecting first adaptaion content type to be audio, actual %s", root.Periods[0].AdaptationSets[0].ContentType)
	}

	if root.Periods[0].AdaptationSets[1].ContentType != "video" {
		t.Errorf("expecting first adaptaion content type to be audio, actual %s", root.Periods[0].AdaptationSets[0].ContentType)
	}

	if root.Periods[0].AdaptationSets[0].SegmentTemplate == nil {
		t.Errorf("audio segment template should not be nil")
	}

	if root.Periods[0].AdaptationSets[0].SegmentTemplate.MediaUrlTemplate != "$RepresentationID$/seg-$Number$.m4f" {
		t.Errorf("expecting audio segment template media url template to be $RepresentationID$/seg-$Number$.m4f, actual %s", root.Periods[0].AdaptationSets[0].SegmentTemplate.MediaUrlTemplate)
	}

	if root.Periods[0].AdaptationSets[0].SegmentTemplate.InitializationUrlTemplate != "$RepresentationID$/init.mp4" {
		t.Errorf("expecting audio segment template init template to be $RepresentationID$/init.mp4, actual %s", root.Periods[0].AdaptationSets[0].SegmentTemplate.InitializationUrlTemplate)
	}

	if root.Periods[0].AdaptationSets[0].Representations == nil {
		t.Errorf("represenration of first adaptation set is nil")
	}

	if len(root.Periods[0].AdaptationSets[0].Representations) != 1 {
		t.Errorf("expecting first adaption set to have only one representation, got %d", len(root.Periods[0].AdaptationSets[0].Representations))
	}

	if root.Periods[0].AdaptationSets[0].Representations[0].MimeType != "audio/mp4" {
		t.Errorf("expecting first representation mime type to be audio/mp4 got %s", root.Periods[0].AdaptationSets[0].Representations[0].MimeType)
	}

	if root.Periods[0].AdaptationSets[0].Representations[0].SegmentTemplate == nil {
		t.Errorf("expecting first representation to have a segment template, got nil")
	}

	if root.Periods[0].AdaptationSets[0].Representations[0].SegmentTemplate.MediaUrlTemplate != "$RepresentationID$/seg-$Number$.m4f" {
		t.Errorf("expecting first representation segment tempate media url template to be $RepresentationID$/seg-$Number$.m4f got %s", root.Periods[0].AdaptationSets[0].Representations[0].SegmentTemplate.MediaUrlTemplate)
	}

	if root.Periods[0].AdaptationSets[0].Representations[0].SegmentTemplate.InitializationUrlTemplate != "$RepresentationID$/init.mp4" {
		t.Errorf("expecting first representation segment tempate initialization url template to be $RepresentationID$/init.mp4 got %s", root.Periods[0].AdaptationSets[0].Representations[0].SegmentTemplate.InitializationUrlTemplate)
	}

	//-------------------------------------------------------------------------------------------------

	if root.Periods[0].AdaptationSets[1].SegmentTemplate == nil {
		t.Errorf("audio segment template should not be nil")
	}

	if root.Periods[0].AdaptationSets[1].SegmentTemplate.MediaUrlTemplate != "$RepresentationID$/seg-$Number$.m4f" {
		t.Errorf("expecting audio segment template media url template to be $RepresentationID$/seg-$Number$.m4f, actual %s", root.Periods[0].AdaptationSets[0].SegmentTemplate.MediaUrlTemplate)
	}

	if root.Periods[0].AdaptationSets[1].SegmentTemplate.InitializationUrlTemplate != "$RepresentationID$/init.mp4" {
		t.Errorf("expecting audio segment template init template to be $RepresentationID$/init.mp4, actual %s", root.Periods[0].AdaptationSets[0].SegmentTemplate.InitializationUrlTemplate)
	}

	if root.Periods[0].AdaptationSets[1].Representations == nil {
		t.Errorf("represenration of first adaptation set is nil")
	}

	if len(root.Periods[0].AdaptationSets[1].Representations) != 1 {
		t.Errorf("expecting first adaption set to have only one representation, got %d", len(root.Periods[0].AdaptationSets[0].Representations))
	}

	if root.Periods[0].AdaptationSets[1].Representations[0].MimeType != "video/mp4" {
		t.Errorf("expecting first representation mime type to be video/mp4 got %s", root.Periods[0].AdaptationSets[0].Representations[0].MimeType)
	}

	if root.Periods[0].AdaptationSets[1].Representations[0].SegmentTemplate == nil {
		t.Errorf("expecting first representation to have a segment template, got nil")
	}

	if root.Periods[0].AdaptationSets[1].Representations[0].SegmentTemplate.MediaUrlTemplate != "$RepresentationID$/seg-$Number$.m4f" {
		t.Errorf("expecting first representation segment tempate media url template to be $RepresentationID$/seg-$Number$.m4f got %s", root.Periods[0].AdaptationSets[0].Representations[0].SegmentTemplate.MediaUrlTemplate)
	}

	if root.Periods[0].AdaptationSets[1].Representations[0].SegmentTemplate.InitializationUrlTemplate != "$RepresentationID$/init.mp4" {
		t.Errorf("expecting first representation segment tempate initialization url template to be $RepresentationID$/init.mp4 got %s", root.Periods[0].AdaptationSets[0].Representations[0].SegmentTemplate.InitializationUrlTemplate)
	}
}
