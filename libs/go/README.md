# metric — Go library for the Metric Calendar

A Go library for working with the Metric Calendar — a precision decimal calendar system built on the spring equinox, ten-day weeks, and thirty-day months.

The Metric Calendar year begins on the spring equinox (March 20) with three **Turning Days** (Vigil, Balance, Dawn), followed by nine regular months of thirty days each, a **Yule** intercalary period of two or three days, and then three final months of thirty days. Every date aligns with a whole number of ten-day weeks, and special solar anchor points — Midsummer and Spiral — fall on predictable calendar positions each year.

Year 0 corresponds to the spring equinox of 1970.

More information: https://metricweek.com

---

## Installation

```bash
go get github.com/alexrabarts/metric-calendar/libs/go
```

Import as:

```go
import metric "github.com/alexrabarts/metric-calendar/libs/go"
```

---

## Quick start

### Convert today's date

```go
package main

import (
    "fmt"
    metric "github.com/alexrabarts/metric-calendar/libs/go"
    "time"
)

func main() {
    d := metric.FromGregorian(time.Now())

    switch {
    case d.IsTurning:
        fmt.Printf("Metric Year %d — %s (Turning Day)\n", d.Year, d.SpecialDay)
    case d.IsYule:
        fmt.Printf("Metric Year %d — %s (Yule)\n", d.Year, d.SpecialDay)
    default:
        fmt.Printf("Metric Year %d — %s %d, %s (Week %d)\n",
            d.Year, d.MonthName, d.Day, d.DayName, d.Week)
        if d.IsMidsummer {
            fmt.Println("  (Midsummer)")
        }
        if d.IsSpiral {
            fmt.Println("  (Spiral Day)")
        }
        if d.IsRest {
            fmt.Println("  (Rest day)")
        }
    }
}
```

### Convert a specific date

```go
d := metric.FromGregorian(time.Date(2026, 3, 23, 0, 0, 0, 0, time.UTC))
// d.Year      = 56
// d.Month     = 1
// d.MonthName = "Unil"
// d.Day       = 1
// d.WeekDay   = 1
// d.DayName   = "Primday"
// d.Week      = 1
```

### Reverse conversion: Metric to Gregorian

```go
// First day of Unil, Metric Year 56
t, err := metric.ToGregorian(56, "month", 1, 1)
// t == 2026-03-23 (UTC)

// Turning Day: Vigil
t, err = metric.ToGregorian(56, "turning", 0, 0)
// t == 2026-03-20 (UTC)

// Yule Eve
t, err = metric.ToGregorian(56, "yule", 0, 0)
// t == 2026-12-18 (UTC)
```

### Check rest days

```go
if metric.IsRestDay(time.Now()) {
    fmt.Println("Today is a rest day (Octday, Novday, or Decday).")
}
```

### Today shorthand

```go
d := metric.Today()
```

---

## API reference

### `FromGregorian(t time.Time) Date`

Converts a Gregorian `time.Time` to a Metric Calendar `Date`. The time-of-day component is ignored; only the calendar date is used.

### `ToGregorian(year int, periodType string, periodValue int, dayOfMonth int) (time.Time, error)`

Converts a Metric Calendar date back to a Gregorian `time.Time` (UTC midnight).

| `periodType` | `periodValue` | `dayOfMonth` |
|---|---|---|
| `"turning"` | 0=Vigil, 1=Balance, 2=Dawn | ignored |
| `"month"` | 1–12 | 1–30 |
| `"yule"` | 0=Yule Eve, 1=Midwinter, 2=Kindling | ignored |

Kindling (yule day 2) only occurs in leap years; passing `periodValue=2` in a non-leap year returns an error.

### `IsRestDay(t time.Time) bool`

Reports whether the given Gregorian date falls on a Metric rest day (WeekDay ≥ 8: Octday, Novday, or Decday).

### `Today() Date`

Returns the Metric Calendar Date for the current local date.

---

## The `Date` type

| Field | Type | Description |
|---|---|---|
| `Year` | `int` | Metric year (0 = spring equinox 1970) |
| `Month` | `int` | Month number 1–12; 0 on Turning/Yule days |
| `MonthName` | `string` | e.g. `"Unil"`, `"Quadril"` |
| `Day` | `int` | Day of month 1–30; 0 on Turning/Yule days |
| `WeekDay` | `int` | Day of ten-day week 1–10; 0 on Turning/Yule days |
| `DayName` | `string` | e.g. `"Primday"`, `"Decday"` |
| `Week` | `int` | Week of year 1–36; 0 on Turning/Yule days |
| `SeasonIndex` | `int` | 0–3 (spring/summer/autumn/winter); -1 on Turning/Yule days |
| `IsLeapYear` | `bool` | Whether the Metric year has a third Yule day (Kindling) |
| `IsTurning` | `bool` | True for the three Turning Days at the year's start |
| `IsYule` | `bool` | True for Yule intercalary days |
| `IsMidsummer` | `bool` | True on Month 4, Day 1 (summer solstice anchor) |
| `IsSpiral` | `bool` | True on Month 5, Day 18 (golden ratio solar marker) |
| `IsRest` | `bool` | True when WeekDay ≥ 8 (Octday, Novday, Decday) |
| `SpecialDay` | `string` | Name of the special day for Turning/Yule dates |

### Month names

| # | Name | # | Name |
|---|---|---|---|
| 1 | Unil | 7 | Septil |
| 2 | Duil | 8 | Octil |
| 3 | Tril | 9 | Novil |
| 4 | Quadril | 10 | Decil |
| 5 | Quintil | 11 | Undecil |
| 6 | Sextil | 12 | Duodecil |

### Day names

| WeekDay | Name | WeekDay | Name |
|---|---|---|---|
| 1 | Primday | 6 | Hexday |
| 2 | Duoday | 7 | Septday |
| 3 | Triday | 8 | Octday* |
| 4 | Quadday | 9 | Novday* |
| 5 | Quintday | 10 | Decday* |

\* Rest days (`IsRest = true`)

---

## License

See the repository root for license information.
