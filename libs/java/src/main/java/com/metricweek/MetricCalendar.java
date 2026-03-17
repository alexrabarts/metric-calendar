package com.metricweek;

import java.util.Calendar;
import java.util.TimeZone;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

/**
 * Metric Calendar date conversion utility.
 *
 * <p>The Metric Calendar is a rational decimal calendar system with 10-day weeks
 * and 12 months of 30 days. Year 0 = spring equinox 1970.
 *
 * @see <a href="https://metricweek.com">metricweek.com</a>
 */
public class MetricCalendar {

    private static final String[] DAY_NAMES = {
        "Primday", "Duoday", "Triday", "Quadday", "Quintday",
        "Hexday", "Septday", "Octday", "Novday", "Decday"
    };

    private static final String[] MONTH_NAMES = {
        "Unil", "Duil", "Tril", "Quadril", "Quintil", "Sextil",
        "Septil", "Octil", "Novil", "Decil", "Undecil", "Duodecil"
    };

    private static final String[] SEASON_NAMES = {
        "Rising", "Flourishing", "Gathering", "Stillness"
    };

    private static final String[] TURNING_DAY_NAMES = {"Vigil", "Balance", "Dawn"};
    private static final String[] YULE_DAY_NAMES = {"Yule Eve", "Midwinter", "Kindling"};

    private static final Pattern FORMAT_RE = Pattern.compile("MMM|MM|M|DD|D|WW|W|Y|S");

    private MetricCalendar() {}

    private static boolean isLeapYear(int y) {
        return y % 4 == 0 && (y % 100 != 0 || y % 400 == 0);
    }

    /** Compute Julian Day Number for a Gregorian date. */
    private static int julianDayNumber(int year, int month, int day) {
        int a = (14 - month) / 12;
        int y = year + 4800 - a;
        int m = month + 12 * a - 3;
        return day + (153 * m + 2) / 5 + 365 * y + y / 4 - y / 100 + y / 400 - 32045;
    }

    /** Convert Julian Day Number to Gregorian date as int[]{year, month, day}. */
    private static int[] gregorianFromJDN(int jdn) {
        int a = jdn + 32044;
        int b = (4 * a + 3) / 146097;
        int c = a - (146097 * b) / 4;
        int d = (4 * c + 3) / 1461;
        int e = c - (1461 * d) / 4;
        int m = (5 * e + 2) / 153;
        int day = e - (153 * m + 2) / 5 + 1;
        int month = m + 3 - 12 * (m / 10);
        int year = 100 * b + d - 4800 + m / 10;
        return new int[]{year, month, day};
    }

    /**
     * Converts a Gregorian date to a {@link MetricDate}.
     *
     * @param year  Gregorian year
     * @param month Gregorian month (1-12)
     * @param day   Gregorian day (1-31)
     * @return the corresponding Metric Calendar date
     */
    public static MetricDate fromGregorian(int year, int month, int day) {
        int dateJDN = julianDayNumber(year, month, day);
        int equinoxJDN = julianDayNumber(year, 3, 20);
        int daysFromEquinox = dateJDN - equinoxJDN;

        int metricYear, dayOfYear;
        if (daysFromEquinox >= 0) {
            metricYear = year - 1970;
            dayOfYear = daysFromEquinox + 1;
        } else {
            metricYear = year - 1 - 1970;
            int prevEquinoxJDN = julianDayNumber(year - 1, 3, 20);
            dayOfYear = dateJDN - prevEquinoxJDN + 1;
        }

        boolean leap = isLeapYear(metricYear + 1971);
        int yuleDayCount = leap ? 3 : 2;

        // The Turning (days 1-3)
        if (dayOfYear <= 3) {
            String name = TURNING_DAY_NAMES[dayOfYear - 1];
            return new MetricDate(
                metricYear, 0, "", 0, 0, "", 0, -1,
                leap, true, false,
                false, false, false, false, false, false, false, false,
                false, name, name
            );
        }

        int adjusted = dayOfYear - 3;
        int m, d;

        if (adjusted <= 270) {
            m = (adjusted - 1) / 30 + 1;
            d = (adjusted - 1) % 30 + 1;
        } else if (adjusted <= 270 + yuleDayCount) {
            String name = YULE_DAY_NAMES[adjusted - 271];
            return new MetricDate(
                metricYear, 0, "", 0, 0, "", 0, -1,
                leap, false, true,
                false, false, false, false, false, false, false, false,
                false, name, name
            );
        } else {
            int postYule = adjusted - 270 - yuleDayCount;
            m = 9 + (postYule - 1) / 30 + 1;
            d = (postYule - 1) % 30 + 1;
        }

        int weekDay = (d - 1) % 10 + 1;
        int week = (m - 1) * 3 + (d - 1) / 10 + 1;

        boolean isMidsummer   = m == 4 && d == 1;
        boolean isSextant     = m == 2 && d == 30;
        boolean isTrine       = m == 4 && d == 30;
        boolean isSpiral      = m == 5 && d == 18;
        boolean isConvergence = m == 5 && d == 24;
        boolean isMeridian    = m == 6 && d == 30;
        boolean isMask        = m == 8 && d == 13;
        boolean isHarmony     = m == 8 && d == 30;

        String observance;
        if      (isMidsummer)   observance = "Midsummer";
        else if (isSextant)     observance = "The Sextant";
        else if (isTrine)       observance = "The Trine";
        else if (isSpiral)      observance = "The Spiral";
        else if (isConvergence) observance = "Convergence";
        else if (isMeridian)    observance = "The Meridian";
        else if (isMask)        observance = "The Mask";
        else if (isHarmony)     observance = "Harmony";
        else                    observance = "";

        return new MetricDate(
            metricYear,
            m, MONTH_NAMES[m - 1],
            d, weekDay, DAY_NAMES[weekDay - 1],
            week, (m - 1) / 3,
            leap,
            false, false,
            isMidsummer, isSextant, isTrine, isSpiral,
            isConvergence, isMeridian, isMask, isHarmony,
            weekDay >= 8,
            "", observance
        );
    }

    /**
     * Formats a MetricDate using a pattern string.
     *
     * <p>Tokens:
     * <ul>
     *   <li>{@code MMM} — month name (e.g. "Unil")</li>
     *   <li>{@code MM}  — month zero-padded (e.g. "01")</li>
     *   <li>{@code M}   — month number (e.g. "1")</li>
     *   <li>{@code DD}  — day zero-padded (e.g. "04")</li>
     *   <li>{@code D}   — day number (e.g. "4")</li>
     *   <li>{@code WW}  — weekday name (e.g. "Quintday")</li>
     *   <li>{@code W}   — weekday number (e.g. "5")</li>
     *   <li>{@code Y}   — year number (e.g. "56")</li>
     *   <li>{@code S}   — season name (e.g. "Rising")</li>
     * </ul>
     *
     * <p>Example: {@code format(d, "WW, MMM D, Year Y")} → {@code "Quintday, Unil 4, Year 56"}
     *
     * @param date    the MetricDate to format
     * @param pattern the pattern string
     * @return the formatted string
     */
    public static String format(MetricDate date, String pattern) {
        String seasonName = (date.seasonIndex >= 0 && date.seasonIndex <= 3)
            ? SEASON_NAMES[date.seasonIndex] : "";
        Matcher m = FORMAT_RE.matcher(pattern);
        StringBuilder sb = new StringBuilder();
        while (m.find()) {
            String replacement;
            switch (m.group()) {
                case "MMM": replacement = date.monthName; break;
                case "MM":  replacement = String.format("%02d", date.month); break;
                case "M":   replacement = Integer.toString(date.month); break;
                case "DD":  replacement = String.format("%02d", date.day); break;
                case "D":   replacement = Integer.toString(date.day); break;
                case "WW":  replacement = date.dayName; break;
                case "W":   replacement = Integer.toString(date.weekDay); break;
                case "Y":   replacement = Integer.toString(date.year); break;
                case "S":   replacement = seasonName; break;
                default:    replacement = m.group(); break;
            }
            m.appendReplacement(sb, Matcher.quoteReplacement(replacement));
        }
        m.appendTail(sb);
        return sb.toString();
    }

    /**
     * Converts a Metric Calendar date back to a Gregorian date.
     *
     * @param year       Metric year
     * @param periodType "turning", "month", or "yule"
     * @param periodValue 0-indexed for turning/yule (0-2), or 1-12 for month
     * @param dayOfMonth 1-30, used only when periodType is "month"
     * @return int array {gregorianYear, gregorianMonth, gregorianDay}
     * @throws IllegalArgumentException if inputs are invalid
     */
    public static int[] toGregorian(int year, String periodType, int periodValue, int dayOfMonth) {
        int equinoxYear = year + 1970;
        boolean leap = isLeapYear(year + 1971);
        int yuleDayCount = leap ? 3 : 2;
        int offset;

        switch (periodType) {
            case "turning":
                if (periodValue < 0 || periodValue > 2)
                    throw new IllegalArgumentException("turning periodValue must be 0-2");
                offset = periodValue;
                break;
            case "month":
                if (periodValue < 1 || periodValue > 12)
                    throw new IllegalArgumentException("month must be 1-12");
                if (dayOfMonth < 1 || dayOfMonth > 30)
                    throw new IllegalArgumentException("day must be 1-30");
                int mv = periodValue, dv = dayOfMonth;
                if (mv <= 9) {
                    offset = 3 + (mv - 1) * 30 + (dv - 1);
                } else {
                    offset = 3 + 270 + yuleDayCount + (mv - 10) * 30 + (dv - 1);
                }
                break;
            case "yule":
                if (periodValue < 0 || periodValue > 2)
                    throw new IllegalArgumentException("yule periodValue must be 0-2");
                if (periodValue == 2 && !leap)
                    throw new IllegalArgumentException("Kindling only occurs in leap years");
                offset = 3 + 270 + periodValue;
                break;
            default:
                throw new IllegalArgumentException("Unknown periodType: " + periodType);
        }

        int equinoxJDN = julianDayNumber(equinoxYear, 3, 20);
        return gregorianFromJDN(equinoxJDN + offset);
    }

    /**
     * Returns true if the given Gregorian date is a rest day.
     */
    public static boolean isRestDay(int year, int month, int day) {
        return fromGregorian(year, month, day).isRest;
    }

    /**
     * Returns the current Metric Calendar date using the system clock (UTC).
     */
    public static MetricDate today() {
        Calendar cal = Calendar.getInstance(TimeZone.getTimeZone("UTC"));
        return fromGregorian(
            cal.get(Calendar.YEAR),
            cal.get(Calendar.MONTH) + 1,
            cal.get(Calendar.DAY_OF_MONTH)
        );
    }
}
