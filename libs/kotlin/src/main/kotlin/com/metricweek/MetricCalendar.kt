package com.metricweek

import java.time.LocalDate
import java.time.temporal.ChronoUnit

private val DAY_NAMES = arrayOf(
    "Primday", "Duoday", "Triday", "Quadday", "Quintday",
    "Hexday", "Septday", "Octday", "Novday", "Decday"
)

private val MONTH_NAMES = arrayOf(
    "Unil", "Duil", "Tril", "Quadril", "Quintil", "Sextil",
    "Septil", "Octil", "Novil", "Decil", "Undecil", "Duodecil"
)

private val SEASON_NAMES = arrayOf("Rising", "Flourishing", "Gathering", "Stillness")

private val TURNING_DAY_NAMES = arrayOf("Vigil", "Balance", "Dawn")
private val YULE_DAY_NAMES = arrayOf("Yule Eve", "Midwinter", "Kindling")

private val FORMAT_RE = Regex("MMM|MM|M|DD|D|WW|W|Y|S")

/**
 * A date in the Metric Calendar system.
 *
 * @param year Metric year (Year 0 = spring equinox 1970)
 * @param month 1-12 for regular days, 0 for Turning/Yule
 * @param monthName e.g. "Unil", empty for Turning/Yule
 * @param day 1-30 for regular days, 0 for Turning/Yule
 * @param weekDay 1-10 for regular days, 0 for Turning/Yule
 * @param dayName e.g. "Primday", empty for Turning/Yule
 * @param week 1-36 for regular days, 0 for Turning/Yule
 * @param seasonIndex 0-3 for regular days, -1 for Turning/Yule
 * @param isLeapYear true if this metric year has 3 Yule days
 * @param isTurning true during The Turning (3 days at spring equinox)
 * @param isYule true during Yule
 * @param isMidsummer true on Quadril 1 (summer solstice)
 * @param isSextant true on Duil 30
 * @param isTrine true on Quadril 30
 * @param isSpiral true on Quintil 18 (golden angle day)
 * @param isConvergence true on Quintil 24
 * @param isMeridian true on Sextil 30
 * @param isMask true on Octil 13
 * @param isHarmony true on Octil 30
 * @param isRest true on days 8-10 of any 10-day week
 * @param specialDay "Vigil", "Balance", "Dawn", "Yule Eve", "Midwinter", "Kindling", or empty
 * @param observance name of the observance, or empty string if none
 */
data class MetricDate(
    val year: Int,
    val month: Int,
    val monthName: String,
    val day: Int,
    val weekDay: Int,
    val dayName: String,
    val week: Int,
    val seasonIndex: Int,
    val isLeapYear: Boolean,
    val isTurning: Boolean,
    val isYule: Boolean,
    val isMidsummer: Boolean,
    val isSextant: Boolean,
    val isTrine: Boolean,
    val isSpiral: Boolean,
    val isConvergence: Boolean,
    val isMeridian: Boolean,
    val isMask: Boolean,
    val isHarmony: Boolean,
    val isRest: Boolean,
    val specialDay: String,
    val observance: String
)

object MetricCalendar {

    private fun isLeapYear(y: Int): Boolean =
        y % 4 == 0 && (y % 100 != 0 || y % 400 == 0)

    /**
     * Converts a Gregorian date to a Metric Calendar date.
     *
     * @param year Gregorian year
     * @param month Gregorian month (1-12)
     * @param day Gregorian day (1-31)
     */
    fun gregorianToMetric(year: Int, month: Int, day: Int): MetricDate {
        val date = LocalDate.of(year, month, day)
        val equinox = LocalDate.of(year, 3, 20)
        val daysFromEquinox = ChronoUnit.DAYS.between(equinox, date).toInt()

        val metricYear: Int
        val dayOfYear: Int

        if (daysFromEquinox >= 0) {
            metricYear = year - 1970
            dayOfYear = daysFromEquinox + 1
        } else {
            metricYear = year - 1 - 1970
            val prevEquinox = LocalDate.of(year - 1, 3, 20)
            dayOfYear = ChronoUnit.DAYS.between(prevEquinox, date).toInt() + 1
        }

        val leap = isLeapYear(metricYear + 1971)
        val yuleDayCount = if (leap) 3 else 2

        val base = MetricDate(
            year = metricYear, month = 0, monthName = "", day = 0, weekDay = 0,
            dayName = "", week = 0, seasonIndex = -1,
            isLeapYear = leap, isTurning = false, isYule = false,
            isMidsummer = false, isSextant = false, isTrine = false,
            isSpiral = false, isConvergence = false, isMeridian = false,
            isMask = false, isHarmony = false,
            isRest = false, specialDay = "", observance = ""
        )

        // The Turning (days 1-3)
        if (dayOfYear <= 3) {
            val name = TURNING_DAY_NAMES[dayOfYear - 1]
            return base.copy(isTurning = true, specialDay = name, observance = name)
        }

        val adjusted = dayOfYear - 3
        val m: Int
        val d: Int

        when {
            adjusted <= 270 -> {
                m = (adjusted - 1) / 30 + 1
                d = (adjusted - 1) % 30 + 1
            }
            adjusted <= 270 + yuleDayCount -> {
                val name = YULE_DAY_NAMES[adjusted - 271]
                return base.copy(isYule = true, specialDay = name, observance = name)
            }
            else -> {
                val postYule = adjusted - 270 - yuleDayCount
                m = 9 + (postYule - 1) / 30 + 1
                d = (postYule - 1) % 30 + 1
            }
        }

        val weekDay = (d - 1) % 10 + 1
        val week = (m - 1) * 3 + (d - 1) / 10 + 1

        val isMidsummer   = m == 4 && d == 1
        val isSextant     = m == 2 && d == 30
        val isTrine       = m == 4 && d == 30
        val isSpiral      = m == 5 && d == 18
        val isConvergence = m == 5 && d == 24
        val isMeridian    = m == 6 && d == 30
        val isMask        = m == 8 && d == 13
        val isHarmony     = m == 8 && d == 30

        val observance = when {
            isMidsummer   -> "Midsummer"
            isSextant     -> "The Sextant"
            isTrine       -> "The Trine"
            isSpiral      -> "The Spiral"
            isConvergence -> "Convergence"
            isMeridian    -> "The Meridian"
            isMask        -> "The Mask"
            isHarmony     -> "Harmony"
            else          -> ""
        }

        return MetricDate(
            year = metricYear,
            month = m, monthName = MONTH_NAMES[m - 1],
            day = d, weekDay = weekDay, dayName = DAY_NAMES[weekDay - 1],
            week = week, seasonIndex = (m - 1) / 3,
            isLeapYear = leap,
            isTurning = false, isYule = false,
            isMidsummer = isMidsummer,
            isSextant = isSextant,
            isTrine = isTrine,
            isSpiral = isSpiral,
            isConvergence = isConvergence,
            isMeridian = isMeridian,
            isMask = isMask,
            isHarmony = isHarmony,
            isRest = weekDay >= 8,
            specialDay = "",
            observance = observance
        )
    }

    /**
     * Formats a MetricDate using a pattern string.
     *
     * Tokens:
     * - `MMM`  month name (e.g. "Unil")
     * - `MM`   month zero-padded (e.g. "01")
     * - `M`    month number (e.g. "1")
     * - `DD`   day zero-padded (e.g. "04")
     * - `D`    day number (e.g. "4")
     * - `WW`   weekday name (e.g. "Quintday")
     * - `W`    weekday number (e.g. "5")
     * - `Y`    year number (e.g. "56")
     * - `S`    season name (e.g. "Rising")
     *
     * Example: `format(d, "WW, MMM D, Year Y")` → `"Quintday, Unil 4, Year 56"`
     */
    fun format(date: MetricDate, pattern: String): String {
        val seasonName = if (date.seasonIndex in 0..3) SEASON_NAMES[date.seasonIndex] else ""
        return FORMAT_RE.replace(pattern) { match ->
            when (match.value) {
                "MMM" -> date.monthName
                "MM"  -> "%02d".format(date.month)
                "M"   -> date.month.toString()
                "DD"  -> "%02d".format(date.day)
                "D"   -> date.day.toString()
                "WW"  -> date.dayName
                "W"   -> date.weekDay.toString()
                "Y"   -> date.year.toString()
                "S"   -> seasonName
                else  -> match.value
            }
        }
    }

    /**
     * Converts a Metric Calendar date back to a Gregorian date.
     *
     * @param year Metric year
     * @param periodType "turning", "month", or "yule"
     * @param periodValue 0-indexed: turning (0-2), month (1-12), yule (0-2)
     * @param dayOfMonth 1-30, used only when periodType is "month"
     * @return [LocalDate]
     * @throws IllegalArgumentException if inputs are invalid
     */
    fun metricToGregorian(year: Int, periodType: String, periodValue: Int, dayOfMonth: Int = 1): LocalDate {
        val equinoxYear = year + 1970
        val leap = isLeapYear(year + 1971)
        val yuleDayCount = if (leap) 3 else 2

        val offset = when (periodType) {
            "turning" -> {
                require(periodValue in 0..2) { "turning period value must be 0-2" }
                periodValue
            }
            "month" -> {
                require(periodValue in 1..12) { "month must be 1-12" }
                require(dayOfMonth in 1..30) { "day must be 1-30" }
                val m = periodValue
                val d = dayOfMonth
                if (m <= 9) {
                    3 + (m - 1) * 30 + (d - 1)
                } else {
                    3 + 270 + yuleDayCount + (m - 10) * 30 + (d - 1)
                }
            }
            "yule" -> {
                require(periodValue in 0..2) { "yule period value must be 0-2" }
                require(!(periodValue == 2 && !leap)) { "Kindling only occurs in leap years" }
                3 + 270 + periodValue
            }
            else -> throw IllegalArgumentException("Unknown periodType: $periodType")
        }

        return LocalDate.of(equinoxYear, 3, 20).plusDays(offset.toLong())
    }

    /**
     * Returns true if the given Gregorian date is a rest day (days 8-10 of any 10-day week).
     */
    fun isRestDay(year: Int, month: Int, day: Int): Boolean =
        gregorianToMetric(year, month, day).isRest

    /**
     * Returns the current Metric Calendar date.
     */
    fun today(): MetricDate {
        val now = LocalDate.now()
        return gregorianToMetric(now.year, now.monthValue, now.dayOfMonth)
    }
}
