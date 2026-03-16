import Foundation

// Calendar constants
private let dayNames: [String] = [
    "Primday", "Duoday", "Triday", "Quadday", "Quintday",
    "Hexday", "Septday", "Octday", "Novday", "Decday",
]

private let monthNames: [String] = [
    "Unil", "Duil", "Tril", "Quadril", "Quintil", "Sextil",
    "Septil", "Octil", "Novil", "Decil", "Undecil", "Duodecil",
]

private let turningDayNames: [String] = ["Vigil", "Balance", "Dawn"]
private let yuleDayNames: [String] = ["Yule Eve", "Midwinter", "Kindling"]

/// A date in the Metric Calendar system.
public struct MetricDate: Equatable, Hashable, CustomStringConvertible {
    /// Metric year (Year 0 = spring equinox 1970)
    public let year: Int
    /// 1-12 for regular days, 0 for Turning/Yule
    public let month: Int
    public let monthName: String
    /// 1-30 for regular days, 0 for Turning/Yule
    public let day: Int
    /// 1-10 for regular days, 0 for Turning/Yule
    public let weekDay: Int
    public let dayName: String
    /// 1-36 for regular days, 0 for Turning/Yule
    public let week: Int
    /// 0-3 for regular days, -1 for Turning/Yule
    public let seasonIndex: Int
    public let isLeapYear: Bool
    /// True during The Turning (3 days at spring equinox)
    public let isTurning: Bool
    /// True during Yule
    public let isYule: Bool
    /// True on Quadril 1 (summer solstice)
    public let isMidsummer: Bool
    /// True on Quintil 18 (golden angle day)
    public let isSpiral: Bool
    /// True on days 8-10 of any 10-day week
    public let isRest: Bool
    /// "Vigil", "Balance", "Dawn", "Yule Eve", "Midwinter", "Kindling", or ""
    public let specialDay: String

    public var description: String {
        if isTurning { return "Year \(year), The Turning — \(specialDay)" }
        if isYule { return "Year \(year), Yule — \(specialDay)" }
        return "Year \(year), \(monthName) \(day) (\(dayName))"
    }
}

/// Errors thrown by metric-to-Gregorian conversion.
public enum MetricCalendarError: Error, LocalizedError {
    case invalidTurningValue
    case invalidMonth
    case invalidDay
    case invalidYuleValue
    case kindlingRequiresLeapYear
    case unknownPeriodType(String)

    public var errorDescription: String? {
        switch self {
        case .invalidTurningValue: return "Turning period value must be 0-2"
        case .invalidMonth: return "Month must be 1-12"
        case .invalidDay: return "Day must be 1-30"
        case .invalidYuleValue: return "Yule period value must be 0-2"
        case .kindlingRequiresLeapYear: return "Kindling only occurs in leap years"
        case .unknownPeriodType(let t): return "Unknown period type: \(t)"
        }
    }
}

// Helper: UTC midnight timestamp for a given Gregorian date using Julian Day Number arithmetic
private func utcTimestamp(year: Int, month: Int, day: Int) -> TimeInterval {
    let a = (14 - month) / 12
    let y = year + 4800 - a
    let m = month + 12 * a - 3
    let jdn = day + (153 * m + 2) / 5 + 365 * y + y / 4 - y / 100 + y / 400 - 32045
    let unixEpochJDN = 2440588  // JDN for 1970-01-01
    return TimeInterval((jdn - unixEpochJDN) * 86400)
}

// Helper: check Gregorian leap year
private func isLeapYear(_ y: Int) -> Bool {
    return y % 4 == 0 && (y % 100 != 0 || y % 400 == 0)
}

// Helper: convert Unix timestamp back to Gregorian date components
private func gregorianFromTimestamp(_ ts: TimeInterval) -> (year: Int, month: Int, day: Int) {
    let days = Int(ts / 86400)
    let jdn = days + 2440588  // 1970-01-01 is JDN 2440588
    let a = jdn + 32044
    let b = (4 * a + 3) / 146097
    let c = a - (146097 * b) / 4
    let d = (4 * c + 3) / 1461
    let e = c - (1461 * d) / 4
    let m = (5 * e + 2) / 153
    let day = e - (153 * m + 2) / 5 + 1
    let month = m + 3 - 12 * (m / 10)
    let year = 100 * b + d - 4800 + m / 10
    return (year, month, day)
}

/// Converts a Gregorian `Date` to a `MetricDate`.
/// The date is interpreted as UTC.
public func gregorianToMetric(_ date: Date) -> MetricDate {
    // Normalize to UTC midnight using raw Unix timestamp math
    let ts = date.timeIntervalSince1970
    let daysSinceEpoch = Int(floor(ts / 86400))
    let utcMidnight = TimeInterval(daysSinceEpoch * 86400)

    let (year, month, day) = gregorianFromTimestamp(utcMidnight)
    let _ = (month, day)  // suppress unused warnings; used implicitly via tuple

    let equinoxTs = utcTimestamp(year: year, month: 3, day: 20)
    let daysFromEquinox = Int((utcMidnight - equinoxTs) / 86400)

    let metricYear: Int
    let dayOfYear: Int

    if daysFromEquinox >= 0 {
        metricYear = year - 1970
        dayOfYear = daysFromEquinox + 1
    } else {
        metricYear = year - 1 - 1970
        let prevEquinoxTs = utcTimestamp(year: year - 1, month: 3, day: 20)
        dayOfYear = Int((utcMidnight - prevEquinoxTs) / 86400) + 1
    }

    // Leap year is determined by whether the Gregorian year that follows the equinox is a leap year.
    // The metric year starts at the equinox of `equinoxYear`. The following calendar year is
    // equinoxYear + 1, which determines whether Yule has a Kindling day.
    let leap = isLeapYear(metricYear + 1971)
    let yuleDayCount = leap ? 3 : 2

    func makeTurning(_ idx: Int) -> MetricDate {
        MetricDate(
            year: metricYear, month: 0, monthName: "", day: 0, weekDay: 0,
            dayName: "", week: 0, seasonIndex: -1, isLeapYear: leap,
            isTurning: true, isYule: false, isMidsummer: false, isSpiral: false,
            isRest: false, specialDay: turningDayNames[idx]
        )
    }

    func makeYule(_ idx: Int) -> MetricDate {
        MetricDate(
            year: metricYear, month: 0, monthName: "", day: 0, weekDay: 0,
            dayName: "", week: 0, seasonIndex: -1, isLeapYear: leap,
            isTurning: false, isYule: true, isMidsummer: false, isSpiral: false,
            isRest: false, specialDay: yuleDayNames[idx]
        )
    }

    // The Turning: days 1-3 of the metric year
    if dayOfYear <= 3 {
        return makeTurning(dayOfYear - 1)
    }

    // adjusted is 1-indexed position within the post-Turning calendar
    let adjusted = dayOfYear - 3

    let m: Int
    let d: Int

    if adjusted <= 270 {
        // Months 1-9: 270 days total
        m = (adjusted - 1) / 30 + 1
        d = (adjusted - 1) % 30 + 1
    } else if adjusted <= 270 + yuleDayCount {
        // Yule: 2 or 3 days
        return makeYule(adjusted - 271)
    } else {
        // Months 10-12: post-Yule
        let postYule = adjusted - 270 - yuleDayCount
        m = 9 + (postYule - 1) / 30 + 1
        d = (postYule - 1) % 30 + 1
    }

    let weekDay = (d - 1) % 10 + 1
    let week = (m - 1) * 3 + (d - 1) / 10 + 1

    return MetricDate(
        year: metricYear,
        month: m, monthName: monthNames[m - 1],
        day: d, weekDay: weekDay, dayName: dayNames[weekDay - 1],
        week: week, seasonIndex: (m - 1) / 3,
        isLeapYear: leap,
        isTurning: false, isYule: false,
        isMidsummer: m == 4 && d == 1,
        isSpiral: m == 5 && d == 18,
        isRest: weekDay >= 8,
        specialDay: ""
    )
}

/// Converts a Metric Calendar date back to a Gregorian `Date`.
///
/// - Parameters:
///   - year: Metric year
///   - periodType: "turning", "month", or "yule"
///   - periodValue: 0-indexed: turning (0-2), month (1-12), yule (0-2)
///   - dayOfMonth: 1-30, used only when periodType is "month"
/// - Returns: UTC midnight `Date` for the corresponding Gregorian date
/// - Throws: `MetricCalendarError` if inputs are invalid
public func metricToGregorian(
    year: Int,
    periodType: String,
    periodValue: Int,
    dayOfMonth: Int = 1
) throws -> Date {
    let equinoxYear = year + 1970
    let leap = isLeapYear(year + 1971)
    let yuleDayCount = leap ? 3 : 2

    // offset is 0-indexed days from the equinox (day 0 = Vigil = equinox itself)
    let offset: Int
    switch periodType {
    case "turning":
        guard (0...2).contains(periodValue) else { throw MetricCalendarError.invalidTurningValue }
        offset = periodValue
    case "month":
        guard (1...12).contains(periodValue) else { throw MetricCalendarError.invalidMonth }
        guard (1...30).contains(dayOfMonth) else { throw MetricCalendarError.invalidDay }
        let m = periodValue
        let d = dayOfMonth
        if m <= 9 {
            // Months 1-9 come after The Turning (3 days), then 0-indexed within month
            offset = 3 + (m - 1) * 30 + (d - 1)
        } else {
            // Months 10-12 come after Turning + months 1-9 (270 days) + Yule
            offset = 3 + 270 + yuleDayCount + (m - 10) * 30 + (d - 1)
        }
    case "yule":
        guard (0...2).contains(periodValue) else { throw MetricCalendarError.invalidYuleValue }
        if periodValue == 2 && !leap { throw MetricCalendarError.kindlingRequiresLeapYear }
        // Yule comes after Turning (3) + months 1-9 (270)
        offset = 3 + 270 + periodValue
    default:
        throw MetricCalendarError.unknownPeriodType(periodType)
    }

    let equinoxTs = utcTimestamp(year: equinoxYear, month: 3, day: 20)
    return Date(timeIntervalSince1970: equinoxTs + TimeInterval(offset * 86400))
}

/// Returns true if the given date is a rest day (days 8-10 of any 10-day week).
public func isRestDay(_ date: Date) -> Bool {
    return gregorianToMetric(date).isRest
}

/// Returns the current Metric Calendar date (using UTC clock).
public func today() -> MetricDate {
    return gregorianToMetric(Date())
}
