# metric_calendar

A Python library for the [Metric Calendar](https://metricweek.com) — a rational decimal calendar system built on the March equinox.

The Metric Calendar divides the solar year into 12 months of exactly 30 days each, organised into 10-day weeks. The year begins at the vernal equinox (March 20) with three intercalary Turning days, and includes two or three Yule days at midsummer to keep the calendar aligned with the solar year. Every date is fully deterministic from standard Gregorian input.

## Installation

```bash
pip install metric_calendar
```

Requires Python 3.8 or later. No external dependencies.

## Quick Start

```python
import datetime
from metric_calendar import today, gregorian_to_metric, metric_to_gregorian

# Get today's Metric date
md = today()
print(md.month_name, md.day, md.year)  # e.g. "Unil 1 56"

# Convert a specific Gregorian date
md = gregorian_to_metric(datetime.date(2026, 3, 23))
print(md)
# MetricDate(year=56, month=1, month_name='Unil', day=1, week_day=1,
#            day_name='Primday', week=1, season_index=0, ...)

# Convert back to Gregorian
d = metric_to_gregorian(56, 'month', 1, 1)
print(d)  # 2026-03-23
```

## Calendar Structure

The Metric Calendar year starts at the March equinox (March 20) and is structured as follows:

| Period | Days | Notes |
|---|---|---|
| The Turning | 3 | Vigil, Balance, Dawn — year-opening intercalary days |
| Months 1–9 (Unil–Novil) | 270 | 9 months × 30 days |
| Yule | 2–3 | Yule Eve, Midwinter, and Kindling (leap years only) |
| Months 10–12 (Decil–Duodecil) | 90 | 3 months × 30 days |

Each month contains three 10-day weeks. Days 8, 9, and 10 of each week (Octday, Novday, Decday) are rest days.

**Metric year epoch:** Year 0 = 1970 CE. Year 56 spans March 2026 – March 2027.

### Month Names

| # | Name | # | Name |
|---|---|---|---|
| 1 | Unil | 7 | Septil |
| 2 | Duil | 8 | Octil |
| 3 | Tril | 9 | Novil |
| 4 | Quadril | 10 | Decil |
| 5 | Quintil | 11 | Undecil |
| 6 | Sextil | 12 | Duodecil |

### Day Names

| # | Name | Rest? |
|---|---|---|
| 1 | Primday | |
| 2 | Duoday | |
| 3 | Triday | |
| 4 | Quadday | |
| 5 | Quintday | |
| 6 | Hexday | |
| 7 | Septday | |
| 8 | Octday | rest |
| 9 | Novday | rest |
| 10 | Decday | rest |

### Special Days

| Name | Description |
|---|---|
| Vigil | First Turning day (March 20) |
| Balance | Second Turning day (March 21) |
| Dawn | Third Turning day (March 22) |
| Yule Eve | First Yule day |
| Midwinter | Second Yule day |
| Kindling | Third Yule day (leap years only) |
| Midsummer | Quadril 1 — first day of the fourth month |
| Spiral | Quintil 18 — golden-ratio day of the year |

## API Reference

### `gregorian_to_metric(date: datetime.date) -> MetricDate`

Convert a Gregorian `datetime.date` to a `MetricDate`.

```python
import datetime
from metric_calendar import gregorian_to_metric

md = gregorian_to_metric(datetime.date(2026, 6, 21))
print(md.is_midsummer)   # True
print(md.month_name)     # 'Quadril'
print(md.day)            # 1
```

### `metric_to_gregorian(year, period_type, period_value, day_of_month=1) -> datetime.date`

Convert a Metric Calendar date back to a Gregorian `datetime.date`.

**Parameters:**

- `year` — Metric year (0 = 1970 CE)
- `period_type` — one of `'turning'`, `'month'`, or `'yule'`
- `period_value`:
  - `'turning'`: `0` (Vigil), `1` (Balance), `2` (Dawn)
  - `'month'`: month number `1`–`12`
  - `'yule'`: `0` (Yule Eve), `1` (Midwinter), `2` (Kindling, leap years only)
- `day_of_month` — day within the month, `1`–`30` (only used when `period_type='month'`)

```python
from metric_calendar import metric_to_gregorian

# Unil 1, Year 56
metric_to_gregorian(56, 'month', 1, 1)       # datetime.date(2026, 3, 23)

# The Turning — Vigil, Year 56
metric_to_gregorian(56, 'turning', 0)         # datetime.date(2026, 3, 20)

# Yule Eve, Year 56
metric_to_gregorian(56, 'yule', 0)            # datetime.date(2026, 12, 18)
```

### `is_rest_day(date: datetime.date) -> bool`

Return `True` if the given Gregorian date falls on a Metric rest day (Octday, Novday, or Decday).

```python
from metric_calendar import is_rest_day
import datetime

is_rest_day(datetime.date(2026, 4, 1))   # True  — Decday
is_rest_day(datetime.date(2026, 3, 23))  # False — Primday
```

### `today() -> MetricDate`

Return today's date as a `MetricDate`.

```python
from metric_calendar import today

md = today()
print(f"{md.day_name}, {md.month_name} {md.day}, Year {md.year}")
```

### `MetricDate`

An immutable dataclass with the following fields:

| Field | Type | Description |
|---|---|---|
| `year` | `int` | Metric year (epoch 1970) |
| `month` | `int` | Month number 1–12; `0` for Turning/Yule days |
| `month_name` | `str` | e.g. `'Unil'`; empty for Turning/Yule days |
| `day` | `int` | Day of month 1–30; `0` for Turning/Yule days |
| `week_day` | `int` | Day of week 1–10; `0` for Turning/Yule days |
| `day_name` | `str` | e.g. `'Primday'`; empty for Turning/Yule days |
| `week` | `int` | Week of year 1–36; `0` for Turning/Yule days |
| `season_index` | `int` | Season 0–3 (Spring/Summer/Autumn/Winter); `-1` for Turning/Yule days |
| `is_leap_year` | `bool` | Whether this Metric year has a Kindling day |
| `is_turning` | `bool` | True for Vigil, Balance, or Dawn |
| `is_yule` | `bool` | True for Yule Eve, Midwinter, or Kindling |
| `is_midsummer` | `bool` | True on Quadril 1 |
| `is_spiral` | `bool` | True on Quintil 18 (golden-ratio day) |
| `is_rest` | `bool` | True when `week_day >= 8` |
| `special_day` | `str` | Name of the special day, or empty string |

## Running the Tests

```bash
python -m pytest tests/
# or
python -m unittest discover tests/
```

## License

MIT. See [metricweek.com](https://metricweek.com) for the full calendar specification.
