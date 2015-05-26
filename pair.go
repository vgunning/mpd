package mpd

type Pair struct {
	Earliest int
	Current  int
}

func NewPair(earliestSegmentNumber, currentSegmentNumber int) Pair {
	return Pair{
		Earliest: earliestSegmentNumber,
		Current:  currentSegmentNumber,
	}
}
