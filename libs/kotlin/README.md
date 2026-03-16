# Metric Calendar — Kotlin Library

A Kotlin/JVM library for converting between Gregorian dates and the [Metric Calendar](https://metricweek.com) system — a precision decimal calendar anchored to the spring equinox.

The Metric Calendar divides the year into:
- **The Turning** — 3 intercalary days at the spring equinox (Vigil, Balance, Dawn)
- **12 months** of 30 days each, organised into 3 ten-day weeks per month (36 weeks total)
- **Yule** — 2 intercalary days at midwinter (3 in a leap year: Yule Eve, Midwinter, Kindling)

Year 0 begins at the spring equinox of 1970.

---

## Requirements

- JVM 8+
- Gradle 8.7+ (see [Setup](#setup) below)

---

## Setup

### Generate the Gradle wrapper scripts

This repository includes `gradle/wrapper/gradle-wrapper.properties` but not the binary wrapper scripts (`gradlew`, `gradlew.bat`, `gradle-wrapper.jar`). Before building for the first time, generate them with:

```bash
gradle wrapper --gradle-version 8.7
```

This requires Gradle to be installed locally (e.g. `brew install gradle` on macOS). After running the command, you can use `./gradlew` for all subsequent builds.

---

## Building

```bash
./gradlew build
```

### Running tests

```bash
./gradlew test
```

### Publishing to local Maven repository

```bash
./gradlew publishToMavenLocal
```

---

## Installation

### From local Maven (after `publishToMavenLocal`)

Add to your `build.gradle.kts`:

```kotlin
repositories {
    mavenLocal()
    mavenCentral()
}

dependencies {
    implementation("com.metricweek:metric-calendar:1.0.0")
}
```

### As a source dependency (Git submodule or included build)

```kotlin
// settings.gradle.kts
includeBuild("path/to/metric-calendar")
```

---

## Usage

```kotlin
import com.metricweek.MetricCalendar

// Convert a Gregorian date to Metric
val date = MetricCalendar.gregorianToMetric(2026, 3, 23)
println(date.year)       // 56
println(date.monthName)  // "Unil"
println(date.day)        // 1
println(date.dayName)    // "Primday"
println(date.week)       // 1
println(date.isRest)     // false

// The Turning (spring equinox days)
val vigil = MetricCalendar.gregorianToMetric(2026, 3, 20)
println(vigil.isTurning)   // true
println(vigil.specialDay)  // "Vigil"

// Yule
val yuleEve = MetricCalendar.gregorianToMetric(2026, 12, 18)
println(yuleEve.isYule)     // true
println(yuleEve.specialDay) // "Yule Eve"

// Convert back to Gregorian
val gregorian = MetricCalendar.metricToGregorian(56, "month", 1, 1)
println(gregorian) // 2026-03-23

val turning = MetricCalendar.metricToGregorian(56, "turning", 0)
println(turning) // 2026-03-20

val yule = MetricCalendar.metricToGregorian(56, "yule", 0)
println(yule) // 2026-12-18

// Check if today is a rest day
println(MetricCalendar.isRestDay(2026, 4, 1)) // true  (Decday)
println(MetricCalendar.isRestDay(2026, 3, 23)) // false (Primday)

// Get today's Metric date
val today = MetricCalendar.today()
println("Today is ${today.dayName}, ${today.monthName} ${today.day}, Year ${today.year}")
```

---

## API Reference

### `MetricCalendar` (object)

#### `gregorianToMetric(year: Int, month: Int, day: Int): MetricDate`

Converts a Gregorian date to a `MetricDate`.

| Parameter | Description |
|-----------|-------------|
| `year` | Gregorian year (e.g. 2026) |
| `month` | Gregorian month, 1–12 |
| `day` | Gregorian day of month, 1–31 |

#### `metricToGregorian(year: Int, periodType: String, periodValue: Int, dayOfMonth: Int = 1): LocalDate`

Converts a Metric Calendar date back to a Gregorian `LocalDate`.

| Parameter | Description |
|-----------|-------------|
| `year` | Metric year |
| `periodType` | `"turning"`, `"month"`, or `"yule"` |
| `periodValue` | Turning: 0–2; Month: 1–12; Yule: 0–2 |
| `dayOfMonth` | Day within month, 1–30 (only used when `periodType = "month"`) |

Throws `IllegalArgumentException` if inputs are invalid (e.g. requesting Kindling in a non-leap year).

#### `isRestDay(year: Int, month: Int, day: Int): Boolean`

Returns `true` if the given Gregorian date falls on a rest day (days 8–10 of any 10-day week).

#### `today(): MetricDate`

Returns the current date as a `MetricDate` using the system clock.

---

### `MetricDate` (data class)

| Field | Type | Description |
|-------|------|-------------|
| `year` | `Int` | Metric year (0 = equinox of 1970) |
| `month` | `Int` | 1–12 for regular days; 0 for Turning/Yule |
| `monthName` | `String` | e.g. `"Unil"`; empty for Turning/Yule |
| `day` | `Int` | 1–30 for regular days; 0 for Turning/Yule |
| `weekDay` | `Int` | 1–10 for regular days; 0 for Turning/Yule |
| `dayName` | `String` | e.g. `"Primday"`; empty for Turning/Yule |
| `week` | `Int` | 1–36 for regular days; 0 for Turning/Yule |
| `seasonIndex` | `Int` | 0–3 (spring/summer/autumn/winter); -1 for Turning/Yule |
| `isLeapYear` | `Boolean` | `true` if this metric year has 3 Yule days |
| `isTurning` | `Boolean` | `true` during The Turning (3 days at spring equinox) |
| `isYule` | `Boolean` | `true` during Yule |
| `isMidsummer` | `Boolean` | `true` on Quadril 1 (summer solstice, ~June 21) |
| `isSpiral` | `Boolean` | `true` on Quintil 18 (golden angle day) |
| `isRest` | `Boolean` | `true` on days 8–10 of any 10-day week |
| `specialDay` | `String` | `"Vigil"`, `"Balance"`, `"Dawn"`, `"Yule Eve"`, `"Midwinter"`, `"Kindling"`, or empty |

#### Month names

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

#### Day names (10-day week)

| # | Name |
|---|------|
| 1 | Primday |
| 2 | Duoday |
| 3 | Triday |
| 4 | Quadday |
| 5 | Quintday |
| 6 | Hexday |
| 7 | Septday |
| 8 | Octday *(rest)* |
| 9 | Novday *(rest)* |
| 10 | Decday *(rest)* |

---

## Year numbering

Metric year = Gregorian year − 1970, counted from the spring equinox (March 20).

Dates before the equinox in a given Gregorian year belong to the **previous** metric year. For example, January 1 2025 falls before the March 20 2025 equinox, so it is in Metric Year 54 (2024 − 1970).

---

## License

MIT — see [LICENSE](https://opensource.org/licenses/MIT).

---

## Links

- [metricweek.com](https://metricweek.com)
