package metric_test

import (
	"testing"
	"time"

	metric "github.com/alexrabarts/metric-calendar/libs/go"
)

func TestFromGregorian(t *testing.T) {
	tests := []struct {
		date        time.Time
		wantYear    int
		wantMonth   int
		wantDay     int
		wantWeekDay int
		wantDayName string
		wantMonthName string
		wantWeek    int
		wantIsRest  bool
		wantIsTurning bool
		wantIsYule  bool
		wantSpecialDay string
	}{
		{
			date:           time.Date(2026, 3, 20, 0, 0, 0, 0, time.UTC),
			wantYear:       56,
			wantIsTurning:  true,
			wantSpecialDay: "Vigil",
		},
		{
			date:           time.Date(2026, 3, 21, 0, 0, 0, 0, time.UTC),
			wantYear:       56,
			wantIsTurning:  true,
			wantSpecialDay: "Balance",
		},
		{
			date:           time.Date(2026, 3, 22, 0, 0, 0, 0, time.UTC),
			wantYear:       56,
			wantIsTurning:  true,
			wantSpecialDay: "Dawn",
		},
		{
			date:          time.Date(2026, 3, 23, 0, 0, 0, 0, time.UTC),
			wantYear:      56,
			wantMonth:     1,
			wantDay:       1,
			wantWeekDay:   1,
			wantDayName:   "Primday",
			wantMonthName: "Unil",
			wantWeek:      1,
			wantIsRest:    false,
		},
		{
			date:          time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC),
			wantYear:      56,
			wantMonth:     1,
			wantDay:       10,
			wantWeekDay:   10,
			wantDayName:   "Decday",
			wantMonthName: "Unil",
			wantWeek:      1,
			wantIsRest:    true,
		},
		{
			date:           time.Date(2026, 12, 18, 0, 0, 0, 0, time.UTC),
			wantYear:       56,
			wantIsYule:     true,
			wantSpecialDay: "Yule Eve",
		},
		{
			// 2025-01-01 is before the spring equinox, so it belongs to metric year 54.
			// Metric year 55 begins on 2025-03-20.
			date:          time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			wantYear:      54,
			wantMonth:     10,
			wantDay:       13,
			wantWeekDay:   3,
			wantDayName:   "Triday",
			wantMonthName: "Decil",
			wantWeek:      29,
			wantIsRest:    false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.date.Format("2006-01-02"), func(t *testing.T) {
			got := metric.FromGregorian(tc.date)

			if got.Year != tc.wantYear {
				t.Errorf("Year: got %d, want %d", got.Year, tc.wantYear)
			}
			if got.IsTurning != tc.wantIsTurning {
				t.Errorf("IsTurning: got %v, want %v", got.IsTurning, tc.wantIsTurning)
			}
			if got.IsYule != tc.wantIsYule {
				t.Errorf("IsYule: got %v, want %v", got.IsYule, tc.wantIsYule)
			}
			if tc.wantSpecialDay != "" && got.SpecialDay != tc.wantSpecialDay {
				t.Errorf("SpecialDay: got %q, want %q", got.SpecialDay, tc.wantSpecialDay)
			}
			if !tc.wantIsTurning && !tc.wantIsYule {
				if got.Month != tc.wantMonth {
					t.Errorf("Month: got %d, want %d", got.Month, tc.wantMonth)
				}
				if got.Day != tc.wantDay {
					t.Errorf("Day: got %d, want %d", got.Day, tc.wantDay)
				}
				if got.WeekDay != tc.wantWeekDay {
					t.Errorf("WeekDay: got %d, want %d", got.WeekDay, tc.wantWeekDay)
				}
				if tc.wantDayName != "" && got.DayName != tc.wantDayName {
					t.Errorf("DayName: got %q, want %q", got.DayName, tc.wantDayName)
				}
				if tc.wantMonthName != "" && got.MonthName != tc.wantMonthName {
					t.Errorf("MonthName: got %q, want %q", got.MonthName, tc.wantMonthName)
				}
				if got.Week != tc.wantWeek {
					t.Errorf("Week: got %d, want %d", got.Week, tc.wantWeek)
				}
				if got.IsRest != tc.wantIsRest {
					t.Errorf("IsRest: got %v, want %v", got.IsRest, tc.wantIsRest)
				}
			}
		})
	}
}

func TestMidsummer(t *testing.T) {
	// Month 4, Day 1 — the summer solstice anchor
	d := metric.FromGregorian(time.Date(2026, 6, 21, 0, 0, 0, 0, time.UTC))
	if !d.IsMidsummer {
		t.Errorf("expected IsMidsummer=true for 2026-06-21, got Month=%d Day=%d", d.Month, d.Day)
	}
	if d.Month != 4 || d.Day != 1 {
		t.Errorf("expected Month=4 Day=1, got Month=%d Day=%d", d.Month, d.Day)
	}
}

func TestSpiral(t *testing.T) {
	// Month 5, Day 18 — the golden-ratio spiral day.
	// 3 turning days + (5-1)*30 + 17 days offset from equinox = 140 days after 2026-03-20 = 2026-08-07.
	d := metric.FromGregorian(time.Date(2026, 8, 7, 0, 0, 0, 0, time.UTC))
	if !d.IsSpiral {
		t.Errorf("expected IsSpiral=true for 2026-08-07, got Month=%d Day=%d", d.Month, d.Day)
	}
	if d.Month != 5 || d.Day != 18 {
		t.Errorf("expected Month=5 Day=18, got Month=%d Day=%d", d.Month, d.Day)
	}
}

func TestIsRestDay(t *testing.T) {
	tests := []struct {
		date     time.Time
		wantRest bool
	}{
		{time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC), true},  // Decday (WeekDay 10)
		{time.Date(2026, 3, 23, 0, 0, 0, 0, time.UTC), false}, // Primday (WeekDay 1)
	}
	for _, tc := range tests {
		t.Run(tc.date.Format("2006-01-02"), func(t *testing.T) {
			got := metric.IsRestDay(tc.date)
			if got != tc.wantRest {
				t.Errorf("IsRestDay(%s): got %v, want %v", tc.date.Format("2006-01-02"), got, tc.wantRest)
			}
		})
	}
}

func TestToGregorian(t *testing.T) {
	tests := []struct {
		name        string
		year        int
		periodType  string
		periodValue int
		dayOfMonth  int
		wantDate    time.Time
		wantErr     bool
	}{
		{
			name:        "month 1 day 1",
			year:        56,
			periodType:  "month",
			periodValue: 1,
			dayOfMonth:  1,
			wantDate:    time.Date(2026, 3, 23, 0, 0, 0, 0, time.UTC),
		},
		{
			name:        "turning Vigil",
			year:        56,
			periodType:  "turning",
			periodValue: 0,
			wantDate:    time.Date(2026, 3, 20, 0, 0, 0, 0, time.UTC),
		},
		{
			name:        "yule Yule Eve",
			year:        56,
			periodType:  "yule",
			periodValue: 0,
			wantDate:    time.Date(2026, 12, 18, 0, 0, 0, 0, time.UTC),
		},
		{
			name:        "invalid month",
			year:        56,
			periodType:  "month",
			periodValue: 13,
			dayOfMonth:  1,
			wantErr:     true,
		},
		{
			name:        "invalid day",
			year:        56,
			periodType:  "month",
			periodValue: 1,
			dayOfMonth:  31,
			wantErr:     true,
		},
		{
			name:        "Kindling in non-leap year",
			year:        56,
			periodType:  "yule",
			periodValue: 2,
			wantErr:     true,
		},
		{
			name:        "invalid period type",
			year:        56,
			periodType:  "quarter",
			periodValue: 1,
			dayOfMonth:  1,
			wantErr:     true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := metric.ToGregorian(tc.year, tc.periodType, tc.periodValue, tc.dayOfMonth)
			if tc.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil (result: %v)", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !got.Equal(tc.wantDate) {
				t.Errorf("got %s, want %s", got.Format("2006-01-02"), tc.wantDate.Format("2006-01-02"))
			}
		})
	}
}

func TestRoundTrip(t *testing.T) {
	// Convert a Gregorian date to Metric then back, and verify we get the same date.
	dates := []time.Time{
		time.Date(2026, 3, 23, 0, 0, 0, 0, time.UTC),
		time.Date(2026, 6, 21, 0, 0, 0, 0, time.UTC),
		time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2026, 12, 18, 0, 0, 0, 0, time.UTC),
	}
	for _, original := range dates {
		original := original
		t.Run(original.Format("2006-01-02"), func(t *testing.T) {
			d := metric.FromGregorian(original)
			var periodType string
			var periodValue, dayOfMonth int
			switch {
			case d.IsTurning:
				periodType = "turning"
				for i, name := range []string{"Vigil", "Balance", "Dawn"} {
					if d.SpecialDay == name {
						periodValue = i
					}
				}
			case d.IsYule:
				periodType = "yule"
				for i, name := range []string{"Yule Eve", "Midwinter", "Kindling"} {
					if d.SpecialDay == name {
						periodValue = i
					}
				}
			default:
				periodType = "month"
				periodValue = d.Month
				dayOfMonth = d.Day
			}
			got, err := metric.ToGregorian(d.Year, periodType, periodValue, dayOfMonth)
			if err != nil {
				t.Fatalf("ToGregorian error: %v", err)
			}
			if !got.Equal(original) {
				t.Errorf("round-trip: got %s, want %s", got.Format("2006-01-02"), original.Format("2006-01-02"))
			}
		})
	}
}

func TestLeapYear(t *testing.T) {
	// Metric year 57 corresponds to Gregorian year 2028, which is a leap year.
	d := metric.FromGregorian(time.Date(2027, 3, 23, 0, 0, 0, 0, time.UTC)) // first day of metric year 57
	if !d.IsLeapYear {
		t.Errorf("metric year 57 (Gregorian 2028) should be a leap year")
	}

	// Metric year 56 (Gregorian 2027) is not a leap year.
	d2 := metric.FromGregorian(time.Date(2026, 3, 23, 0, 0, 0, 0, time.UTC))
	if d2.IsLeapYear {
		t.Errorf("metric year 56 (Gregorian 2027) should not be a leap year")
	}
}

func TestSeasonIndex(t *testing.T) {
	tests := []struct {
		date        time.Time
		wantSeason  int
	}{
		{time.Date(2026, 3, 23, 0, 0, 0, 0, time.UTC), 0},  // Unil, Month 1 → Season 0
		{time.Date(2026, 6, 21, 0, 0, 0, 0, time.UTC), 1},  // Quadril, Month 4 → Season 1
		{time.Date(2026, 9, 23, 0, 0, 0, 0, time.UTC), 2},  // Month 7 → Season 2
		// 2026-12-10 is Month 9 (season 2); Decil (Month 10) starts 2026-12-20.
		{time.Date(2026, 12, 22, 0, 0, 0, 0, time.UTC), 3}, // Month 10 (Decil) → Season 3
	}
	for _, tc := range tests {
		t.Run(tc.date.Format("2006-01-02"), func(t *testing.T) {
			got := metric.FromGregorian(tc.date)
			if got.SeasonIndex != tc.wantSeason {
				t.Errorf("SeasonIndex: got %d, want %d (Month=%d)", got.SeasonIndex, tc.wantSeason, got.Month)
			}
		})
	}
}
