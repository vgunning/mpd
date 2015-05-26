package mpd

import (
	"fmt"
	"testing"
)

func TestMPDProcessingExample1(t *testing.T) {
	var mpd *Mpd
	var err error

	if mpd, err = ParseMpd("http://sdk.streamrail.com/pepsi/cdn/0.0.1/3a5dd80efc3a867e55c69996c7f22051f6c3b94d/dash/manifest.mpd"); err != nil {
		t.Error(err)
	}

	mpdProcessor := NewMpdProcessor()
	mpdProcessor.Process(mpd)

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

	if len(mpdProcessor.ManifestInfo.PeriodInfos[0].StreamSetInfos[0].StreamInfos[0].SegmentIndex.References) != 424 {
		t.Errorf("expecting 424 references, got %d", len(mpdProcessor.ManifestInfo.PeriodInfos[0].StreamSetInfos[0].StreamInfos[0].SegmentIndex.References))
	}

	if len(mpdProcessor.ManifestInfo.PeriodInfos[0].StreamSetInfos[1].StreamInfos[0].SegmentIndex.References) != 424 {
		t.Errorf("expecting 424 references, got %d", len(mpdProcessor.ManifestInfo.PeriodInfos[0].StreamSetInfos[1].StreamInfos[0].SegmentIndex.References))
	}

	if mpdProcessor.ManifestInfo.PeriodInfos[0].StreamSetInfos[0].StreamInfos[0].SegmentIndex.References[423].Url != "http://sdk.streamrail.com/pepsi/cdn/0.0.1/3a5dd80efc3a867e55c69996c7f22051f6c3b94d/dash/470k/audio/und/seg-424.m4f" {
		t.Errorf("expecting last reference to point to %s, got:%s", "http://sdk.streamrail.com/pepsi/cdn/0.0.1/3a5dd80efc3a867e55c69996c7f22051f6c3b94d/dash/470k/audio/und/seg-424.m4f", mpdProcessor.ManifestInfo.PeriodInfos[0].StreamSetInfos[0].StreamInfos[0].SegmentIndex.References[423].Url)
	}

	if mpdProcessor.ManifestInfo.PeriodInfos[0].StreamSetInfos[1].StreamInfos[0].SegmentIndex.References[423].Url != "http://sdk.streamrail.com/pepsi/cdn/0.0.1/3a5dd80efc3a867e55c69996c7f22051f6c3b94d/dash/470k/video/1/seg-424.m4f" {
		t.Errorf("expecting last reference to point to %s, got:%s", "http://sdk.streamrail.com/pepsi/cdn/0.0.1/3a5dd80efc3a867e55c69996c7f22051f6c3b94d/dash/470k/video/1/seg-424.m4f", mpdProcessor.ManifestInfo.PeriodInfos[0].StreamSetInfos[1].StreamInfos[0].SegmentIndex.References[423].Url)
	}
}

func MPDProcessingExample2(t *testing.T) {
	var mpd *Mpd
	var err error

	if mpd, err = ParseMpd("http://sdk.streamrail.com/pepsi/cdn/0.0.1/925e302c164efcbe473977cff27771a3e1184902/dash/manifest.mpd"); err != nil {
		t.Error(err)
	}

	mpdProcessor := NewMpdProcessor()
	mpdProcessor.Process(mpd)

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

	if mpdProcessor.ManifestInfo.PeriodInfos[0].StreamSetInfos[0].StreamInfos[0].SegmentIndex.References[136].Url != "http://sdk.streamrail.com/pepsi/cdn/0.0.1/925e302c164efcbe473977cff27771a3e1184902/dash/700k/audio/und/seg-137.m4f" {
		t.Errorf("expecting last reference to point to %s, got:%s", "http://sdk.streamrail.com/pepsi/cdn/0.0.1/925e302c164efcbe473977cff27771a3e1184902/dash/700k/audio/und/seg-137.m4f", mpdProcessor.ManifestInfo.PeriodInfos[0].StreamSetInfos[0].StreamInfos[0].SegmentIndex.References[136].Url)
	}

	if mpdProcessor.ManifestInfo.PeriodInfos[0].StreamSetInfos[1].StreamInfos[0].SegmentIndex.References[136].Url != "http://sdk.streamrail.com/pepsi/cdn/0.0.1/925e302c164efcbe473977cff27771a3e1184902/dash/700k/video/1/seg-137.m4f" {
		t.Errorf("expecting last reference to point to %s, got:%s", "http://sdk.streamrail.com/pepsi/cdn/0.0.1/925e302c164efcbe473977cff27771a3e1184902/dash/700k/video/1/seg-137.m4f", mpdProcessor.ManifestInfo.PeriodInfos[0].StreamSetInfos[1].StreamInfos[0].SegmentIndex.References[136].Url)
	}

	if mpdProcessor.ManifestInfo.PeriodInfos[0].StreamSetInfos[1].StreamInfos[1].SegmentIndex.References[136].Url != "http://sdk.streamrail.com/pepsi/cdn/0.0.1/925e302c164efcbe473977cff27771a3e1184902/dash/1200k/video/1/seg-137.m4f" {
		t.Errorf("expecting last reference to point to %s, got:%s", "http://sdk.streamrail.com/pepsi/cdn/0.0.1/925e302c164efcbe473977cff27771a3e1184902/dash/1200k/video/1/seg-137.m4f", mpdProcessor.ManifestInfo.PeriodInfos[0].StreamSetInfos[1].StreamInfos[0].SegmentIndex.References[136].Url)
	}

	if mpdProcessor.ManifestInfo.PeriodInfos[0].StreamSetInfos[1].StreamInfos[2].SegmentIndex.References[136].Url != "http://sdk.streamrail.com/pepsi/cdn/0.0.1/925e302c164efcbe473977cff27771a3e1184902/dash/1531k/video/1/seg-137.m4f" {
		t.Errorf("expecting last reference to point to %s, got:%s", "http://sdk.streamrail.com/pepsi/cdn/0.0.1/925e302c164efcbe473977cff27771a3e1184902/dash/1531k/video/1/seg-137.m4f", mpdProcessor.ManifestInfo.PeriodInfos[0].StreamSetInfos[1].StreamInfos[0].SegmentIndex.References[136].Url)
	}
}
