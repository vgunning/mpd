package mpd

import (
	mapset "github.com/deckarep/golang-set"
)

/**
 * The next unique ID to assign to a StreamSetInfo.
 */
var streamSetInfoNextUniqueId int = 0

type StreamSetInfo struct {
	/** @type {number} */
	UniqueId int

	/** @type {?string} */
	Id string

	/** @type {string} */
	ContentType mapset.Set

	/** @type {!Array.<!StreamInfo>} */
	StreamInfos []*StreamInfo

	/** @type {!Array.<!DrmSchemeInfo>} */
	// DrmSchemes = DrmSchemeInfo[]

	/** @type {string} */
	Lang string

	/** @type {boolean} */
	Main bool
}

func NewStreamSetInfo() StreamSetInfo {
	streamSetInfoNextUniqueId++

	return StreamSetInfo{
		UniqueId:    streamSetInfoNextUniqueId,
		Id:          "",
		ContentType: mapset.NewSet(),
		StreamInfos: make([]*StreamInfo, 0),

		/** @type {!Array.<!DrmSchemeInfo>} */
		// DrmSchemes: make(DrmSchemeInfo[], 0),

		Lang: "",
		Main: false,
	}
}
