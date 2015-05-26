package mpd

import (
	"fmt"
)

type SegmentIndex struct {
	References          []*SegmentReference
	TimestampCorrection int
}

/**
 * Creates a SegmentIndex.
 *
 * @param {!Array.<!SegmentReference>} references Sorted by time in
 *     ascending order with no gaps.
 */
func NewSegmentIndex(references []*SegmentReference) SegmentIndex {
	return SegmentIndex{
		References:          references,
		TimestampCorrection: 0,
	}

	// assertCorrectReferences()
}

/**
 * Gets the number of SegmentReferences.
 *
 * @return {number}
 */
func (segmentIndex SegmentIndex) Length() int {
	return len(segmentIndex.References)
}

/**
 * Gets the last SegmentReference.
 *
 * @return {!SegmentReference} The last SegmentReference.
 * @throws {RangeError} when there are no SegmentReferences.
 */
func (segmentIndex SegmentIndex) Last() *SegmentReference {
	if len(segmentIndex.References) == 0 {
		fmt.Println("SegmentIndex: There is no last SegmentReference.")
		// throw new RangeError('SegmentIndex: There is no last SegmentReference.')
		return nil
	}

	return segmentIndex.References[len(segmentIndex.References)-1]
}
