# Metric Calendar — Java Library

A Java 8 library for converting between Gregorian and [Metric Calendar](https://metricweek.com) dates. Designed for Android legacy compatibility — no `java.time`, no external dependencies beyond JUnit for tests.

## Overview

The Metric Calendar is a rational decimal calendar system with:

- **10-day weeks** (Primday through Decday)
- **12 months of 30 days** (Unil through Duodecil)
- **The Turning** — 3 intercalary days at the spring equinox (March 20–22)
- **Yule** — 2–3 intercalary days at the winter solstice (December 18–20)
- **Year 0** = spring equinox 1970

## Requirements

- Java 8 or higher
- Android API 26+ (or use [desugaring](https://developer.android.com/studio/write/java8-support) for lower APIs)

## Installation

### Gradle (Kotlin DSL)

```kotlin
dependencies {
    implementation("com.metricweek:metric-calendar:1.0.0")
}
```

### Gradle (Groovy DSL)

```groovy
dependencies {
    implementation 'com.metricweek:metric-calendar:1.0.0'
}
```

### Maven

```xml
<dependency>
    <groupId>com.metricweek</groupId>
    <artifactId>metric-calendar</artifactId>
    <version>1.0.0</version>
</dependency>
```

## Building from Source

> **Note:** The Gradle wrapper scripts (`gradlew`, `gradlew.bat`) and `gradle-wrapper.jar` are not included in this repository as they are binary files. Generate them once after cloning:
>
> ```bash
> gradle wrapper
> ```
>
> This requires Gradle 8.7+ to be installed on your system. After running `gradle wrapper`, you can use `./gradlew` for all subsequent commands.

```bash
# Run tests
./gradlew test

# Build JAR
./gradlew build

# Publish to local Maven repository
./gradlew publishToMavenLocal
```

## Usage

### Convert Gregorian to Metric

```java
import com.metricweek.MetricCalendar;
import com.metricweek.MetricDate;

// Regular day
MetricDate date = MetricCalendar.fromGregorian(2026, 3, 23);
System.out.println(date);
// → "Year 56, Unil 1 (Primday)"

System.out.println(date.year);       // 56
System.out.println(date.month);      // 1
System.out.println(date.monthName);  // "Unil"
System.out.println(date.day);        // 1
System.out.println(date.weekDay);    // 1
System.out.println(date.dayName);    // "Primday"
System.out.println(date.week);       // 1
System.out.println(date.seasonIndex); // 0 (spring)
System.out.println(date.isRest);     // false

// The Turning
MetricDate turning = MetricCalendar.fromGregorian(2026, 3, 20);
System.out.println(turning);
// → "Year 56, The Turning — Vigil"
System.out.println(turning.isTurning);   // true
System.out.println(turning.specialDay);  // "Vigil"

// Yule
MetricDate yule = MetricCalendar.fromGregorian(2026, 12, 18);
System.out.println(yule);
// → "Year 56, Yule — Yule Eve"
System.out.println(yule.isYule);        // true
System.out.println(yule.specialDay);    // "Yule Eve"

// Get today's date
MetricDate today = MetricCalendar.today();
```

### Convert Metric to Gregorian

```java
// Month day: Year 56, Unil 1 → 2026-03-23
int[] greg = MetricCalendar.toGregorian(56, "month", 1, 1);
// greg = {2026, 3, 23}

// Turning day: Year 56, Vigil (index 0) → 2026-03-20
int[] vigil = MetricCalendar.toGregorian(56, "turning", 0, 0);
// vigil = {2026, 3, 20}

// Yule day: Year 56, Yule Eve (index 0) → 2026-12-18
int[] yuleEve = MetricCalendar.toGregorian(56, "yule", 0, 0);
// yuleEve = {2026, 12, 18}
```

### Check for Rest Days

```java
// Days 8–10 of any 10-day week are rest days
boolean rest = MetricCalendar.isRestDay(2026, 4, 1);  // true (Unil 10)
boolean work = MetricCalendar.isRestDay(2026, 3, 23); // false (Unil 1)
```

## API Reference

### `MetricCalendar`

Static utility class. All methods are thread-safe.

| Method | Description |
|--------|-------------|
| `fromGregorian(int year, int month, int day)` | Convert a Gregorian date to `MetricDate` |
| `toGregorian(int year, String periodType, int periodValue, int dayOfMonth)` | Convert a Metric date to Gregorian `int[]{year, month, day}` |
| `isRestDay(int year, int month, int day)` | Returns `true` if the Gregorian date falls on a Metric rest day |
| `today()` | Returns today's `MetricDate` using the system clock (UTC) |

#### `toGregorian` parameters

| `periodType` | `periodValue` | `dayOfMonth` |
|---|---|---|
| `"turning"` | `0` = Vigil, `1` = Balance, `2` = Dawn | ignored (pass `0`) |
| `"month"` | `1`–`12` | `1`–`30` |
| `"yule"` | `0` = Yule Eve, `1` = Midwinter, `2` = Kindling (leap years only) | ignored (pass `0`) |

### `MetricDate`

Immutable value object with the following public fields:

| Field | Type | Description |
|-------|------|-------------|
| `year` | `int` | Metric year (Year 0 = spring equinox 1970) |
| `month` | `int` | 1–12 for regular days; `0` for Turning/Yule |
| `monthName` | `String` | e.g. `"Unil"`; empty for Turning/Yule |
| `day` | `int` | 1–30 for regular days; `0` for Turning/Yule |
| `weekDay` | `int` | 1–10 for regular days; `0` for Turning/Yule |
| `dayName` | `String` | e.g. `"Primday"`; empty for Turning/Yule |
| `week` | `int` | 1–36 for regular days; `0` for Turning/Yule |
| `seasonIndex` | `int` | `0`=spring, `1`=summer, `2`=autumn, `3`=winter; `-1` for Turning/Yule |
| `isLeapYear` | `boolean` | `true` if this metric year has 3 Yule days |
| `isTurning` | `boolean` | `true` during The Turning (spring equinox days) |
| `isYule` | `boolean` | `true` during Yule |
| `isMidsummer` | `boolean` | `true` on Quadril 1 (summer solstice) |
| `isSpiral` | `boolean` | `true` on Quintil 18 (golden angle day) |
| `isRest` | `boolean` | `true` on days 8–10 of any 10-day week |
| `specialDay` | `String` | `"Vigil"`, `"Balance"`, `"Dawn"`, `"Yule Eve"`, `"Midwinter"`, `"Kindling"`, or `""` |

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
| 8 | Octday (rest) |
| 9 | Novday (rest) |
| 10 | Decday (rest) |

## Year Numbering

A metric year begins at the spring equinox (March 20). Dates before the equinox in a given Gregorian year belong to the previous metric year:

- 2025-01-01 → **Year 54** (before the March 2025 equinox; 2024 − 1970 = 54)
- 2025-03-20 → **Year 55** (on the equinox; 2025 − 1970 = 55)
- 2026-03-20 → **Year 56**

## License

MIT License — see [LICENSE](https://opensource.org/licenses/MIT)

---

[metricweek.com](https://metricweek.com)
