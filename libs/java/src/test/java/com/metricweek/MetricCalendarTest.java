package com.metricweek;

import org.junit.Test;
import static org.junit.Assert.*;

public class MetricCalendarTest {

    @Test public void testVigil() {
        MetricDate r = MetricCalendar.fromGregorian(2026, 3, 20);
        assertEquals(56, r.year);
        assertTrue(r.isTurning);
        assertEquals("Vigil", r.specialDay);
    }

    @Test public void testBalance() {
        MetricDate r = MetricCalendar.fromGregorian(2026, 3, 21);
        assertTrue(r.isTurning);
        assertEquals("Balance", r.specialDay);
    }

    @Test public void testDawn() {
        MetricDate r = MetricCalendar.fromGregorian(2026, 3, 22);
        assertTrue(r.isTurning);
        assertEquals("Dawn", r.specialDay);
    }

    @Test public void testUnil1() {
        MetricDate r = MetricCalendar.fromGregorian(2026, 3, 23);
        assertEquals(56, r.year);
        assertEquals(1, r.month);
        assertEquals("Unil", r.monthName);
        assertEquals(1, r.day);
        assertEquals(1, r.weekDay);
        assertEquals("Primday", r.dayName);
        assertFalse(r.isRest);
        assertEquals(1, r.week);
    }

    @Test public void testUnil10() {
        MetricDate r = MetricCalendar.fromGregorian(2026, 4, 1);
        assertEquals(1, r.month);
        assertEquals(10, r.day);
        assertEquals(10, r.weekDay);
        assertEquals("Decday", r.dayName);
        assertTrue(r.isRest);
    }

    @Test public void testYuleEve() {
        MetricDate r = MetricCalendar.fromGregorian(2026, 12, 18);
        assertEquals(56, r.year);
        assertTrue(r.isYule);
        assertEquals("Yule Eve", r.specialDay);
    }

    @Test public void testPreEquinox() {
        MetricDate r = MetricCalendar.fromGregorian(2025, 1, 1);
        assertEquals(54, r.year);  // before March 2025 equinox → year 54
        assertEquals(10, r.month);
        assertEquals("Decil", r.monthName);
        assertEquals(13, r.day);
        assertEquals(3, r.weekDay);
        assertEquals("Triday", r.dayName);
        assertFalse(r.isRest);
    }

    @Test public void testMidsummer() {
        MetricDate r = MetricCalendar.fromGregorian(2026, 6, 21);
        assertTrue(r.isMidsummer);
        assertEquals(4, r.month);
        assertEquals(1, r.day);
    }

    @Test public void testIsRestDay() {
        assertTrue(MetricCalendar.isRestDay(2026, 4, 1));
        assertFalse(MetricCalendar.isRestDay(2026, 3, 23));
    }

    @Test public void testToGregorianMonth() {
        int[] d = MetricCalendar.toGregorian(56, "month", 1, 1);
        assertArrayEquals(new int[]{2026, 3, 23}, d);
    }

    @Test public void testToGregorianTurning() {
        int[] d = MetricCalendar.toGregorian(56, "turning", 0, 0);
        assertArrayEquals(new int[]{2026, 3, 20}, d);
    }

    @Test public void testToGregorianYule() {
        int[] d = MetricCalendar.toGregorian(56, "yule", 0, 0);
        assertArrayEquals(new int[]{2026, 12, 18}, d);
    }

    @Test(expected = IllegalArgumentException.class)
    public void testKindlingRequiresLeap() {
        MetricCalendar.toGregorian(56, "yule", 2, 0);
    }
}
