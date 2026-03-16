# MetricCalendar

A Ruby library for working with the [Metric Calendar](https://metricweek.com) — a rational decimal calendar system built on the spring equinox, with 10-day weeks and 12 months of 30 days.

## Installation

```sh
gem install MetricCalendar
```

Or add to your `Gemfile`:

```ruby
gem 'MetricCalendar'
```

## Overview

The Metric Calendar anchors Year 0 to the spring equinox of 1970. Each year begins on the spring equinox (March 20) with a three-day intercalary period called **The Turning**, followed by 12 months of 30 days each (arranged in 10-day weeks), a **Yule** intercalary period of 2–3 days near the winter solstice, and then the final 3 months of the year.

| Period | Days | Notes |
|---|---|---|
| The Turning | 3 | Vigil, Balance, Dawn — start of year |
| Months 1–9 | 270 | Spring through early autumn |
| Yule | 2–3 | Yule Eve, Midwinter, Kindling (leap years only) |
| Months 10–12 | 90 | Late autumn through winter |

Weeks run Primday through Decday (days 1–10), with days 8–10 (Octday, Novday, Decday) as rest days.

## Usage

```ruby
require 'metric_calendar'
require 'date'

# Convert today's Gregorian date
md = MetricCalendar.today
puts "#{md.day_name}, #{md.month_name} #{md.day}, Year #{md.year}"
# => "Triday, Decil 13, Year 54"  (for 2025-01-01)

# Convert a specific date
md = MetricCalendar.gregorian_to_metric(Date.new(2026, 3, 23))
puts md.month_name   # => "Unil"
puts md.day          # => 1
puts md.week_day     # => 1
puts md.day_name     # => "Primday"
puts md.week         # => 1
puts md.year         # => 56

# Check special days
md = MetricCalendar.gregorian_to_metric(Date.new(2026, 3, 20))
puts md.is_turning   # => true
puts md.special_day  # => "Vigil"

# Yule
md = MetricCalendar.gregorian_to_metric(Date.new(2026, 12, 18))
puts md.is_yule      # => true
puts md.special_day  # => "Yule Eve"

# Midsummer (Quadril 1) and the Spiral Day (Quintil 18)
md = MetricCalendar.gregorian_to_metric(Date.new(2026, 6, 21))
puts md.is_midsummer # => true

# Rest days (Octday, Novday, Decday)
MetricCalendar.is_rest_day(Date.new(2026, 4, 1))  # => true  (Decday)
MetricCalendar.is_rest_day(Date.new(2026, 3, 23)) # => false (Primday)

# Convert Metric back to Gregorian
MetricCalendar.metric_to_gregorian(56, 'month', 1, 1)
# => #<Date: 2026-03-23>

MetricCalendar.metric_to_gregorian(56, 'turning', 0)
# => #<Date: 2026-03-20>  (Vigil)

MetricCalendar.metric_to_gregorian(56, 'yule', 0)
# => #<Date: 2026-12-18>  (Yule Eve)
```

## API Reference

### `MetricCalendar.gregorian_to_metric(date) -> MetricDate`

Converts a `Date` object to a `MetricDate`.

**Parameters:**
- `date` — a `Date` instance

**Returns:** a `MetricDate` struct (see below).

---

### `MetricCalendar.metric_to_gregorian(year, period_type, period_value, day_of_month = 1) -> Date`

Converts a Metric Calendar date back to a Gregorian `Date`.

**Parameters:**
- `year` — Integer, Metric year (e.g. `56`)
- `period_type` — String: `"turning"`, `"month"`, or `"yule"`
- `period_value` — Integer:
  - `"turning"`: `0` (Vigil), `1` (Balance), `2` (Dawn)
  - `"month"`: `1`–`12`
  - `"yule"`: `0` (Yule Eve), `1` (Midwinter), `2` (Kindling — leap years only)
- `day_of_month` — Integer `1`–`30`, used only when `period_type` is `"month"`

Raises `ArgumentError` for out-of-range values or attempting Kindling in a non-leap year.

---

### `MetricCalendar.is_rest_day(date) -> Boolean`

Returns `true` if the given Gregorian `Date` falls on a rest day (week days 8–10: Octday, Novday, Decday).

---

### `MetricCalendar.today -> MetricDate`

Returns the current date as a `MetricDate`.

---

### `MetricDate` Struct Fields

| Field | Type | Description |
|---|---|---|
| `year` | Integer | Metric year (Year 0 = spring equinox 1970) |
| `month` | Integer | `1`–`12`; `0` during Turning or Yule |
| `month_name` | String | e.g. `"Unil"`, `""` during Turning/Yule |
| `day` | Integer | `1`–`30`; `0` during Turning or Yule |
| `week_day` | Integer | `1`–`10`; `0` during Turning or Yule |
| `day_name` | String | e.g. `"Primday"`, `""` during Turning/Yule |
| `week` | Integer | `1`–`36`; `0` during Turning or Yule |
| `season_index` | Integer | `0`=Spring, `1`=Summer, `2`=Autumn, `3`=Winter; `-1` during Turning/Yule |
| `is_leap_year` | Boolean | `true` if this Metric year has 3 Yule days (Kindling) |
| `is_turning` | Boolean | `true` during The Turning (Vigil, Balance, Dawn) |
| `is_yule` | Boolean | `true` during Yule |
| `is_midsummer` | Boolean | `true` on Quadril 1 (summer solstice) |
| `is_spiral` | Boolean | `true` on Quintil 18 (golden angle day) |
| `is_rest` | Boolean | `true` on week days 8–10 (Octday, Novday, Decday) |
| `special_day` | String | `"Vigil"`, `"Balance"`, `"Dawn"`, `"Yule Eve"`, `"Midwinter"`, `"Kindling"`, or `""` |

### Month Names

| # | Name | Season |
|---|---|---|
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

`Primday`, `Duoday`, `Triday`, `Quadday`, `Quintday`, `Hexday`, `Septday`, `Octday`, `Novday`, `Decday`

Days 8–10 (Octday, Novday, Decday) are rest days.

## Running Tests

```sh
ruby test/test_metric_calendar.rb
```

## Requirements

- Ruby >= 2.5
- No external gems — stdlib `date` only

## License

MIT

## Links

- [metricweek.com](https://metricweek.com) — the Metric Calendar
- [Source](https://github.com/alexrabarts/metric-calendar)
