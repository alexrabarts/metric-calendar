import XCTest
@testable import MetricCalendar

final class MetricCalendarTests: XCTestCase {

    // Helper: create a UTC Date from year/month/day
    private func utcDate(_ year: Int, _ month: Int, _ day: Int) -> Date {
        var comps = DateComponents()
        comps.year = year; comps.month = month; comps.day = day
        comps.hour = 0; comps.minute = 0; comps.second = 0
        var cal = Calendar(identifier: .gregorian)
        cal.timeZone = TimeZone(identifier: "UTC")!
        return cal.date(from: comps)!
    }

    func testVigil() {
        let r = gregorianToMetric(utcDate(2026, 3, 20))
        XCTAssertEqual(r.year, 56)
        XCTAssertTrue(r.isTurning)
        XCTAssertEqual(r.specialDay, "Vigil")
    }

    func testBalance() {
        let r = gregorianToMetric(utcDate(2026, 3, 21))
        XCTAssertTrue(r.isTurning)
        XCTAssertEqual(r.specialDay, "Balance")
    }

    func testDawn() {
        let r = gregorianToMetric(utcDate(2026, 3, 22))
        XCTAssertTrue(r.isTurning)
        XCTAssertEqual(r.specialDay, "Dawn")
    }

    func testUnil1() {
        let r = gregorianToMetric(utcDate(2026, 3, 23))
        XCTAssertEqual(r.year, 56)
        XCTAssertEqual(r.month, 1)
        XCTAssertEqual(r.monthName, "Unil")
        XCTAssertEqual(r.day, 1)
        XCTAssertEqual(r.weekDay, 1)
        XCTAssertEqual(r.dayName, "Primday")
        XCTAssertFalse(r.isRest)
        XCTAssertEqual(r.week, 1)
    }

    func testUnil10() {
        let r = gregorianToMetric(utcDate(2026, 4, 1))
        XCTAssertEqual(r.month, 1)
        XCTAssertEqual(r.day, 10)
        XCTAssertEqual(r.weekDay, 10)
        XCTAssertEqual(r.dayName, "Decday")
        XCTAssertTrue(r.isRest)
    }

    func testYuleEve() {
        let r = gregorianToMetric(utcDate(2026, 12, 18))
        XCTAssertEqual(r.year, 56)
        XCTAssertTrue(r.isYule)
        XCTAssertEqual(r.specialDay, "Yule Eve")
    }

    func testPreEquinox() {
        let r = gregorianToMetric(utcDate(2025, 1, 1))
        XCTAssertEqual(r.year, 54)  // before March 2025 equinox → year 54
        XCTAssertEqual(r.month, 10)
        XCTAssertEqual(r.monthName, "Decil")
        XCTAssertEqual(r.day, 13)
        XCTAssertEqual(r.weekDay, 3)
        XCTAssertEqual(r.dayName, "Triday")
        XCTAssertFalse(r.isRest)
    }

    func testMidsummer() {
        let r = gregorianToMetric(utcDate(2026, 6, 21))
        XCTAssertTrue(r.isMidsummer)
        XCTAssertEqual(r.month, 4)
        XCTAssertEqual(r.day, 1)
    }

    func testIsRestDay() {
        XCTAssertTrue(isRestDay(utcDate(2026, 4, 1)))
        XCTAssertFalse(isRestDay(utcDate(2026, 3, 23)))
    }

    func testMetricToGregorianMonth() throws {
        let d = try metricToGregorian(year: 56, periodType: "month", periodValue: 1, dayOfMonth: 1)
        var cal = Calendar(identifier: .gregorian)
        cal.timeZone = TimeZone(identifier: "UTC")!
        XCTAssertEqual(cal.component(.year, from: d), 2026)
        XCTAssertEqual(cal.component(.month, from: d), 3)
        XCTAssertEqual(cal.component(.day, from: d), 23)
    }

    func testMetricToGregorianTurning() throws {
        let d = try metricToGregorian(year: 56, periodType: "turning", periodValue: 0)
        var cal = Calendar(identifier: .gregorian)
        cal.timeZone = TimeZone(identifier: "UTC")!
        XCTAssertEqual(cal.component(.year, from: d), 2026)
        XCTAssertEqual(cal.component(.month, from: d), 3)
        XCTAssertEqual(cal.component(.day, from: d), 20)
    }

    func testMetricToGregorianYule() throws {
        let d = try metricToGregorian(year: 56, periodType: "yule", periodValue: 0)
        var cal = Calendar(identifier: .gregorian)
        cal.timeZone = TimeZone(identifier: "UTC")!
        XCTAssertEqual(cal.component(.year, from: d), 2026)
        XCTAssertEqual(cal.component(.month, from: d), 12)
        XCTAssertEqual(cal.component(.day, from: d), 18)
    }

    func testKindlingRequiresLeap() {
        XCTAssertThrowsError(try metricToGregorian(year: 56, periodType: "yule", periodValue: 2))
    }

    // Additional edge-case tests

    func testYearBoundary_JustBeforeEquinox() {
        // March 19, 2026 is the last day of metric year 55
        let r = gregorianToMetric(utcDate(2026, 3, 19))
        XCTAssertEqual(r.year, 55)
    }

    func testYearBoundary_EquinoxIsNewYear() {
        // March 20, 2026 is the first day (Vigil) of metric year 56
        let r = gregorianToMetric(utcDate(2026, 3, 20))
        XCTAssertEqual(r.year, 56)
        XCTAssertTrue(r.isTurning)
        XCTAssertEqual(r.specialDay, "Vigil")
    }

    func testSeasonIndex() {
        // Unil (month 1) → season 0
        let spring = gregorianToMetric(utcDate(2026, 3, 23))
        XCTAssertEqual(spring.seasonIndex, 0)

        // Quadril (month 4) → season 1
        let summer = gregorianToMetric(utcDate(2026, 6, 21))
        XCTAssertEqual(summer.seasonIndex, 1)

        // Septil (month 7) → season 2
        let autumn = gregorianToMetric(utcDate(2026, 9, 19))
        XCTAssertEqual(autumn.seasonIndex, 2)

        // Decil (month 10) → season 3
        let winter = gregorianToMetric(utcDate(2026, 12, 21))
        XCTAssertEqual(winter.seasonIndex, 3)
    }

    func testTurningIsNotRest() {
        let vigil = gregorianToMetric(utcDate(2026, 3, 20))
        XCTAssertFalse(vigil.isRest)
        XCTAssertEqual(vigil.weekDay, 0)
    }

    func testYuleIsNotRest() {
        let yuleEve = gregorianToMetric(utcDate(2026, 12, 18))
        XCTAssertFalse(yuleEve.isRest)
        XCTAssertEqual(yuleEve.weekDay, 0)
    }

    func testDescriptionFormatTurning() {
        let r = gregorianToMetric(utcDate(2026, 3, 20))
        XCTAssertEqual(r.description, "Year 56, The Turning — Vigil")
    }

    func testDescriptionFormatYule() {
        let r = gregorianToMetric(utcDate(2026, 12, 18))
        XCTAssertEqual(r.description, "Year 56, Yule — Yule Eve")
    }

    func testDescriptionFormatRegular() {
        let r = gregorianToMetric(utcDate(2026, 3, 23))
        XCTAssertEqual(r.description, "Year 56, Unil 1 (Primday)")
    }

    func testInvalidTurningValue() {
        XCTAssertThrowsError(try metricToGregorian(year: 56, periodType: "turning", periodValue: 3))
    }

    func testInvalidMonth() {
        XCTAssertThrowsError(try metricToGregorian(year: 56, periodType: "month", periodValue: 13, dayOfMonth: 1))
        XCTAssertThrowsError(try metricToGregorian(year: 56, periodType: "month", periodValue: 0, dayOfMonth: 1))
    }

    func testInvalidDay() {
        XCTAssertThrowsError(try metricToGregorian(year: 56, periodType: "month", periodValue: 1, dayOfMonth: 31))
        XCTAssertThrowsError(try metricToGregorian(year: 56, periodType: "month", periodValue: 1, dayOfMonth: 0))
    }

    func testUnknownPeriodType() {
        XCTAssertThrowsError(try metricToGregorian(year: 56, periodType: "quarter", periodValue: 1))
    }

    func testRoundTrip_Month() throws {
        // Convert metric → gregorian → metric and check we get the same date
        let original = gregorianToMetric(utcDate(2026, 7, 15))
        let gregorian = try metricToGregorian(
            year: original.year,
            periodType: "month",
            periodValue: original.month,
            dayOfMonth: original.day
        )
        let roundTripped = gregorianToMetric(gregorian)
        XCTAssertEqual(roundTripped.year, original.year)
        XCTAssertEqual(roundTripped.month, original.month)
        XCTAssertEqual(roundTripped.day, original.day)
    }

    func testRoundTrip_Turning() throws {
        let gregorian = try metricToGregorian(year: 56, periodType: "turning", periodValue: 1)
        let r = gregorianToMetric(gregorian)
        XCTAssertTrue(r.isTurning)
        XCTAssertEqual(r.specialDay, "Balance")
        XCTAssertEqual(r.year, 56)
    }

    func testWeekNumbers() {
        // Week 1: Unil 1-10, Week 2: Unil 11-20, Week 3: Unil 21-30
        let w1 = gregorianToMetric(utcDate(2026, 3, 23))  // Unil 1
        XCTAssertEqual(w1.week, 1)

        let w2 = gregorianToMetric(utcDate(2026, 4, 2))   // Unil 11
        XCTAssertEqual(w2.week, 2)

        let w3 = gregorianToMetric(utcDate(2026, 4, 12))  // Unil 21
        XCTAssertEqual(w3.week, 3)

        let w4 = gregorianToMetric(utcDate(2026, 4, 22))  // Duil 1
        XCTAssertEqual(w4.week, 4)
    }

    func testSpiralDay() {
        // Quintil 18 is the golden angle day
        // offset from equinox = 3 (Turning) + 4*30 (months 1-4) + 17 (day 18, 0-indexed) = 140
        // 2026-03-20 + 140 days = 2026-08-07
        let r = gregorianToMetric(utcDate(2026, 8, 7))  // Quintil 18
        XCTAssertTrue(r.isSpiral)
        XCTAssertEqual(r.month, 5)
        XCTAssertEqual(r.day, 18)
    }
}
