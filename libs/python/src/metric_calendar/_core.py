from __future__ import annotations
import datetime
from dataclasses import dataclass

_DAY_NAMES = (
    "Primday", "Duoday", "Triday", "Quadday", "Quintday",
    "Hexday", "Septday", "Octday", "Novday", "Decday",
)

_MONTH_NAMES = (
    "Unil", "Duil", "Tril", "Quadril", "Quintil", "Sextil",
    "Septil", "Octil", "Novil", "Decil", "Undecil", "Duodecil",
)

_TURNING_DAY_NAMES = ("Vigil", "Balance", "Dawn")
_YULE_DAY_NAMES = ("Yule Eve", "Midwinter", "Kindling")


@dataclass(frozen=True)
class MetricDate:
    year: int
    month: int          # 1-12, 0 for Turning/Yule
    month_name: str
    day: int            # 1-30, 0 for Turning/Yule
    week_day: int       # 1-10, 0 for Turning/Yule
    day_name: str
    week: int           # 1-36, 0 for Turning/Yule
    season_index: int   # 0-3, -1 for Turning/Yule
    is_leap_year: bool
    is_turning: bool
    is_yule: bool
    is_midsummer: bool  # month == 4 and day == 1
    is_spiral: bool     # month == 5 and day == 18
    is_rest: bool       # week_day >= 8
    special_day: str    # set for Turning and Yule days


def _is_leap_year(y: int) -> bool:
    return y % 4 == 0 and (y % 100 != 0 or y % 400 == 0)


def gregorian_to_metric(date: datetime.date) -> MetricDate:
    year = date.year
    equinox = datetime.date(year, 3, 20)
    days_from_equinox = (date - equinox).days

    if days_from_equinox >= 0:
        metric_year = year - 1970
        day_of_year = days_from_equinox + 1
    else:
        metric_year = year - 1 - 1970
        prev_equinox = datetime.date(year - 1, 3, 20)
        day_of_year = (date - prev_equinox).days + 1

    leap = _is_leap_year(metric_year + 1971)
    yule_day_count = 3 if leap else 2

    base = dict(
        year=metric_year, month=0, month_name="", day=0, week_day=0,
        day_name="", week=0, season_index=-1,
        is_leap_year=leap, is_turning=False, is_yule=False,
        is_midsummer=False, is_spiral=False, is_rest=False, special_day="",
    )

    # The Turning (days 1-3)
    if day_of_year <= 3:
        return MetricDate(**{**base, "is_turning": True, "special_day": _TURNING_DAY_NAMES[day_of_year - 1]})

    adjusted = day_of_year - 3

    if adjusted <= 270:
        m = (adjusted - 1) // 30 + 1
        d = (adjusted - 1) % 30 + 1
    elif adjusted <= 270 + yule_day_count:
        return MetricDate(**{**base, "is_yule": True, "special_day": _YULE_DAY_NAMES[adjusted - 271]})
    else:
        post_yule = adjusted - 270 - yule_day_count
        m = 9 + (post_yule - 1) // 30 + 1
        d = (post_yule - 1) % 30 + 1

    week_day = (d - 1) % 10 + 1
    week = (m - 1) * 3 + (d - 1) // 10 + 1

    return MetricDate(
        year=metric_year,
        month=m, month_name=_MONTH_NAMES[m - 1],
        day=d, week_day=week_day, day_name=_DAY_NAMES[week_day - 1],
        week=week, season_index=(m - 1) // 3,
        is_leap_year=leap,
        is_turning=False, is_yule=False,
        is_midsummer=(m == 4 and d == 1),
        is_spiral=(m == 5 and d == 18),
        is_rest=week_day >= 8,
        special_day="",
    )


def metric_to_gregorian(
    year: int,
    period_type: str,
    period_value: int,
    day_of_month: int = 1,
) -> datetime.date:
    equinox_year = year + 1970
    leap = _is_leap_year(year + 1971)
    yule_day_count = 3 if leap else 2

    if period_type == "turning":
        if not (0 <= period_value <= 2):
            raise ValueError("turning period_value must be 0-2")
        offset = period_value
    elif period_type == "month":
        m, d = period_value, day_of_month
        if not (1 <= m <= 12):
            raise ValueError("month must be 1-12")
        if not (1 <= d <= 30):
            raise ValueError("day must be 1-30")
        if m <= 9:
            offset = 3 + (m - 1) * 30 + (d - 1)
        else:
            offset = 3 + 270 + yule_day_count + (m - 10) * 30 + (d - 1)
    elif period_type == "yule":
        if period_value == 2 and not leap:
            raise ValueError("Kindling only occurs in leap years")
        if not (0 <= period_value <= 2):
            raise ValueError("yule period_value must be 0-2")
        offset = 3 + 270 + period_value
    else:
        raise ValueError(f"Unknown period_type: {period_type!r}")

    equinox = datetime.date(equinox_year, 3, 20)
    return equinox + datetime.timedelta(days=offset)


def is_rest_day(date: datetime.date) -> bool:
    return gregorian_to_metric(date).is_rest


def today() -> MetricDate:
    return gregorian_to_metric(datetime.date.today())
