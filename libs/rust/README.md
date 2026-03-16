# metric_calendar

A Rust library for the [Metric Calendar](https://metricweek.com) — a rational decimal calendar system built on natural astronomical cycles.

The Metric Calendar divides the solar year into 12 months of 30 days each, grouped into 3-month seasons. Each month contains three 10-day weeks. The year begins at the spring equinox (March 20), marked by a 3-day intercalary period called The Turning. At midwinter, 2 or 3 Yule days bridge the transition into the new year's final months.

**Year 0** corresponds to the spring equinox of 1970.

## Install

```toml
[dependencies]
metric_calendar = "1.0.0"
```

Or via cargo:

```sh
cargo add metric_calendar
```

Zero external dependencies — uses only `std`.

## Usage

### Convert today's date

```rust
use metric_calendar::today;

let date = today();
println!("Metric date: Year {}, {} {}", date.year, date.month_name, date.day);
// e.g. "Metric date: Year 56, Unil 7"
```

### Convert a Gregorian date to Metric

```rust
use metric_calendar::gregorian_to_metric;

let date = gregorian_to_metric(2026, 3, 23);
assert_eq!(date.year, 56);
assert_eq!(date.month_name, "Unil");
assert_eq!(date.day, 1);
assert_eq!(date.day_name, "Primday");
assert_eq!(date.week, 1);
```

### Special calendar periods

```rust
use metric_calendar::gregorian_to_metric;

// The Turning — 3 intercalary days at the spring equinox
let vigil = gregorian_to_metric(2026, 3, 20);
assert!(vigil.is_turning);
assert_eq!(vigil.special_day, "Vigil");

// Yule — midwinter intercalary days
let yule_eve = gregorian_to_metric(2026, 12, 18);
assert!(yule_eve.is_yule);
assert_eq!(yule_eve.special_day, "Yule Eve");

// Midsummer — Quadril 1, the summer solstice anchor
let midsummer = gregorian_to_metric(2026, 6, 21);
assert!(midsummer.is_midsummer);
```

### Check rest days

Rest days fall on days 8, 9, and 10 of every 10-day week.

```rust
use metric_calendar::is_rest_day;

assert!(is_rest_day(2026, 4, 1));   // Unil 10 — a rest day
assert!(!is_rest_day(2026, 3, 23)); // Unil 1 — not a rest day
```

### Convert a Metric date back to Gregorian

```rust
use metric_calendar::metric_to_gregorian;

// Regular month day
let (y, m, d) = metric_to_gregorian(56, "month", 1, 1).unwrap();
assert_eq!((y, m, d), (2026, 3, 23));

// Turning day (0 = Vigil, 1 = Balance, 2 = Dawn)
let (y, m, d) = metric_to_gregorian(56, "turning", 0, 0).unwrap();
assert_eq!((y, m, d), (2026, 3, 20));

// Yule day (0 = Yule Eve, 1 = Midwinter, 2 = Kindling — leap years only)
let (y, m, d) = metric_to_gregorian(56, "yule", 0, 0).unwrap();
assert_eq!((y, m, d), (2026, 12, 18));
```

## API Reference

### Types

#### `MetricDate`

```rust
pub struct MetricDate {
    pub year: i32,           // Metric year (0 = spring equinox 1970)
    pub month: u8,           // 1–12 for regular days, 0 for Turning/Yule
    pub month_name: &'static str, // e.g. "Unil", "Quadril"
    pub day: u8,             // 1–30 for regular days, 0 for Turning/Yule
    pub week_day: u8,        // 1–10 within the 10-day week, 0 for special days
    pub day_name: &'static str,   // e.g. "Primday", "Decday"
    pub week: u8,            // 1–36 within the year, 0 for special days
    pub season_index: i8,    // 0–3 (Spring/Summer/Autumn/Winter), -1 for special days
    pub is_leap_year: bool,  // true if the metric year has a Kindling day
    pub is_turning: bool,    // true during The Turning (spring equinox)
    pub is_yule: bool,       // true during Yule (midwinter)
    pub is_midsummer: bool,  // true on Quadril 1 (summer solstice anchor)
    pub is_spiral: bool,     // true on Quintil 18 (golden angle / spiral day)
    pub is_rest: bool,       // true on week days 8, 9, 10
    pub special_day: &'static str, // "Vigil", "Balance", "Dawn", "Yule Eve",
                             //   "Midwinter", "Kindling", or "" for regular days
}
```

### Functions

#### `gregorian_to_metric(year: i32, month: u32, day: u32) -> MetricDate`

Converts a Gregorian calendar date to a `MetricDate`. Accepts any valid proleptic Gregorian date.

#### `metric_to_gregorian(year: i32, period_type: &str, period_value: i32, day_of_month: u32) -> Result<(i32, u32, u32), String>`

Converts a Metric Calendar date back to a Gregorian `(year, month, day)` tuple.

- `period_type`: `"turning"`, `"month"`, or `"yule"`
- `period_value`:
  - For `"turning"`: `0` (Vigil), `1` (Balance), `2` (Dawn)
  - For `"month"`: `1`–`12`
  - For `"yule"`: `0` (Yule Eve), `1` (Midwinter), `2` (Kindling — leap years only)
- `day_of_month`: `1`–`30`, only used when `period_type` is `"month"`

Returns `Err` for invalid inputs (e.g. Kindling in a non-leap year, out-of-range values).

#### `is_rest_day(year: i32, month: u32, day: u32) -> bool`

Returns `true` if the Gregorian date falls on a Metric rest day (week days 8, 9, or 10).

#### `today() -> MetricDate`

Returns the current Metric Calendar date derived from the system clock (UTC).

## Calendar Structure

| Period | Days | Description |
|--------|------|-------------|
| The Turning | 3 | Vigil, Balance, Dawn — spring equinox |
| Months 1–9 (Unil–Novil) | 270 | 9 months × 30 days, Spring through early Winter |
| Yule | 2 or 3 | Yule Eve, Midwinter, Kindling (leap years) |
| Months 10–12 (Decil–Duodecil) | 90 | 3 months × 30 days, late Winter |

**Total**: 365 days (366 in leap years)

### Month Names

| # | Name | Season |
|---|------|--------|
| 1 | Unil | Spring |
| 2 | Duil | Spring |
| 3 | Tril | Spring |
| 4 | Quadril | Summer |
| 5 | Quintil | Summer |
| 6 | Sextil | Summer |
| 7 | Septil | Autumn |
| 8 | Octil | Autumn |
| 9 | Novil | Autumn |
| 10 | Decil | Winter |
| 11 | Undecil | Winter |
| 12 | Duodecil | Winter |

### Day Names

| # | Name |
|---|------|
| 1 | Primday |
| 2 | Duoday |
| 3 | Triday |
| 4 | Quadday |
| 5 | Quintday |
| 6 | Hexday |
| 7 | Septday |
| 8 | Octday |
| 9 | Novday |
| 10 | Decday |

Days 8–10 (Octday, Novday, Decday) are rest days in every week.

## License

MIT

---

Learn more at [metricweek.com](https://metricweek.com).
