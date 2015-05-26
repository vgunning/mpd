package mpd

type SegmentMetadataInfo struct {
	Url string

	StartByte int

	EndByte int
}

func NewSegmentMetadataInfo() SegmentMetadataInfo {
	return SegmentMetadataInfo{
		Url:       "",
		StartByte: 0,
		EndByte:   -1,
	}
}
