//! Metric Calendar — a precision decimal calendar system.
//!
//! Year 0 = spring equinox 1970. Each year begins at the spring equinox (March 20).

const DAY_NAMES: [&str; 10] = [
    "Primday", "Duoday", "Triday", "Quadday", "Quintday",
    "Hexday", "Septday", "Octday", "Novday", "Decday",
];

const MONTH_NAMES: [&str; 12] = [
    "Unil", "Duil", "Tril", "Quadril", "Quintil", "Sextil",
    "Septil", "Octil", "Novil", "Decil", "Undecil", "Duodecil",
];

const TURNING_DAY_NAMES: [&str; 3] = ["Vigil", "Balance", "Dawn"];
const YULE_DAY_NAMES: [&str; 3] = ["Yule Eve", "Midwinter", "Kindling"];

/// A date in the Metric Calendar system.
#[derive(Debug, Clone, PartialEq, Eq)]
pub struct MetricDate {
    pub year: i32,
    /// 1-12 for regular days, 0 for Turning/Yule
    pub month: u8,
    pub month_name: &'static str,
    /// 1-30 for regular days, 0 for Turning/Yule
    pub day: u8,
    /// 1-10 for regular days, 0 for Turning/Yule
    pub week_day: u8,
    pub day_name: &'static str,
    /// 1-36 for regular days, 0 for Turning/Yule
    pub week: u8,
    /// 0-3 for regular days, -1 for Turning/Yule (use i8)
    pub season_index: i8,
    pub is_leap_year: bool,
    /// true during The Turning (3 days at spring equinox)
    pub is_turning: bool,
    /// true during Yule
    pub is_yule: bool,
    /// true on Quadril 1 (summer solstice)
    pub is_midsummer: bool,
    /// true on Quintil 18 (golden angle day)
    pub is_spiral: bool,
    /// true on days 8-10 of any 10-day week
    pub is_rest: bool,
    /// name for Turning or Yule days, empty string otherwise
    pub special_day: &'static str,
}

fn is_leap_year(y: i32) -> bool {
    y % 4 == 0 && (y % 100 != 0 || y % 400 == 0)
}

/// Compute Julian Day Number for a Gregorian calendar date.
/// This allows integer day arithmetic without external dependencies.
fn julian_day_number(year: i32, month: u32, day: u32) -> i32 {
    let a = (14 - month as i32) / 12;
    let y = year + 4800 - a;
    let m = month as i32 + 12 * a - 3;
    day as i32 + (153 * m + 2) / 5 + 365 * y + y / 4 - y / 100 + y / 400 - 32045
}

/// Convert a Gregorian date to a Metric Calendar date.
///
/// # Arguments
/// * `year` - Gregorian year
/// * `month` - Gregorian month (1-12)
/// * `day` - Gregorian day (1-31)
pub fn gregorian_to_metric(year: i32, month: u32, day: u32) -> MetricDate {
    let date_jdn = julian_day_number(year, month, day);
    let equinox_jdn = julian_day_number(year, 3, 20);
    let days_from_equinox = date_jdn - equinox_jdn;

    let (metric_year, day_of_year) = if days_from_equinox >= 0 {
        (year - 1970, days_from_equinox + 1)
    } else {
        let prev_equinox_jdn = julian_day_number(year - 1, 3, 20);
        (year - 1 - 1970, date_jdn - prev_equinox_jdn + 1)
    };

    let leap = is_leap_year(metric_year + 1971);
    let yule_day_count = if leap { 3 } else { 2 };

    let base = MetricDate {
        year: metric_year,
        month: 0, month_name: "",
        day: 0, week_day: 0, day_name: "",
        week: 0, season_index: -1,
        is_leap_year: leap,
        is_turning: false, is_yule: false,
        is_midsummer: false, is_spiral: false, is_rest: false,
        special_day: "",
    };

    // The Turning (days 1-3)
    if day_of_year <= 3 {
        return MetricDate {
            is_turning: true,
            special_day: TURNING_DAY_NAMES[(day_of_year - 1) as usize],
            ..base
        };
    }

    let adjusted = day_of_year - 3;

    let (m, d) = if adjusted <= 270 {
        (((adjusted - 1) / 30 + 1) as u8, ((adjusted - 1) % 30 + 1) as u8)
    } else if adjusted <= 270 + yule_day_count {
        return MetricDate {
            is_yule: true,
            special_day: YULE_DAY_NAMES[(adjusted - 271) as usize],
            ..base
        };
    } else {
        let post_yule = adjusted - 270 - yule_day_count;
        ((9 + (post_yule - 1) / 30 + 1) as u8, ((post_yule - 1) % 30 + 1) as u8)
    };

    let week_day = ((d - 1) % 10 + 1) as u8;
    let week = ((m as i32 - 1) * 3 + (d as i32 - 1) / 10 + 1) as u8;

    MetricDate {
        year: metric_year,
        month: m,
        month_name: MONTH_NAMES[(m - 1) as usize],
        day: d,
        week_day,
        day_name: DAY_NAMES[(week_day - 1) as usize],
        week,
        season_index: ((m as i32 - 1) / 3) as i8,
        is_leap_year: leap,
        is_turning: false,
        is_yule: false,
        is_midsummer: m == 4 && d == 1,
        is_spiral: m == 5 && d == 18,
        is_rest: week_day >= 8,
        special_day: "",
    }
}

/// Convert a Metric Calendar date back to a Gregorian date.
///
/// Returns `(year, month, day)` as a tuple, or `Err` if the input is invalid.
///
/// # Arguments
/// * `year` - Metric year
/// * `period_type` - "turning", "month", or "yule"
/// * `period_value` - 0-indexed: turning (0-2), month (1-12), yule (0-2)
/// * `day_of_month` - 1-30, used only when period_type is "month"
pub fn metric_to_gregorian(
    year: i32,
    period_type: &str,
    period_value: i32,
    day_of_month: u32,
) -> Result<(i32, u32, u32), String> {
    let equinox_year = year + 1970;
    let leap = is_leap_year(year + 1971);
    let yule_day_count = if leap { 3 } else { 2 };

    let offset = match period_type {
        "turning" => {
            if !(0..=2).contains(&period_value) {
                return Err("turning period_value must be 0-2".to_string());
            }
            period_value
        }
        "month" => {
            let m = period_value;
            let d = day_of_month as i32;
            if !(1..=12).contains(&m) {
                return Err("month must be 1-12".to_string());
            }
            if !(1..=30).contains(&d) {
                return Err("day must be 1-30".to_string());
            }
            if m <= 9 {
                3 + (m - 1) * 30 + (d - 1)
            } else {
                3 + 270 + yule_day_count + (m - 10) * 30 + (d - 1)
            }
        }
        "yule" => {
            if period_value == 2 && !leap {
                return Err("Kindling only occurs in leap years".to_string());
            }
            if !(0..=2).contains(&period_value) {
                return Err("yule period_value must be 0-2".to_string());
            }
            3 + 270 + period_value
        }
        _ => return Err(format!("Unknown period_type: {}", period_type)),
    };

    // Start from March 20 of equinox year, add offset days
    // Use JDN arithmetic to convert back to y/m/d
    let equinox_jdn = julian_day_number(equinox_year, 3, 20);
    let target_jdn = equinox_jdn + offset;
    let (gy, gm, gd) = jdn_to_gregorian(target_jdn);
    Ok((gy, gm, gd))
}

/// Convert Julian Day Number back to Gregorian date.
fn jdn_to_gregorian(jdn: i32) -> (i32, u32, u32) {
    let a = jdn + 32044;
    let b = (4 * a + 3) / 146097;
    let c = a - (146097 * b) / 4;
    let d = (4 * c + 3) / 1461;
    let e = c - (1461 * d) / 4;
    let m = (5 * e + 2) / 153;
    let day = (e - (153 * m + 2) / 5 + 1) as u32;
    let month = (m + 3 - 12 * (m / 10)) as u32;
    let year = (100 * b + d - 4800 + m / 10) as i32;
    (year, month, day)
}

/// Returns true if the given date falls on a rest day (days 8-10 of any 10-day week).
pub fn is_rest_day(year: i32, month: u32, day: u32) -> bool {
    gregorian_to_metric(year, month, day).is_rest
}

/// Returns the current Metric Calendar date using the system clock (UTC).
pub fn today() -> MetricDate {
    // Get current UTC date from system time
    let secs = std::time::SystemTime::now()
        .duration_since(std::time::UNIX_EPOCH)
        .expect("system time before epoch")
        .as_secs();

    // Convert Unix timestamp to UTC date
    let days_since_epoch = (secs / 86400) as i64;

    // Use JDN: Unix epoch (1970-01-01) = JDN 2440588
    let jdn = (days_since_epoch + 2440588) as i32;
    let (y, m, d) = jdn_to_gregorian(jdn);
    gregorian_to_metric(y, m, d)
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_vigil() {
        let r = gregorian_to_metric(2026, 3, 20);
        assert_eq!(r.year, 56);
        assert!(r.is_turning);
        assert_eq!(r.special_day, "Vigil");
    }

    #[test]
    fn test_balance() {
        let r = gregorian_to_metric(2026, 3, 21);
        assert!(r.is_turning);
        assert_eq!(r.special_day, "Balance");
    }

    #[test]
    fn test_dawn() {
        let r = gregorian_to_metric(2026, 3, 22);
        assert!(r.is_turning);
        assert_eq!(r.special_day, "Dawn");
    }

    #[test]
    fn test_unil_1() {
        let r = gregorian_to_metric(2026, 3, 23);
        assert_eq!(r.year, 56);
        assert_eq!(r.month, 1);
        assert_eq!(r.month_name, "Unil");
        assert_eq!(r.day, 1);
        assert_eq!(r.week_day, 1);
        assert_eq!(r.day_name, "Primday");
        assert!(!r.is_rest);
        assert_eq!(r.week, 1);
    }

    #[test]
    fn test_unil_10() {
        let r = gregorian_to_metric(2026, 4, 1);
        assert_eq!(r.month, 1);
        assert_eq!(r.day, 10);
        assert_eq!(r.week_day, 10);
        assert_eq!(r.day_name, "Decday");
        assert!(r.is_rest);
    }

    #[test]
    fn test_yule_eve() {
        let r = gregorian_to_metric(2026, 12, 18);
        assert_eq!(r.year, 56);
        assert!(r.is_yule);
        assert_eq!(r.special_day, "Yule Eve");
    }

    #[test]
    fn test_pre_equinox() {
        let r = gregorian_to_metric(2025, 1, 1);
        assert_eq!(r.year, 54);
        assert_eq!(r.month, 10);
        assert_eq!(r.month_name, "Decil");
        assert_eq!(r.day, 13);
        assert_eq!(r.week_day, 3);
        assert_eq!(r.day_name, "Triday");
        assert!(!r.is_rest);
    }

    #[test]
    fn test_midsummer() {
        let r = gregorian_to_metric(2026, 6, 21);
        assert!(r.is_midsummer);
        assert_eq!(r.month, 4);
        assert_eq!(r.day, 1);
    }

    #[test]
    fn test_is_rest_day() {
        assert!(is_rest_day(2026, 4, 1));
        assert!(!is_rest_day(2026, 3, 23));
    }

    #[test]
    fn test_metric_to_gregorian_month() {
        let (y, m, d) = metric_to_gregorian(56, "month", 1, 1).unwrap();
        assert_eq!((y, m, d), (2026, 3, 23));
    }

    #[test]
    fn test_metric_to_gregorian_turning() {
        let (y, m, d) = metric_to_gregorian(56, "turning", 0, 1).unwrap();
        assert_eq!((y, m, d), (2026, 3, 20));
    }

    #[test]
    fn test_metric_to_gregorian_yule() {
        let (y, m, d) = metric_to_gregorian(56, "yule", 0, 1).unwrap();
        assert_eq!((y, m, d), (2026, 12, 18));
    }

    #[test]
    fn test_kindling_requires_leap() {
        assert!(metric_to_gregorian(56, "yule", 2, 1).is_err());
    }

    #[test]
    fn test_leap_year_57() {
        // Metric year 57 uses Gregorian leap year 2028 for its Yule day count
        // (is_leap_year checks metric_year + 1971 = 57 + 1971 = 2028, which is a leap year)
        // so metric year 57 has 3 Yule days and Kindling exists
        let kindling = metric_to_gregorian(57, "yule", 2, 1);
        assert!(kindling.is_ok());
    }
}
