package main

type PeriodInfo struct {
	Id string

	Start int

	/**
	 * The period's duration, in seconds.
	 */
	Duration int

	StreamSetInfos []StreamSetInfo
}

func NewPeriodInfo() PeriodInfo {
	return PeriodInfo{
		Id:             "",
		Start:          0,
		Duration:       -1,
		StreamSetInfos: make([]StreamSetInfo, 0),
	}
}
