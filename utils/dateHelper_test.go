package utils

import (
	"testing"
	"time"
)

func TestRangeDate(t *testing.T) {
	start := time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2022, time.January, 5, 0, 0, 0, 0, time.UTC)

	rangeDate := RangeDate(start, end)

	expectedDates := []time.Time{
		time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2022, time.January, 2, 0, 0, 0, 0, time.UTC),
		time.Date(2022, time.January, 3, 0, 0, 0, 0, time.UTC),
		time.Date(2022, time.January, 4, 0, 0, 0, 0, time.UTC),
		time.Date(2022, time.January, 5, 0, 0, 0, 0, time.UTC),
	}

	for _, expected := range expectedDates {
		date := rangeDate()
		if !date.Equal(expected) {
			t.Errorf("RangeDate() = %s; want %s", date.Format(LayoutPostgres), expected.Format(LayoutPostgres))
		}
	}
}

func TestRangeDateMonthly(t *testing.T) {
	start := time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2022, time.March, 1, 0, 0, 0, 0, time.UTC)

	rangeDateMonthly := RangeDateMonthly(start, end)

	expectedDates := []time.Time{
		time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2022, time.February, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2022, time.March, 1, 0, 0, 0, 0, time.UTC),
	}

	for _, expected := range expectedDates {
		date := rangeDateMonthly()
		if !date.Equal(expected) {
			t.Errorf("RangeDateMonthly() = %s; want %s", date.Format(LayoutPostgres), expected.Format(LayoutPostgres))
		}
	}
}

func TestRangeDateOffsetByDate(t *testing.T) {
	start := time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2022, time.January, 5, 0, 0, 0, 0, time.UTC)
	offsetDays := 2

	rangeDateOffsetByDate := RangeDateOffsetByDate(start, end, &offsetDays)

	expectedDates := []time.Time{
		time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2022, time.January, 3, 0, 0, 0, 0, time.UTC),
		time.Date(2022, time.January, 5, 0, 0, 0, 0, time.UTC),
	}

	for _, expected := range expectedDates {
		date := rangeDateOffsetByDate()
		if !date.Equal(expected) {
			t.Errorf("RangeDateOffsetByDate() = %s; want %s", date.Format(LayoutPostgres), expected.Format(LayoutPostgres))
		}
	}
}

func TestConvertDateToUTCZero(t *testing.T) {
	t1 := time.Date(2022, time.January, 1, 10, 30, 0, 0, time.UTC)
	expected := time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC)

	result := ConvertDateToUTCZero(t1)

	if !result.Equal(expected) {
		t.Errorf("ConvertDateToUTCZero() = %s; want %s", result.Format(LayoutPostgres), expected.Format(LayoutPostgres))
	}
}

func TestEndOfMonthDate(t *testing.T) {
	t1 := time.Date(2022, time.January, 15, 0, 0, 0, 0, time.UTC)
	expected := time.Date(2022, time.January, 31, 0, 0, 0, 0, time.UTC)

	result := EndOfMonthDate(t1)

	if !result.Equal(expected) {
		t.Errorf("EndOfMonthDate() = %s; want %s", result.Format(LayoutPostgres), expected.Format(LayoutPostgres))
	}
}

func TestOffsetDateByDay(t *testing.T) {
	t1 := time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC)
	offsetDay := 5
	expected := time.Date(2022, time.January, 6, 23, 59, 59, int(time.Second-time.Nanosecond), time.UTC)

	result := OffsetDateByDay(t1, &offsetDay)

	if !result.Equal(expected) {
		t.Errorf("OffsetDateByDay() = %s; want %s", result.Format(LayoutPostgres), expected.Format(LayoutPostgres))
	}
}

func TestEndDate(t *testing.T) {
	t1 := time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC)
	expected := time.Date(2022, time.January, 2, 0, 0, 0, -1, time.UTC)

	result := EndDate(t1)

	if !result.Equal(expected) {
		t.Errorf("EndDate() = %s; want %s", result.Format(LayoutPostgres), expected.Format(LayoutPostgres))
	}
}

func TestParseStringToDate(t *testing.T) {
	s := "2022-01-01"
	layout := LayoutISO
	expected := time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC)

	result, err := ParseStringToDate(s, layout)

	if err != nil {
		t.Errorf("ParseStringToDate() returned an error: %v", err)
	}

	if !result.Equal(expected) {
		t.Errorf("ParseStringToDate() = %s; want %s", result.Format(LayoutPostgres), expected.Format(LayoutPostgres))
	}
}

func TestTimeDiff(t *testing.T) {
	t1 := time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC)
	t2 := time.Date(2022, time.January, 2, 0, 0, 0, 0, time.UTC)
	expected := time.Hour * 24

	result := TimeDiff(t1, t2)

	if !result.Equal(time.Time{}.Add(expected)) {
		t.Errorf("TimeDiff() = %s; want %s", result.String(), expected.String())
	}
}
