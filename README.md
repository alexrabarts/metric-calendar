# Metric Calendar

A rational decimal calendar system built on 10-day weeks, 30-day months, and the spring equinox.

**Website:** [metricweek.com](https://metricweek.com)

---

## The calendar

The Metric Calendar divides the solar year into 12 months of exactly 30 days each, grouped into four 3-month seasons. Each month contains three 10-day weeks. Days 8, 9, and 10 of each week (Octday, Novday, Decday) are rest days.

Two intercalary periods sit outside the regular months:

- **The Turning** — 3 days at the spring equinox (March 20–22), marking the year transition: Vigil, Balance, Dawn.
- **Yule** — 2 days (3 in leap years) at midwinter (December 18–20), bridging the transition into the year's final three months: Yule Eve, Midwinter, Kindling (leap years only).

**Year 0** is the spring equinox of 1970 — the Unix epoch.

360 regular days. The same number as degrees in a circle. Special solar and geometric landmarks (Midsummer, The Spiral, The Trine, The Meridian, and others) fall on exact calendar dates each year without being placed there — they emerge from the arithmetic.

---

## Libraries

Zero-dependency libraries are available for 8 languages, all at v1.0.0:

| Language | Install |
|---|---|
| **Go** | `go get github.com/alexrabarts/metric-calendar/libs/go` |
| **TypeScript / JavaScript** | `npm install metric-calendar` |
| **Python** | `pip install metric_calendar` |
| **Rust** | `cargo add metric_calendar` |
| **Ruby** | `gem install metric_calendar` |
| **Swift** | Swift Package Manager — see `libs/swift/` |
| **Kotlin / JVM** | Gradle / Maven — see `libs/kotlin/` |
| **Java** | Gradle / Maven — see `libs/java/` |

Each library lives under `libs/<language>/` and has its own README with full API docs and examples.

### Quick example (Go)

```go
import metric "github.com/alexrabarts/metric-calendar/libs/go"

d := metric.Today()
fmt.Printf("%s %d, Year %d\n", d.MonthName, d.Day, d.Year)
// e.g. "Unil 4, Year 56"
```

### Quick example (TypeScript)

```ts
import { today } from 'metric-calendar'

const d = today()
console.log(`${d.monthName} ${d.day}, Year ${d.year}`)
// e.g. "Unil 4, Year 56"
```

### Quick example (Python)

```python
from metric_calendar import today

d = today()
print(f"{d.month_name} {d.day}, Year {d.year}")
# e.g. "Unil 4, Year 56"
```

---

## Repository layout

```
libs/           Language libraries (Go, JS/TS, Python, Rust, Ruby, Swift, Kotlin, Java)
www/            Website source (metricweek.com)
  index.html    Single-page app
  metric-calendar.js  Browser bundle (built from libs/js)
cmd/
  generate-ics/ Go CLI that generates the ICS calendar feeds
justfile        Common development tasks
```

---

## Local development

Prerequisites: [just](https://just.systems), Node.js, Go 1.23+, Python 3.

```bash
# Serve the website locally
just serve          # http://localhost:8000

# Rebuild the JS library and copy the browser bundle into www/
just build-js

# Regenerate the ICS calendar feeds into www/calendar/
just generate-calendar
```

The website (`www/index.html`) loads `metric-calendar.js` as an IIFE and uses it for all date conversion logic. If you change `libs/js/src/`, run `just build-js` to update the bundle.

---

## Deploying

```bash
just deploy
```

Runs `build-js` and `generate-calendar`, then rsyncs `www/` to the production server and reloads Caddy.

---

## License

MIT — see [LICENSE](LICENSE).
