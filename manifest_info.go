package main

type ManifestInfo struct {
	Live bool

	MinBufferTime int

	PeriodInfos []PeriodInfo
}

func NewManifestInfo() ManifestInfo {
	return ManifestInfo{
		Live:          false,
		MinBufferTime: 0,
		PeriodInfos:   make([]PeriodInfo, 0),
	}
}
