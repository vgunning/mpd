package mpd

type TimeLine struct {
	Start uint64
	End   uint64
}

func NewTimeLine(start, end uint64) TimeLine {
	return TimeLine{
		Start: start,
		End:   end,
	}
}
