package com.metricweek

import org.junit.Test
import org.junit.Assert.*
import java.time.LocalDate

class MetricCalendarTest {

    @Test fun testVigil() {
        val r = MetricCalendar.gregorianToMetric(2026, 3, 20)
        assertEquals(56, r.year)
        assertTrue(r.isTurning)
        assertEquals("Vigil", r.specialDay)
    }

    @Test fun testBalance() {
        val r = MetricCalendar.gregorianToMetric(2026, 3, 21)
        assertTrue(r.isTurning)
        assertEquals("Balance", r.specialDay)
    }

    @Test fun testDawn() {
        val r = MetricCalendar.gregorianToMetric(2026, 3, 22)
        assertTrue(r.isTurning)
        assertEquals("Dawn", r.specialDay)
    }

    @Test fun testUnil1() {
        val r = MetricCalendar.gregorianToMetric(2026, 3, 23)
        assertEquals(56, r.year)
        assertEquals(1, r.month)
        assertEquals("Unil", r.monthName)
        assertEquals(1, r.day)
        assertEquals(1, r.weekDay)
        assertEquals("Primday", r.dayName)
        assertFalse(r.isRest)
        assertEquals(1, r.week)
    }

    @Test fun testUnil10() {
        val r = MetricCalendar.gregorianToMetric(2026, 4, 1)
        assertEquals(1, r.month)
        assertEquals(10, r.day)
        assertEquals(10, r.weekDay)
        assertEquals("Decday", r.dayName)
        assertTrue(r.isRest)
    }

    @Test fun testYuleEve() {
        val r = MetricCalendar.gregorianToMetric(2026, 12, 18)
        assertEquals(56, r.year)
        assertTrue(r.isYule)
        assertEquals("Yule Eve", r.specialDay)
    }

    @Test fun testPreEquinox() {
        val r = MetricCalendar.gregorianToMetric(2025, 1, 1)
        assertEquals(54, r.year)  // before March 2025 equinox → year 54
        assertEquals(10, r.month)
        assertEquals("Decil", r.monthName)
        assertEquals(13, r.day)
        assertEquals(3, r.weekDay)
        assertEquals("Triday", r.dayName)
        assertFalse(r.isRest)
    }

    @Test fun testMidsummer() {
        val r = MetricCalendar.gregorianToMetric(2026, 6, 21)
        assertTrue(r.isMidsummer)
        assertEquals(4, r.month)
        assertEquals(1, r.day)
    }

    @Test fun testIsRestDay() {
        assertTrue(MetricCalendar.isRestDay(2026, 4, 1))
        assertFalse(MetricCalendar.isRestDay(2026, 3, 23))
    }

    @Test fun testMetricToGregorianMonth() {
        val d = MetricCalendar.metricToGregorian(56, "month", 1, 1)
        assertEquals(LocalDate.of(2026, 3, 23), d)
    }

    @Test fun testMetricToGregorianTurning() {
        val d = MetricCalendar.metricToGregorian(56, "turning", 0)
        assertEquals(LocalDate.of(2026, 3, 20), d)
    }

    @Test fun testMetricToGregorianYule() {
        val d = MetricCalendar.metricToGregorian(56, "yule", 0)
        assertEquals(LocalDate.of(2026, 12, 18), d)
    }

    @Test(expected = IllegalArgumentException::class)
    fun testKindlingRequiresLeap() {
        MetricCalendar.metricToGregorian(56, "yule", 2)
    }
}
