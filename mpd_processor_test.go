package main

import (
	"fmt"
	"testing"
)

func TestMPDProcessingExample2(t *testing.T) {
	var mpd *Mpd
	var err error

	if mpd, err = ParseMpd("example2.xml", "streamrail.com/"); err != nil {
		t.Error(err)
	}

	mpdProcessor := NewMpdProcessor()
	mpdProcessor.Process(*mpd)

	if len(mpdProcessor.ManifestInfo.PeriodInfos) != 1 {
		t.Errorf("expecting 2 periods, got %d", len(mpdProcessor.ManifestInfo.PeriodInfos))
	}

	if len(mpdProcessor.ManifestInfo.PeriodInfos[0].StreamSetInfos) != 2 {
		t.Errorf("expecting 2 strean sets, got %d", len(mpdProcessor.ManifestInfo.PeriodInfos[0].StreamSetInfos))
	}

	if len(mpdProcessor.ManifestInfo.PeriodInfos[0].StreamSetInfos) != 2 {
		t.Errorf("expecting 2 strean sets, got %d", len(mpdProcessor.ManifestInfo.PeriodInfos[0].StreamSetInfos))
	}

	if len(mpdProcessor.ManifestInfo.PeriodInfos[0].StreamSetInfos[0].StreamInfos) != 1 {
		t.Errorf("expecting 1 strean info, got %d", len(mpdProcessor.ManifestInfo.PeriodInfos[0].StreamSetInfos[0].StreamInfos))
	}

	if len(mpdProcessor.ManifestInfo.PeriodInfos[0].StreamSetInfos[1].StreamInfos) != 1 {
		t.Errorf("expecting 1 strean info, got %d", len(mpdProcessor.ManifestInfo.PeriodInfos[0].StreamSetInfos[1].StreamInfos))
	}

	if mpdProcessor.ManifestInfo.PeriodInfos[0].StreamSetInfos[0].StreamInfos[0].SegmentIndex == nil {
		t.Error("missing segment index")
	}

	if mpdProcessor.ManifestInfo.PeriodInfos[0].StreamSetInfos[1].StreamInfos[0].SegmentIndex == nil {
		t.Error("missing segment index")
	}

	if len(mpdProcessor.ManifestInfo.PeriodInfos[0].StreamSetInfos[0].StreamInfos[0].SegmentIndex.References) != 386 {
		t.Errorf("expecting 386 references, got %d", len(mpdProcessor.ManifestInfo.PeriodInfos[0].StreamSetInfos[0].StreamInfos[0].SegmentIndex.References))
	}

	if len(mpdProcessor.ManifestInfo.PeriodInfos[0].StreamSetInfos[1].StreamInfos[0].SegmentIndex.References) != 386 {
		t.Errorf("expecting 386 references, got %d", len(mpdProcessor.ManifestInfo.PeriodInfos[0].StreamSetInfos[1].StreamInfos[0].SegmentIndex.References))
	}

	if mpdProcessor.ManifestInfo.PeriodInfos[0].StreamSetInfos[0].StreamInfos[0].SegmentIndex.References[385].Url != "streamrail.com/302k/audio/und/seg-386.m4f" {
		t.Errorf("expecting last reference to point to %s, got:%s", "streamrail.com/302k/audio/und/seg-386.m4f", mpdProcessor.ManifestInfo.PeriodInfos[0].StreamSetInfos[0].StreamInfos[0].SegmentIndex.References[365].Url)
	}

	if mpdProcessor.ManifestInfo.PeriodInfos[0].StreamSetInfos[1].StreamInfos[0].SegmentIndex.References[385].Url != "streamrail.com/302k/video/1/seg-386.m4f" {
		t.Errorf("expecting last reference to point to %s, got:%s", "streamrail.com/302k/video/1/seg-386.m4f", mpdProcessor.ManifestInfo.PeriodInfos[0].StreamSetInfos[1].StreamInfos[0].SegmentIndex.References[365].Url)
	}
}

func TestMPDProcessingExample3(t *testing.T) {
	var mpd *Mpd
	var err error

	if mpd, err = ParseMpd("example3.xml", "streamrail.com/"); err != nil {
		t.Error(err)
	}

	mpdProcessor := NewMpdProcessor()
	mpdProcessor.Process(*mpd)

	if len(mpdProcessor.ManifestInfo.PeriodInfos) != 1 {
		t.Errorf("expecting 2 periods, got %d", len(mpdProcessor.ManifestInfo.PeriodInfos))
	}

	if len(mpdProcessor.ManifestInfo.PeriodInfos[0].StreamSetInfos) != 2 {
		t.Errorf("expecting 2 strean sets, got %d", len(mpdProcessor.ManifestInfo.PeriodInfos[0].StreamSetInfos))
	}

	if len(mpdProcessor.ManifestInfo.PeriodInfos[0].StreamSetInfos) != 2 {
		t.Errorf("expecting 2 strean sets, got %d", len(mpdProcessor.ManifestInfo.PeriodInfos[0].StreamSetInfos))
	}

	if len(mpdProcessor.ManifestInfo.PeriodInfos[0].StreamSetInfos[0].StreamInfos) != 1 {
		t.Errorf("expecting 1 strean info, got %d", len(mpdProcessor.ManifestInfo.PeriodInfos[0].StreamSetInfos[0].StreamInfos))
	}

	if len(mpdProcessor.ManifestInfo.PeriodInfos[0].StreamSetInfos[1].StreamInfos) != 3 {
		t.Errorf("expecting 3 strean infos, got %d", len(mpdProcessor.ManifestInfo.PeriodInfos[0].StreamSetInfos[1].StreamInfos))
	}

	if mpdProcessor.ManifestInfo.PeriodInfos[0].StreamSetInfos[0].StreamInfos[0].SegmentIndex == nil {
		t.Error("missing segment index")
	}

	if mpdProcessor.ManifestInfo.PeriodInfos[0].StreamSetInfos[1].StreamInfos[0].SegmentIndex == nil {
		t.Error("missing segment index")
	}

	if len(mpdProcessor.ManifestInfo.PeriodInfos[0].StreamSetInfos[0].StreamInfos[0].SegmentIndex.References) != 137 {
		t.Errorf("expecting 137 references, got %d", len(mpdProcessor.ManifestInfo.PeriodInfos[0].StreamSetInfos[0].StreamInfos[0].SegmentIndex.References))
	}

	if len(mpdProcessor.ManifestInfo.PeriodInfos[0].StreamSetInfos[1].StreamInfos[0].SegmentIndex.References) != 137 {
		t.Errorf("expecting 137 references, got %d", len(mpdProcessor.ManifestInfo.PeriodInfos[0].StreamSetInfos[1].StreamInfos[0].SegmentIndex.References))
	}

	if len(mpdProcessor.ManifestInfo.PeriodInfos[0].StreamSetInfos[1].StreamInfos[1].SegmentIndex.References) != 137 {
		t.Errorf("expecting 137 references, got %d", len(mpdProcessor.ManifestInfo.PeriodInfos[0].StreamSetInfos[1].StreamInfos[0].SegmentIndex.References))
	}

	if len(mpdProcessor.ManifestInfo.PeriodInfos[0].StreamSetInfos[1].StreamInfos[2].SegmentIndex.References) != 137 {
		t.Errorf("expecting 137 references, got %d", len(mpdProcessor.ManifestInfo.PeriodInfos[0].StreamSetInfos[1].StreamInfos[0].SegmentIndex.References))
	}

	for _, ref := range mpdProcessor.ManifestInfo.PeriodInfos[0].StreamSetInfos[0].StreamInfos[0].SegmentIndex.References {
		fmt.Printf("Url: %s\r\n", ref.Url)
	}

	if mpdProcessor.ManifestInfo.PeriodInfos[0].StreamSetInfos[0].StreamInfos[0].SegmentIndex.References[136].Url != "streamrail.com/700k/audio/und/seg-137.m4f" {
		t.Errorf("expecting last reference to point to %s, got:%s", "streamrail.com/700k/audio/und/seg-137.m4f", mpdProcessor.ManifestInfo.PeriodInfos[0].StreamSetInfos[0].StreamInfos[0].SegmentIndex.References[365].Url)
	}

	if mpdProcessor.ManifestInfo.PeriodInfos[0].StreamSetInfos[1].StreamInfos[0].SegmentIndex.References[136].Url != "streamrail.com/700k/video/1/seg-137.m4f" {
		t.Errorf("expecting last reference to point to %s, got:%s", "streamrail.com/700k/video/1/seg-137.m4f", mpdProcessor.ManifestInfo.PeriodInfos[0].StreamSetInfos[1].StreamInfos[0].SegmentIndex.References[365].Url)
	}

	if mpdProcessor.ManifestInfo.PeriodInfos[0].StreamSetInfos[1].StreamInfos[1].SegmentIndex.References[136].Url != "streamrail.com/1200k/video/1/seg-137.m4f" {
		t.Errorf("expecting last reference to point to %s, got:%s", "streamrail.com/1200k/video/1/seg-137.m4f", mpdProcessor.ManifestInfo.PeriodInfos[0].StreamSetInfos[1].StreamInfos[0].SegmentIndex.References[365].Url)
	}

	if mpdProcessor.ManifestInfo.PeriodInfos[0].StreamSetInfos[1].StreamInfos[2].SegmentIndex.References[136].Url != "streamrail.com/1531k/video/1/seg-137.m4f" {
		t.Errorf("expecting last reference to point to %s, got:%s", "streamrail.com/1531k/video/1/seg-137.m4f", mpdProcessor.ManifestInfo.PeriodInfos[0].StreamSetInfos[1].StreamInfos[0].SegmentIndex.References[365].Url)
	}
}
