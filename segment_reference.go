package mpd

type SegmentReference struct {

	/**
	 * The segment's ID.
	 * @const {number}
	 */
	Id int

	/**
	 * The time, in seconds, that the segment begins.
	 * @const {number}
	 */
	StartTime int

	/**
	 * The time, in seconds, that the segment ends. The segment ends immediately
	 * before this time. A null value indicates that the segment continues to the
	 * end of the stream.
	 * @const {?number}
	 */
	EndTime int

	/**
	 * The position of the segment's first byte.
	 * @const {number}
	 */
	StartByte int

	/**
	 * The position of the segment's last byte, inclusive. A null value indicates
	 * that the segment continues to the end of the file located at |url|.
	 * @const {?number}
	 */
	EndByte int

	/**
	 * The segment's location.
	 * @const {!string}
	 */
	Url string
}

func NewSegmentReference(id, startTime, endTime, startByte, endByte int, url string) SegmentReference {
	// assert((endTime == -1) || (startTime <= endTime), "startTime should be <= endTime")
	assert((endTime == -1) || (startTime <= endTime))

	return SegmentReference{
		Id: id,

		StartTime: startTime,

		EndTime: endTime,

		StartByte: startByte,

		EndByte: endByte,

		Url: url,
	}
}
