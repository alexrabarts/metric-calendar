# MetricCalendar

A Swift Package Manager library for converting between the [Metric Calendar](https://metricweek.com) and the Gregorian calendar.

The Metric Calendar is a rational reform calendar that begins each year at the spring equinox (March 20). It consists of 12 months of 30 days each, a 3-day transitional period called The Turning at the new year, and a 2–3 day mid-winter Yule period. Each 30-day month is divided into three 10-day weeks.

## Installation

Add the package to your `Package.swift` dependencies:

```swift
dependencies: [
    .package(url: "https://github.com/alexrabarts/metric-calendar", from: "1.0.0"),
],
```

Then add it as a dependency to your target:

```swift
.target(
    name: "MyApp",
    dependencies: ["MetricCalendar"]
),
```

## Usage

### Convert today's date

```swift
import MetricCalendar

let metricToday = today()
print(metricToday)
// e.g. "Year 56, Unil 1 (Primday)"
```

### Convert a specific Gregorian date to Metric

```swift
import Foundation
import MetricCalendar

// Dates are interpreted as UTC
var cal = Calendar(identifier: .gregorian)
cal.timeZone = TimeZone(identifier: "UTC")!
var comps = DateComponents()
comps.year = 2025; comps.month = 1; comps.day = 1
let date = cal.date(from: comps)!

let metricDate = gregorianToMetric(date)
print(metricDate.year)       // 54  (before March 2025 equinox)
print(metricDate.monthName)  // "Decil"
print(metricDate.day)        // 13
print(metricDate.dayName)    // "Triday"
```

### Convert a Metric date back to Gregorian

```swift
import MetricCalendar

// Regular month day: Year 56, Unil 1 → 2026-03-23
let date = try metricToGregorian(year: 56, periodType: "month", periodValue: 1, dayOfMonth: 1)

// The Turning (Vigil = 0, Balance = 1, Dawn = 2): Year 56, Vigil → 2026-03-20
let vigil = try metricToGregorian(year: 56, periodType: "turning", periodValue: 0)

// Yule (Yule Eve = 0, Midwinter = 1, Kindling = 2*): Year 56, Yule Eve → 2026-12-18
// * Kindling only occurs in leap years
let yuleEve = try metricToGregorian(year: 56, periodType: "yule", periodValue: 0)
```

### Check if a date is a rest day

```swift
import MetricCalendar

let date = Date()
if isRestDay(date) {
    print("Today is a rest day (days 8–10 of a 10-day week)")
}
```

### Inspect a MetricDate

```swift
let r = gregorianToMetric(Date())

r.year          // Int  — metric year (Year 0 = spring equinox 1970)
r.month         // Int  — 1–12 for regular days, 0 for Turning/Yule
r.monthName     // String — e.g. "Unil", "Quadril", "Duodecil"
r.day           // Int  — 1–30 for regular days, 0 for Turning/Yule
r.weekDay       // Int  — 1–10, 0 for Turning/Yule
r.dayName       // String — e.g. "Primday", "Decday"
r.week          // Int  — 1–36 for regular days, 0 for Turning/Yule
r.seasonIndex   // Int  — 0 (spring), 1 (summer), 2 (autumn), 3 (winter); -1 for Turning/Yule
r.isLeapYear    // Bool — true if this metric year has a Kindling day
r.isTurning     // Bool — true during The Turning (3 days at spring equinox)
r.isYule        // Bool — true during Yule (mid-winter)
r.isMidsummer   // Bool — true on Quadril 1 (summer solstice, ~June 21)
r.isSpiral      // Bool — true on Quintil 18 (golden angle day)
r.isRest        // Bool — true on days 8–10 of any 10-day week
r.specialDay    // String — "Vigil", "Balance", "Dawn", "Yule Eve", "Midwinter", "Kindling", or ""
r.description   // String — human-readable representation
```

## API Reference

### Types

#### `MetricDate`

A `struct` representing a date in the Metric Calendar. Conforms to `Equatable`, `Hashable`, and `CustomStringConvertible`.

| Property | Type | Description |
|---|---|---|
| `year` | `Int` | Metric year. Year 0 begins at the spring equinox of 1970. |
| `month` | `Int` | Month number 1–12 for regular days; 0 for Turning/Yule. |
| `monthName` | `String` | Month name (Unil through Duodecil), or `""` for Turning/Yule. |
| `day` | `Int` | Day of month, 1–30; 0 for Turning/Yule. |
| `weekDay` | `Int` | Day within the 10-day week, 1–10; 0 for Turning/Yule. |
| `dayName` | `String` | Day name (Primday through Decday), or `""` for Turning/Yule. |
| `week` | `Int` | Week number within the year, 1–36; 0 for Turning/Yule. |
| `seasonIndex` | `Int` | Season: 0=spring, 1=summer, 2=autumn, 3=winter; -1 for Turning/Yule. |
| `isLeapYear` | `Bool` | Whether this metric year includes a Kindling day. |
| `isTurning` | `Bool` | Whether this is one of the 3 Turning days at the spring equinox. |
| `isYule` | `Bool` | Whether this is a Yule day (mid-winter). |
| `isMidsummer` | `Bool` | Whether this is Quadril 1 (summer solstice). |
| `isSpiral` | `Bool` | Whether this is Quintil 18 (the golden angle day). |
| `isRest` | `Bool` | Whether this is a rest day (days 8–10 of any 10-day week). |
| `specialDay` | `String` | Special day name, or `""` for ordinary days. |
| `description` | `String` | Human-readable description, e.g. `"Year 56, Unil 1 (Primday)"`. |

#### `MetricCalendarError`

An `enum` conforming to `Error` and `LocalizedError` with cases:

| Case | Description |
|---|---|
| `invalidTurningValue` | Turning period value was not in 0–2. |
| `invalidMonth` | Month was not in 1–12. |
| `invalidDay` | Day was not in 1–30. |
| `invalidYuleValue` | Yule period value was not in 0–2. |
| `kindlingRequiresLeapYear` | Kindling (`yule` value 2) was requested for a non-leap metric year. |
| `unknownPeriodType(String)` | `periodType` was not `"turning"`, `"month"`, or `"yule"`. |

### Functions

#### `gregorianToMetric(_ date: Date) -> MetricDate`

Converts a Gregorian `Date` to a `MetricDate`. The date is interpreted as UTC; time-of-day is ignored.

#### `metricToGregorian(year:periodType:periodValue:dayOfMonth:) throws -> Date`

Converts a Metric Calendar date to a Gregorian `Date` (UTC midnight).

| Parameter | Type | Description |
|---|---|---|
| `year` | `Int` | Metric year number. |
| `periodType` | `String` | `"turning"`, `"month"`, or `"yule"`. |
| `periodValue` | `Int` | For turning: 0 (Vigil), 1 (Balance), 2 (Dawn). For month: 1–12. For yule: 0 (Yule Eve), 1 (Midwinter), 2 (Kindling). |
| `dayOfMonth` | `Int` | Day within the month (1–30). Only used when `periodType` is `"month"`. Defaults to `1`. |

Throws `MetricCalendarError` on invalid input.

#### `isRestDay(_ date: Date) -> Bool`

Returns `true` if the given Gregorian date falls on days 8–10 of a Metric 10-day week.

#### `today() -> MetricDate`

Returns the current Metric Calendar date, using the UTC system clock.

## Calendar Structure

| Period | Days | Notes |
|---|---|---|
| The Turning | 3 | Vigil, Balance, Dawn — spring equinox (March 20) |
| Unil – Novil | 270 | Months 1–9, 30 days each |
| Yule | 2–3 | Yule Eve, Midwinter, [Kindling in leap years] |
| Decil – Duodecil | 90 | Months 10–12, 30 days each |

Each month has three 10-day weeks. Rest days are the 8th, 9th, and 10th day of each week (Octday, Novday, Decday).

**Month names:** Unil, Duil, Tril, Quadril, Quintil, Sextil, Septil, Octil, Novil, Decil, Undecil, Duodecil

**Day names:** Primday, Duoday, Triday, Quadday, Quintday, Hexday, Septday, Octday, Novday, Decday

## Requirements

- Swift 5.9+
- macOS 10.15+ / iOS 13+ / watchOS 6+ / tvOS 13+
- Zero external dependencies (Foundation only)

## Learn More

Visit [metricweek.com](https://metricweek.com) for the full calendar, current date, and more information about the Metric Calendar system.
