package metric

import (
	"fmt"
	"regexp"
	"time"
)

var dayNames = [10]string{
	"Primday", "Duoday", "Triday", "Quadday", "Quintday",
	"Hexday", "Septday", "Octday", "Novday", "Decday",
}

var monthNames = [12]string{
	"Unil", "Duil", "Tril", "Quadril", "Quintil", "Sextil",
	"Septil", "Octil", "Novil", "Decil", "Undecil", "Duodecil",
}

var seasonNames = [4]string{"Rising", "Flourishing", "Gathering", "Stillness"}

var turningDayNames = [3]string{"Vigil", "Balance", "Dawn"}
var yuleDayNames = [3]string{"Yule Eve", "Midwinter", "Kindling"}

var formatRe = regexp.MustCompile(`MMM|MM|M|DD|D|WW|W|Y|S`)

// Date represents a date in the Metric Calendar.
type Date struct {
	Year         int
	Month        int    // 1-12; 0 for Turning/Yule days
	MonthName    string
	Day          int    // 1-30; 0 for Turning/Yule days
	WeekDay      int    // 1-10; 0 for Turning/Yule days
	DayName      string
	Week         int    // 1-36; 0 for Turning/Yule days
	SeasonIndex  int    // 0-3; -1 for Turning/Yule days
	IsLeapYear   bool
	IsTurning    bool
	IsYule       bool
	IsMidsummer  bool   // true when Month==4 && Day==1
	IsSextant    bool   // true when Month==2 && Day==30
	IsTrine      bool   // true when Month==4 && Day==30
	IsSpiral     bool   // true when Month==5 && Day==18
	IsConvergence bool  // true when Month==5 && Day==24
	IsMeridian   bool   // true when Month==6 && Day==30
	IsMask       bool   // true when Month==8 && Day==13
	IsHarmony    bool   // true when Month==8 && Day==30
	IsRest       bool   // true when WeekDay >= 8
	SpecialDay   string // name set for Turning and Yule days
	Observance   string // name of observance, or "" if none
}

// FromGregorian converts a Gregorian time.Time to a Metric Calendar Date.
// The time's clock component is ignored; only the date is used.
func FromGregorian(t time.Time) Date {
	t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
	year := t.Year()
	equinox := time.Date(year, 3, 20, 0, 0, 0, 0, time.UTC)
	daysFromEquinox := int(t.Sub(equinox).Hours() / 24)

	var metricYear, dayOfYear int
	if daysFromEquinox >= 0 {
		metricYear = year - 1970
		dayOfYear = daysFromEquinox + 1
	} else {
		metricYear = year - 1 - 1970
		prevEquinox := time.Date(year-1, 3, 20, 0, 0, 0, 0, time.UTC)
		dayOfYear = int(t.Sub(prevEquinox).Hours()/24) + 1
	}

	leap := isLeapYear(metricYear + 1971)
	yuleDayCount := 2
	if leap {
		yuleDayCount = 3
	}

	result := Date{Year: metricYear, IsLeapYear: leap, SeasonIndex: -1}

	if dayOfYear <= 3 {
		result.IsTurning = true
		result.SpecialDay = turningDayNames[dayOfYear-1]
		result.Observance = result.SpecialDay
		return result
	}

	adjusted := dayOfYear - 3
	if adjusted <= 270 {
		result.Month = (adjusted-1)/30 + 1
		result.Day = (adjusted-1)%30 + 1
	} else if adjusted <= 270+yuleDayCount {
		result.IsYule = true
		result.SpecialDay = yuleDayNames[adjusted-271]
		result.Observance = result.SpecialDay
		return result
	} else {
		postYule := adjusted - 270 - yuleDayCount
		result.Month = 9 + (postYule-1)/30 + 1
		result.Day = (postYule-1)%30 + 1
	}

	result.WeekDay = (result.Day-1)%10 + 1
	result.IsRest = result.WeekDay >= 8
	result.IsMidsummer = result.Month == 4 && result.Day == 1
	result.IsSextant = result.Month == 2 && result.Day == 30
	result.IsTrine = result.Month == 4 && result.Day == 30
	result.IsSpiral = result.Month == 5 && result.Day == 18
	result.IsConvergence = result.Month == 5 && result.Day == 24
	result.IsMeridian = result.Month == 6 && result.Day == 30
	result.IsMask = result.Month == 8 && result.Day == 13
	result.IsHarmony = result.Month == 8 && result.Day == 30
	result.DayName = dayNames[result.WeekDay-1]
	result.MonthName = monthNames[result.Month-1]
	result.SeasonIndex = (result.Month - 1) / 3
	result.Week = (result.Month-1)*3 + (result.Day-1)/10 + 1

	switch {
	case result.IsMidsummer:
		result.Observance = "Midsummer"
	case result.IsSextant:
		result.Observance = "The Sextant"
	case result.IsTrine:
		result.Observance = "The Trine"
	case result.IsSpiral:
		result.Observance = "The Spiral"
	case result.IsConvergence:
		result.Observance = "Convergence"
	case result.IsMeridian:
		result.Observance = "The Meridian"
	case result.IsMask:
		result.Observance = "The Mask"
	case result.IsHarmony:
		result.Observance = "Harmony"
	}

	return result
}

// Format formats a Date using a pattern string.
//
// Tokens:
//
//	MMM  month name (e.g. "Unil")
//	MM   month zero-padded (e.g. "01")
//	M    month number (e.g. "1")
//	DD   day zero-padded (e.g. "04")
//	D    day number (e.g. "4")
//	WW   weekday name (e.g. "Quintday")
//	W    weekday number (e.g. "5")
//	Y    year number (e.g. "56")
//	S    season name (e.g. "Rising")
//
// Example: Format(d, "WW, MMM D, Year Y") → "Quintday, Unil 4, Year 56"
func Format(d Date, pattern string) string {
	seasonName := ""
	if d.SeasonIndex >= 0 && d.SeasonIndex <= 3 {
		seasonName = seasonNames[d.SeasonIndex]
	}
	return formatRe.ReplaceAllStringFunc(pattern, func(tok string) string {
		switch tok {
		case "MMM":
			return d.MonthName
		case "MM":
			return fmt.Sprintf("%02d", d.Month)
		case "M":
			return fmt.Sprintf("%d", d.Month)
		case "DD":
			return fmt.Sprintf("%02d", d.Day)
		case "D":
			return fmt.Sprintf("%d", d.Day)
		case "WW":
			return d.DayName
		case "W":
			return fmt.Sprintf("%d", d.WeekDay)
		case "Y":
			return fmt.Sprintf("%d", d.Year)
		case "S":
			return seasonName
		}
		return tok
	})
}

// ToGregorian converts a Metric Calendar date back to a Gregorian time.Time.
//
// year is the Metric year (0 = spring equinox 1970).
// periodType is one of "turning", "month", or "yule".
// periodValue meaning depends on periodType:
//   - "turning": 0=Vigil, 1=Balance, 2=Dawn
//   - "month":   1-12 (the month number); dayOfMonth must be 1-30
//   - "yule":    0=Yule Eve, 1=Midwinter, 2=Kindling (leap years only)
//
// dayOfMonth is only used when periodType is "month".
func ToGregorian(year int, periodType string, periodValue int, dayOfMonth int) (time.Time, error) {
	equinoxYear := year + 1970
	leap := isLeapYear(year + 1971)
	yuleDayCount := 2
	if leap {
		yuleDayCount = 3
	}

	var offset int
	switch periodType {
	case "turning":
		if periodValue < 0 || periodValue > 2 {
			return time.Time{}, fmt.Errorf("turning period value must be 0-2")
		}
		offset = periodValue
	case "month":
		m := periodValue
		d := dayOfMonth
		if m < 1 || m > 12 {
			return time.Time{}, fmt.Errorf("month must be 1-12")
		}
		if d < 1 || d > 30 {
			return time.Time{}, fmt.Errorf("day must be 1-30")
		}
		if m <= 9 {
			offset = 3 + (m-1)*30 + (d - 1)
		} else {
			offset = 3 + 270 + yuleDayCount + (m-10)*30 + (d - 1)
		}
	case "yule":
		if periodValue == 2 && !leap {
			return time.Time{}, fmt.Errorf("Kindling only occurs in leap years")
		}
		if periodValue < 0 || periodValue > 2 {
			return time.Time{}, fmt.Errorf("yule period value must be 0-2")
		}
		offset = 3 + 270 + periodValue
	default:
		return time.Time{}, fmt.Errorf("unknown period type: %s", periodType)
	}

	equinox := time.Date(equinoxYear, 3, 20, 0, 0, 0, 0, time.UTC)
	return equinox.AddDate(0, 0, offset), nil
}

// IsRestDay reports whether the given Gregorian date falls on a Metric rest day
// (WeekDay 8, 9, or 10 — Octday, Novday, or Decday).
func IsRestDay(t time.Time) bool {
	return FromGregorian(t).IsRest
}

// Today returns the Metric Calendar Date for the current local date.
func Today() Date {
	return FromGregorian(time.Now())
}

func isLeapYear(y int) bool {
	return y%4 == 0 && (y%100 != 0 || y%400 == 0)
}
