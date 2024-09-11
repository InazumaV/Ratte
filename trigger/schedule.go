package trigger

import "time"

type Schedule struct {
	interval int
}

func newSchedule(interval int) *Schedule {
	return &Schedule{
		interval: interval,
	}
}

func (s *Schedule) Next(t time.Time) time.Time {
	return t.Add(time.Duration(s.interval) * time.Second)
}
