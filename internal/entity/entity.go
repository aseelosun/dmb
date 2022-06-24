package entity

import "time"

type MySchedule struct {
	Id    int
	MType string
	MDate string
	MTime string
}

type MyProfile struct {
	Name     string
	Email    string
	PhoneNum string
}

type AllSchedule struct {
	MyProfile
	MySchedule
}

func (s *MySchedule) DayOfWeek() time.Weekday {
	date, _ := time.Parse("2006-01-02", s.MDate)

	return date.Weekday()
}
