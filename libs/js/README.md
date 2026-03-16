# metric-calendar

JavaScript/TypeScript library for the [Metric Calendar](https://metricweek.com) — a precision decimal calendar system built on 10-day weeks, 30-day months, and the spring equinox.

The Metric Calendar divides the year into 12 months of exactly 30 days each, with 3 Turning days at the spring equinox and 2–3 Yule days at midwinter. Every month has three 10-day weeks; days 8, 9, and 10 (Octday, Novday, Decday) are rest days. Year zero is 1970 — a nod to the Unix epoch.

## Install

```bash
npm install metric-calendar
```

## Usage

### ESM

```js
import { gregorianToMetric, metricToGregorian, isRestDay, today } from 'metric-calendar'

const mc = gregorianToMetric(new Date('2026-03-23'))
console.log(mc.year)       // 56
console.log(mc.monthName)  // 'Unil'
console.log(mc.day)        // 1
console.log(mc.dayName)    // 'Primday'
console.log(mc.isRest)     // false
```

### CommonJS

```js
const { gregorianToMetric, today } = require('metric-calendar')

const mc = today()
console.log(`${mc.monthName} ${mc.day}, Year ${mc.year}`)
```

### TypeScript

```ts
import { gregorianToMetric, MetricDate } from 'metric-calendar'

const mc: MetricDate = gregorianToMetric(new Date())

if (mc.isTurning) {
  console.log(`The Turning — ${mc.specialDay}`)
} else if (mc.isYule) {
  console.log(`Yule — ${mc.specialDay}`)
} else {
  console.log(`${mc.dayName}, ${mc.monthName} ${mc.day}, Year ${mc.year}`)
}
```

### Browser (IIFE via CDN)

```html
<script src="https://unpkg.com/metric-calendar/dist/metric-calendar.iife.js"></script>
<script>
  const mc = MetricCalendar.today()
  console.log(mc.monthName, mc.day, mc.year)
</script>
```

## API

### `gregorianToMetric(date: Date): MetricDate`

Converts a Gregorian `Date` to a `MetricDate`. Always reads the date in UTC — pass a UTC-normalised date for predictable results:

```js
const d = new Date(Date.UTC(2026, 2, 23))  // 2026-03-23
gregorianToMetric(d)
```

### `metricToGregorian(year, periodType, periodValue, dayOfMonth?): Date`

Converts a Metric Calendar position back to a Gregorian `Date`.

| Parameter | Type | Description |
|-----------|------|-------------|
| `year` | `number` | Metric year (e.g. 56) |
| `periodType` | `'turning' \| 'month' \| 'yule'` | Which period within the year |
| `periodValue` | `number` | Turning: 0–2; month: 1–12; yule: 0–2 |
| `dayOfMonth` | `number` | Day within the month, 1–30 (month only, default 1) |

```js
// First day of Unil, Year 56
metricToGregorian(56, 'month', 1, 1)  // → 2026-03-23

// Vigil (first Turning day), Year 56
metricToGregorian(56, 'turning', 0)   // → 2026-03-20

// Yule Eve, Year 56
metricToGregorian(56, 'yule', 0)      // → 2026-12-18
```

### `isRestDay(date: Date): boolean`

Returns `true` if the date falls on a rest day (Octday, Novday, or Decday — weekDay >= 8).

### `today(): MetricDate`

Returns the current date as a `MetricDate`.

## MetricDate fields

| Field | Type | Description |
|-------|------|-------------|
| `year` | `number` | Metric year. Year 0 = Gregorian 1970. |
| `month` | `number` | 1–12. `0` during Turning or Yule. |
| `monthName` | `string` | e.g. `'Unil'`, `'Quadril'`, `'Decil'`. Empty during Turning/Yule. |
| `day` | `number` | Day within the month, 1–30. `0` during Turning/Yule. |
| `weekDay` | `number` | Position within the 10-day week, 1–10. `0` during Turning/Yule. |
| `dayName` | `string` | e.g. `'Primday'`, `'Hexday'`, `'Decday'`. Empty during Turning/Yule. |
| `week` | `number` | Week of the year, 1–36. `0` during Turning/Yule. |
| `seasonIndex` | `number` | 0 = Spring, 1 = Summer, 2 = Autumn, 3 = Winter. `-1` during Turning/Yule. |
| `isLeapYear` | `boolean` | Whether this Metric year contains a Gregorian leap day (3 Yule days). |
| `isTurning` | `boolean` | The three Turning days at the spring equinox. |
| `isYule` | `boolean` | The two or three Yule days at midwinter. |
| `isMidsummer` | `boolean` | `true` on Quadril 1 (first day of summer). |
| `isSpiral` | `boolean` | `true` on Quintil 18 (the Spiral day). |
| `isRest` | `boolean` | `true` when `weekDay >= 8` (Octday, Novday, Decday). |
| `specialDay` | `string` | Name for Turning days (`'Vigil'`, `'Balance'`, `'Dawn'`) or Yule days (`'Yule Eve'`, `'Midwinter'`, `'Kindling'`). Empty otherwise. |

## Month names

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

## Day names

| weekDay | Name |
|---------|------|
| 1 | Primday |
| 2 | Duoday |
| 3 | Triday |
| 4 | Quadday |
| 5 | Quintday |
| 6 | Hexday |
| 7 | Septday |
| 8 | Octday (rest) |
| 9 | Novday (rest) |
| 10 | Decday (rest) |

## Learn more

Visit [metricweek.com](https://metricweek.com) for a full calendar view, ICS feeds, and an interactive Gregorian-to-Metric converter.

## License

MIT
