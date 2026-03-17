package com.metricweek;

/**
 * A date in the Metric Calendar system.
 *
 * <p>Year 0 = spring equinox 1970. Each metric year begins at the spring equinox (March 20).
 */
public final class MetricDate {
    /** Metric year (Year 0 = spring equinox 1970) */
    public final int year;
    /** 1-12 for regular days, 0 for Turning/Yule */
    public final int month;
    /** e.g. "Unil", empty string for Turning/Yule */
    public final String monthName;
    /** 1-30 for regular days, 0 for Turning/Yule */
    public final int day;
    /** 1-10 for regular days, 0 for Turning/Yule */
    public final int weekDay;
    /** e.g. "Primday", empty string for Turning/Yule */
    public final String dayName;
    /** 1-36 for regular days, 0 for Turning/Yule */
    public final int week;
    /** 0-3 for regular days, -1 for Turning/Yule */
    public final int seasonIndex;
    /** true if this metric year has 3 Yule days */
    public final boolean isLeapYear;
    /** true during The Turning (3 days at spring equinox) */
    public final boolean isTurning;
    /** true during Yule */
    public final boolean isYule;
    /** true on Quadril 1 (summer solstice) */
    public final boolean isMidsummer;
    /** true on Duil 30 */
    public final boolean isSextant;
    /** true on Quadril 30 */
    public final boolean isTrine;
    /** true on Quintil 18 (golden angle day) */
    public final boolean isSpiral;
    /** true on Quintil 24 */
    public final boolean isConvergence;
    /** true on Sextil 30 */
    public final boolean isMeridian;
    /** true on Octil 13 */
    public final boolean isMask;
    /** true on Octil 30 */
    public final boolean isHarmony;
    /** true on days 8-10 of any 10-day week */
    public final boolean isRest;
    /** "Vigil", "Balance", "Dawn", "Yule Eve", "Midwinter", "Kindling", or empty */
    public final String specialDay;
    /** Name of the observance, or empty string if none */
    public final String observance;

    public MetricDate(int year, int month, String monthName, int day, int weekDay, String dayName,
                      int week, int seasonIndex, boolean isLeapYear, boolean isTurning,
                      boolean isYule, boolean isMidsummer, boolean isSextant, boolean isTrine,
                      boolean isSpiral, boolean isConvergence, boolean isMeridian,
                      boolean isMask, boolean isHarmony, boolean isRest,
                      String specialDay, String observance) {
        this.year = year;
        this.month = month;
        this.monthName = monthName;
        this.day = day;
        this.weekDay = weekDay;
        this.dayName = dayName;
        this.week = week;
        this.seasonIndex = seasonIndex;
        this.isLeapYear = isLeapYear;
        this.isTurning = isTurning;
        this.isYule = isYule;
        this.isMidsummer = isMidsummer;
        this.isSextant = isSextant;
        this.isTrine = isTrine;
        this.isSpiral = isSpiral;
        this.isConvergence = isConvergence;
        this.isMeridian = isMeridian;
        this.isMask = isMask;
        this.isHarmony = isHarmony;
        this.isRest = isRest;
        this.specialDay = specialDay;
        this.observance = observance;
    }

    @Override
    public String toString() {
        if (isTurning) return "Year " + year + ", The Turning — " + specialDay;
        if (isYule) return "Year " + year + ", Yule — " + specialDay;
        return "Year " + year + ", " + monthName + " " + day + " (" + dayName + ")";
    }
}
