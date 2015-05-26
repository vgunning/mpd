package mpd

type TimeLine struct {
	Start int
	End   int
}

func NewTimeLine(start, end int) TimeLine {
	return TimeLine{
		Start: start,
		End:   end,
	}
}
