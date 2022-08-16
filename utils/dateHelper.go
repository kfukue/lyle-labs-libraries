package utils

import (
	"fmt"
	"log"
	"time"
)

const (
	LayoutISO       = "2006-01-02"
	LayoutUS        = "January 2, 2006"
	LayoutPostgres  = "2006-01-02 15:04:05"
	LayoutCoingecko = "02-01-2006"
	LayoutRFC3339   = "2016-06-20T12:41:45.140Z"
)

func RangeDate(start, end time.Time) func() time.Time {
	y, m, d := start.Date()
	start = time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
	y, m, d = end.Date()
	end = time.Date(y, m, d, 0, 0, 0, 0, time.UTC)

	return func() time.Time {
		if start.After(end) {
			return time.Time{}
		}
		date := start
		start = start.AddDate(0, 0, 1)
		return date
	}
}

func RangeDateMonthly(start, end time.Time) func() time.Time {
	y, m, d := start.Date()
	firstRun := true
	start = time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
	initialStartDate := start
	y, m, d = end.Date()
	end = time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
	log.Println(fmt.Sprintf("RangeDateMonthly, start: %s, end : %s", start.Format(LayoutPostgres), end.Format(LayoutPostgres)))
	return func() time.Time {
		// log.Println(fmt.Sprintf("start initial: %s,  start: %s, end initial : %s,  end : %s", initialStartDate.Format(LayoutPostgres), start.Format(LayoutPostgres), end.Format(LayoutPostgres), initialEndDate.Format(LayoutPostgres)))
		if start.After(end) {
			return time.Time{}
		}
		if start.Equal(initialStartDate) && firstRun == true {
			firstRun = false
			return start
		}
		currentYear, currentMonth := start.Year(), start.Month()
		start = time.Date(currentYear, currentMonth+1, 1, 0, 0, 0, 0, time.UTC)
		return start
	}
}

func RangeDateOffsetByDate(start, end time.Time, offSetDays *int) func() time.Time {
	y, m, d := start.Date()
	firstRun := true
	start = time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
	initialStartDate := start
	y, m, d = end.Date()
	end = time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
	log.Println(fmt.Sprintf("RangeDateMonthly, start: %s, end : %s", start.Format(LayoutPostgres), end.Format(LayoutPostgres)))
	return func() time.Time {
		// log.Println(fmt.Sprintf("start initial: %s,  start: %s, end initial : %s,  end : %s", initialStartDate.Format(LayoutPostgres), start.Format(LayoutPostgres), end.Format(LayoutPostgres), initialEndDate.Format(LayoutPostgres)))
		if start.After(end) {
			return time.Time{}
		}
		if start.Equal(initialStartDate) && firstRun == true {
			firstRun = false
			return start
		}
		currentYear, currentMonth, currentDay := start.Year(), start.Month(), start.Day()
		start = time.Date(currentYear, currentMonth, currentDay+*offSetDays, 0, 0, 0, 0, time.UTC)
		return start
	}
}

func ConvertDateToUTCZero(t time.Time) time.Time {
	currentYear, currentMonth, currentDay := t.Year(), t.Month(), t.Day()
	firstDay := time.Date(currentYear, currentMonth, currentDay, 0, 0, 0, 0, time.UTC)
	return firstDay
}

func EndOfMonthDate(t time.Time) time.Time {
	currentYear, currentMonth := t.Year(), t.Month()
	firstDay := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, time.UTC)
	lastDay := firstDay.AddDate(0, 1, 0).Add(time.Nanosecond * -1)
	return lastDay
}

func OffsetDateByDay(t time.Time, offsetDay *int) time.Time {
	currentYear, currentMonth, currentDay := t.Year(), t.Month(), t.Day()
	firstDay := time.Date(currentYear, currentMonth, currentDay+*offsetDay, 0, 0, 0, 0, time.UTC)
	lastDay := firstDay.AddDate(0, 0, 1).Add(time.Nanosecond * -1)
	return lastDay
}

func EndDate(t time.Time) time.Time {
	wd := t.Weekday()
	next := int((wd - t.Weekday() + 7) % 7)
	y, m, d := t.Date()
	return time.Date(y, m, d+next+1, 0, 0, 0, -1, t.Location())
}

func ParseStringToDate(s string, layout string) (time.Time, error) {
	result, err := time.Parse(layout, s)
	if err != nil {
		log.Println(err.Error())
		return time.Time{}, err
	}
	return result, err
}

func TimeDiff(TimeA time.Time, TimeB time.Time) time.Time {
	timeDiff := TimeB.Sub(TimeA)
	out := time.Time{}.Add(timeDiff)
	return out
}
