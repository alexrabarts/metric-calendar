package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Calendar constants -- ported from www/index.html
var dayNames = [10]string{
	"Primday", "Duoday", "Triday", "Quadday", "Quintday",
	"Hexday", "Septday", "Octday", "Novday", "Decday",
}

var monthNames = [12]string{
	"Unil", "Duil", "Tril", "Quadril", "Quintil", "Sextil",
	"Septil", "Octil", "Novil", "Decil", "Undecil", "Duodecil",
}

var seasonNamesNorth = [4]string{"Rising", "Flourishing", "Gathering", "Stillness"}
var seasonNamesSouth = [4]string{"Gathering", "Stillness", "Rising", "Flourishing"}

var turningDayNames = [3]string{"Vigil", "Balance", "Dawn"}
var yuleDayNames = [3]string{"Yule Eve", "Solstice", "Kindling"}

// MetricDate holds the result of converting a Gregorian date.
type MetricDate struct {
	Year        int
	IsLeapYear  bool
	IsTurning   bool
	IsYule      bool
	IsMidsummer    bool
	IsSpiral       bool
	IsSextant      bool
	IsTrine        bool
	IsConvergence  bool
	IsMeridian     bool
	IsHarmony      bool
	IsMask         bool
	IsRest      bool
	Month       int // 1-12
	Day         int // 1-30
	WeekDay     int // 1-10
	DayName     string
	MonthName   string
	SeasonIndex int // 0-3
	SpecialDay  string
	Week        int // 1-36
}

func isGregorianLeapYear(y int) bool {
	return y%4 == 0 && (y%100 != 0 || y%400 == 0)
}

func gregorianToMetric(t time.Time) MetricDate {
	year := t.Year()
	month := int(t.Month())
	day := t.Day()

	equinox := time.Date(year, 3, 20, 0, 0, 0, 0, time.UTC)
	date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)

	daysFromEquinox := int(date.Sub(equinox).Hours() / 24)

	var metricYear int
	var dayOfYear int

	if daysFromEquinox >= 0 {
		metricYear = year - 1970
		dayOfYear = daysFromEquinox + 1
	} else {
		metricYear = year - 1 - 1970
		prevEquinox := time.Date(year-1, 3, 20, 0, 0, 0, 0, time.UTC)
		dayOfYear = int(date.Sub(prevEquinox).Hours()/24) + 1
	}

	// Metric year Y contains Feb 29 of Gregorian year (Y + 1971) if that year is a leap year
	leap := isGregorianLeapYear(metricYear + 1971)
	yuleDays := 2
	if leap {
		yuleDays = 3
	}

	result := MetricDate{
		Year:        metricYear,
		IsLeapYear:  leap,
		SeasonIndex: -1,
	}

	// The Turning (always 3 days)
	if dayOfYear <= 3 {
		result.IsTurning = true
		result.SpecialDay = turningDayNames[dayOfYear-1]
		return result
	}

	adjusted := dayOfYear - 3

	// Months 1-9: adjusted days 1-270
	if adjusted <= 270 {
		result.Month = (adjusted-1)/30 + 1
		result.Day = (adjusted-1)%30 + 1
	} else if adjusted <= 270+yuleDays {
		// Yule
		result.IsYule = true
		result.SpecialDay = yuleDayNames[adjusted-271]
		return result
	} else {
		// Months 10-12: after Yule
		postYule := adjusted - 270 - yuleDays
		result.Month = 9 + (postYule-1)/30 + 1
		result.Day = (postYule-1)%30 + 1
	}

	result.WeekDay = (result.Day-1)%10 + 1
	result.IsRest = result.WeekDay >= 8
	result.IsMidsummer = result.Month == 4 && result.Day == 1
	result.IsSpiral = result.Month == 5 && result.Day == 18
	result.IsSextant = result.Month == 2 && result.Day == 30
	result.IsTrine = result.Month == 4 && result.Day == 30
	result.IsConvergence = result.Month == 5 && result.Day == 24
	result.IsMeridian = result.Month == 6 && result.Day == 30
	result.IsHarmony = result.Month == 8 && result.Day == 30
	result.IsMask = result.Month == 8 && result.Day == 13
	result.DayName = dayNames[result.WeekDay-1]
	result.MonthName = monthNames[result.Month-1]
	result.SeasonIndex = (result.Month - 1) / 3

	// Week within the year (each month has 3 weeks of 10 days)
	result.Week = (result.Month-1)*3 + (result.Day-1)/10 + 1

	return result
}

func seasonName(index int, hemisphere string) string {
	if hemisphere == "north" {
		return seasonNamesNorth[index]
	}
	return seasonNamesSouth[index]
}

// ICS date format: YYYYMMDD
func icsDate(t time.Time) string {
	return t.Format("20060102")
}

// ICS timestamp format
func icsTimestamp(t time.Time) string {
	return t.Format("20060102T150405Z")
}

// Fold long ICS lines per RFC 5545 (max 75 octets per line).
func foldLine(line string) string {
	if len(line) <= 75 {
		return line
	}
	var b strings.Builder
	b.WriteString(line[:75])
	line = line[75:]
	for len(line) > 0 {
		end := 74 // continuation lines: space + 74 chars = 75 octets
		if end > len(line) {
			end = len(line)
		}
		b.WriteString("\r\n ")
		b.WriteString(line[:end])
		line = line[end:]
	}
	return b.String()
}

type icsEvent struct {
	uid         string
	dtstart     string // YYYYMMDD
	dtend       string // YYYYMMDD (exclusive, day after)
	summary     string
	description string
	categories  string
}

func formatEvent(e icsEvent, dtstamp string) string {
	var b strings.Builder
	b.WriteString("BEGIN:VEVENT\r\n")
	b.WriteString(foldLine("UID:" + e.uid))
	b.WriteString("\r\n")
	b.WriteString("DTSTAMP:" + dtstamp + "\r\n")
	b.WriteString("DTSTART;VALUE=DATE:" + e.dtstart + "\r\n")
	b.WriteString("DTEND;VALUE=DATE:" + e.dtend + "\r\n")
	b.WriteString(foldLine("SUMMARY:" + e.summary))
	b.WriteString("\r\n")
	if e.description != "" {
		b.WriteString(foldLine("DESCRIPTION:" + e.description))
		b.WriteString("\r\n")
	}
	b.WriteString("TRANSP:TRANSPARENT\r\n")
	if e.categories != "" {
		b.WriteString(foldLine("CATEGORIES:" + e.categories))
		b.WriteString("\r\n")
	}
	b.WriteString("END:VEVENT\r\n")
	return b.String()
}

func nextDay(t time.Time) time.Time {
	return t.AddDate(0, 0, 1)
}

// Season start dates (Gregorian) for a given metric year.
// Season 0 starts at the first regular day (Mar 23 = Unil 1).
// Season 1 starts at Quadril 1 (Jun 21).
// Season 2 starts at Septil 1 (Sep 19).
// Season 3 starts at Decil 1 (Dec 20 or 21 depending on leap).
func seasonStartGregorian(metricYear int) [4]time.Time {
	gregYear := metricYear + 1970
	leap := isGregorianLeapYear(metricYear + 1971)

	// Season 0: Mar 23 of gregYear
	s0 := time.Date(gregYear, 3, 23, 0, 0, 0, 0, time.UTC)
	// Season 1: Jun 21 of gregYear
	s1 := time.Date(gregYear, 6, 21, 0, 0, 0, 0, time.UTC)
	// Season 2: Sep 19 of gregYear
	s2 := time.Date(gregYear, 9, 19, 0, 0, 0, 0, time.UTC)
	// Season 3: After Yule. Yule starts Dec 18 and lasts 2 or 3 days.
	// Yule Eve = Dec 18, Solstice = Dec 19, Kindling = Dec 20 (leap only)
	// Decil 1 = day after Yule ends
	yuleEnd := 20
	if leap {
		yuleEnd = 21
	}
	s3 := time.Date(gregYear, 12, yuleEnd, 0, 0, 0, 0, time.UTC)

	return [4]time.Time{s0, s1, s2, s3}
}

// Month start dates (Gregorian) for a given metric year.
func monthStartGregorian(metricYear int) [12]time.Time {
	gregYear := metricYear + 1970
	leap := isGregorianLeapYear(metricYear + 1971)

	// Month 1 (Unil): Mar 23
	m1 := time.Date(gregYear, 3, 23, 0, 0, 0, 0, time.UTC)

	var months [12]time.Time
	// Months 1-9: each 30 days, starting Mar 23
	for i := 0; i < 9; i++ {
		months[i] = m1.AddDate(0, 0, i*30)
	}

	// After month 9 comes Yule (2 or 3 days)
	yuleDays := 2
	if leap {
		yuleDays = 3
	}
	// Month 10 starts after Yule
	month10Start := m1.AddDate(0, 0, 270+yuleDays)
	for i := 9; i < 12; i++ {
		months[i] = month10Start.AddDate(0, 0, (i-9)*30)
	}

	return months
}

type feedConfig struct {
	hemisphere string
	daily      bool
	filename   string
	calName    string
}

func generateFeed(cfg feedConfig, outDir string, startDate, endDate time.Time, dtstamp string) error {
	var events []string

	for d := startDate; d.Before(endDate); d = nextDay(d) {
		mc := gregorianToMetric(d)
		dateStr := icsDate(d)
		endStr := icsDate(nextDay(d))

		if mc.IsTurning {
			events = append(events, formatEvent(icsEvent{
				uid:        fmt.Sprintf("turning-%s-%s@metricweek.com", mc.SpecialDay, dateStr),
				dtstart:    dateStr,
				dtend:      endStr,
				summary:    "The Turning: " + mc.SpecialDay,
				categories: "Metric Calendar,Special Day",
			}, dtstamp))
			continue
		}

		if mc.IsYule {
			events = append(events, formatEvent(icsEvent{
				uid:        fmt.Sprintf("yule-%s-%s@metricweek.com", strings.ReplaceAll(strings.ToLower(mc.SpecialDay), " ", "-"), dateStr),
				dtstart:    dateStr,
				dtend:      endStr,
				summary:    "Yule: " + mc.SpecialDay,
				categories: "Metric Calendar,Special Day",
			}, dtstamp))
			continue
		}

		// Special days (always included in both lean and daily feeds)
		if mc.IsMidsummer {
			events = append(events, formatEvent(icsEvent{
				uid:        fmt.Sprintf("midsummer-%s@metricweek.com", dateStr),
				dtstart:    dateStr,
				dtend:      endStr,
				summary:    fmt.Sprintf("☀️ Midsummer (%s 1)", monthNames[3]),
				categories: "Metric Calendar,Special Day",
			}, dtstamp))
		}

		if mc.IsSpiral {
			events = append(events, formatEvent(icsEvent{
				uid:        fmt.Sprintf("spiral-%s@metricweek.com", dateStr),
				dtstart:    dateStr,
				dtend:      endStr,
				summary:    fmt.Sprintf("🌀 The Spiral (%s 18)", monthNames[4]),
				categories: "Metric Calendar,Special Day",
			}, dtstamp))
		}

		if mc.IsSextant {
			events = append(events, formatEvent(icsEvent{
				uid:        fmt.Sprintf("sextant-%s@metricweek.com", dateStr),
				dtstart:    dateStr,
				dtend:      endStr,
				summary:    fmt.Sprintf("🧭 The Sextant (%s 30)", monthNames[1]),
				categories: "Metric Calendar,Special Day",
			}, dtstamp))
		}

		if mc.IsTrine {
			events = append(events, formatEvent(icsEvent{
				uid:        fmt.Sprintf("trine-%s@metricweek.com", dateStr),
				dtstart:    dateStr,
				dtend:      endStr,
				summary:    fmt.Sprintf("🔺 The Trine (%s 30)", monthNames[3]),
				categories: "Metric Calendar,Special Day",
			}, dtstamp))
		}

		if mc.IsConvergence {
			events = append(events, formatEvent(icsEvent{
				uid:        fmt.Sprintf("convergence-%s@metricweek.com", dateStr),
				dtstart:    dateStr,
				dtend:      endStr,
				summary:    fmt.Sprintf("🔀 Convergence (%s 24)", monthNames[4]),
				categories: "Metric Calendar,Special Day",
			}, dtstamp))
		}

		if mc.IsMeridian {
			events = append(events, formatEvent(icsEvent{
				uid:        fmt.Sprintf("meridian-%s@metricweek.com", dateStr),
				dtstart:    dateStr,
				dtend:      endStr,
				summary:    fmt.Sprintf("🌗 The Meridian (%s 30)", monthNames[5]),
				categories: "Metric Calendar,Special Day",
			}, dtstamp))
		}

		if mc.IsHarmony {
			events = append(events, formatEvent(icsEvent{
				uid:        fmt.Sprintf("harmony-%s@metricweek.com", dateStr),
				dtstart:    dateStr,
				dtend:      endStr,
				summary:    fmt.Sprintf("🎵 Harmony (%s 30)", monthNames[7]),
				categories: "Metric Calendar,Special Day",
			}, dtstamp))
		}

		if mc.IsMask {
			events = append(events, formatEvent(icsEvent{
				uid:        fmt.Sprintf("mask-%s@metricweek.com", dateStr),
				dtstart:    dateStr,
				dtend:      endStr,
				summary:    fmt.Sprintf("🎭 The Mask (%s 13)", monthNames[7]),
				categories: "Metric Calendar,Special Day",
			}, dtstamp))
		}

		// Season starts
		if mc.Day == 1 && mc.Month%3 == 1 {
			sn := seasonName(mc.SeasonIndex, cfg.hemisphere)
			events = append(events, formatEvent(icsEvent{
				uid:        fmt.Sprintf("season-%s-%s@metricweek.com", strings.ToLower(sn), dateStr),
				dtstart:    dateStr,
				dtend:      endStr,
				summary:    sn + " begins",
				categories: "Metric Calendar,Season",
			}, dtstamp))
		}

		// Month starts
		if mc.Day == 1 {
			events = append(events, formatEvent(icsEvent{
				uid:        fmt.Sprintf("month-%d-%s@metricweek.com", mc.Month, dateStr),
				dtstart:    dateStr,
				dtend:      endStr,
				summary:    fmt.Sprintf("%s begins (Month %d)", mc.MonthName, mc.Month),
				categories: "Metric Calendar,Month",
			}, dtstamp))
		}

		// Daily events
		if cfg.daily {
			sn := seasonName(mc.SeasonIndex, cfg.hemisphere)
			desc := fmt.Sprintf("Week %d · %s · Year %d", mc.Week, sn, mc.Year)
			if mc.IsRest {
				desc += " · The Rest"
			}
			events = append(events, formatEvent(icsEvent{
				uid:         fmt.Sprintf("day-%s@metricweek.com", dateStr),
				dtstart:     dateStr,
				dtend:       endStr,
				summary:     fmt.Sprintf("%s %d · %s", mc.MonthName, mc.Day, mc.DayName),
				description: desc,
				categories:  "Metric Calendar",
			}, dtstamp))
		}
	}

	var b strings.Builder
	b.WriteString("BEGIN:VCALENDAR\r\n")
	b.WriteString("VERSION:2.0\r\n")
	b.WriteString("PRODID:-//Metric Calendar//metricweek.com//EN\r\n")
	b.WriteString(foldLine("X-WR-CALNAME:"+cfg.calName) + "\r\n")
	b.WriteString(foldLine("NAME:"+cfg.calName) + "\r\n")
	b.WriteString("CALSCALE:GREGORIAN\r\n")
	for _, e := range events {
		b.WriteString(e)
	}
	b.WriteString("END:VCALENDAR\r\n")

	outPath := filepath.Join(outDir, cfg.filename)
	return os.WriteFile(outPath, []byte(b.String()), 0644)
}

func main() {
	outDir := "www/calendar"
	if err := os.MkdirAll(outDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "error creating output directory: %v\n", err)
		os.Exit(1)
	}

	dtstamp := icsTimestamp(time.Now().UTC())
	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2036, 3, 23, 0, 0, 0, 0, time.UTC) // through end of metric year 65

	feeds := []feedConfig{
		{hemisphere: "north", daily: false, filename: "metric-north.ics", calName: "Metric Calendar (Northern)"},
		{hemisphere: "south", daily: false, filename: "metric-south.ics", calName: "Metric Calendar (Southern)"},
		{hemisphere: "north", daily: true, filename: "metric-daily-north.ics", calName: "Metric Calendar Daily (Northern)"},
		{hemisphere: "south", daily: true, filename: "metric-daily-south.ics", calName: "Metric Calendar Daily (Southern)"},
	}

	for _, cfg := range feeds {
		if err := generateFeed(cfg, outDir, startDate, endDate, dtstamp); err != nil {
			fmt.Fprintf(os.Stderr, "error generating %s: %v\n", cfg.filename, err)
			os.Exit(1)
		}
		fmt.Printf("Generated %s/%s\n", outDir, cfg.filename)
	}
}
