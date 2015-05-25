package main

/**
 * The next unique ID to assign to a StreamSetInfo.
 */
var nextUniqueId int = 0

type StreamInfo struct {
	UniqueId int

	Id string

	/**
	 * An offset, in seconds, to apply to each timestamp within each media
	 * segment that's put in buffer.
	 */
	TimestampOffset int

	/**
	 * Indicates the stream's current segment's start time, i.e., its live-edge.
	 * This value is non-null if the stream is both live and available;
	 * otherwise, this value is null.
	 */
	CurrentSegmentStartTime int

	/**
	 * Bandwidth required, in bits per second, to assure uninterrupted playback,
	 * assuming that |minBufferTime| seconds of video are in buffer before
	 * playback begins.
	 */
	Bandwidth int

	Width int

	Height int

	MimeType string

	Codecs string

	MediaUrl string

	Enabled bool

	/**
	 * The stream's SegmentIndex metadata.
	 * @see {StreamInfo.isAvailable}
	 */
	SegmentIndexInfo *SegmentMetadataInfo

	/**
	 * The stream's segment initialization metadata.
	 * @type {SegmentMetadataInfo}
	 */
	SegmentInitializationInfo *SegmentMetadataInfo

	/**
	 * The stream's SegmentIndex.
	 * @see {StreamInfo.isAvailable}
	 * @type {SegmentIndex}
	 */
	SegmentIndex *SegmentIndex

	/** @type {ArrayBuffer} */
	SegmentInitializationData []byte

	/** @private {ArrayBuffer} */
	SegmentIndexData []byte
}

func NewStreamInfo() StreamInfo {
	nextUniqueId++

	return StreamInfo{
		UniqueId:                  nextUniqueId,
		Id:                        "",
		TimestampOffset:           0,
		CurrentSegmentStartTime:   -1,
		Bandwidth:                 -1,
		Width:                     -1,
		Height:                    -1,
		MimeType:                  "",
		Codecs:                    "",
		MediaUrl:                  "",
		Enabled:                   true,
		SegmentIndexInfo:          nil,
		SegmentInitializationInfo: nil,
		SegmentIndex:              nil,
		SegmentInitializationData: nil,
		SegmentIndexData:          nil,
	}
}
