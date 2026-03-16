// src/index.ts
var DAY_NAMES = [
  "Primday",
  "Duoday",
  "Triday",
  "Quadday",
  "Quintday",
  "Hexday",
  "Septday",
  "Octday",
  "Novday",
  "Decday"
];
var MONTH_NAMES = [
  "Unil",
  "Duil",
  "Tril",
  "Quadril",
  "Quintil",
  "Sextil",
  "Septil",
  "Octil",
  "Novil",
  "Decil",
  "Undecil",
  "Duodecil"
];
var TURNING_DAY_NAMES = ["Vigil", "Balance", "Dawn"];
var YULE_DAY_NAMES = ["Yule Eve", "Midwinter", "Kindling"];
function gregorianToMetric(date) {
  const year = date.getUTCFullYear();
  const month = date.getUTCMonth() + 1;
  const day = date.getUTCDate();
  const equinoxMs = Date.UTC(year, 2, 20);
  const dateMs = Date.UTC(year, month - 1, day);
  const daysFromEquinox = Math.floor((dateMs - equinoxMs) / 864e5);
  let metricYear;
  let dayOfYear;
  if (daysFromEquinox >= 0) {
    metricYear = year - 1970;
    dayOfYear = daysFromEquinox + 1;
  } else {
    metricYear = year - 1 - 1970;
    const prevEquinoxMs = Date.UTC(year - 1, 2, 20);
    dayOfYear = Math.floor((dateMs - prevEquinoxMs) / 864e5) + 1;
  }
  const leap = isLeapYear(metricYear + 1971);
  const yuleDayCount = leap ? 3 : 2;
  const base = {
    year: metricYear,
    month: 0,
    monthName: "",
    day: 0,
    weekDay: 0,
    dayName: "",
    week: 0,
    seasonIndex: -1,
    isLeapYear: leap,
    isTurning: false,
    isYule: false,
    isMidsummer: false,
    isSpiral: false,
    isRest: false,
    specialDay: ""
  };
  if (dayOfYear <= 3) {
    return { ...base, isTurning: true, specialDay: TURNING_DAY_NAMES[dayOfYear - 1] };
  }
  const adjusted = dayOfYear - 3;
  let m = 0, d = 0;
  if (adjusted <= 270) {
    m = Math.ceil(adjusted / 30);
    d = (adjusted - 1) % 30 + 1;
  } else if (adjusted <= 270 + yuleDayCount) {
    return { ...base, isYule: true, specialDay: YULE_DAY_NAMES[adjusted - 271] };
  } else {
    const postYule = adjusted - 270 - yuleDayCount;
    m = 9 + Math.ceil(postYule / 30);
    d = (postYule - 1) % 30 + 1;
  }
  const weekDay = (d - 1) % 10 + 1;
  const week = (m - 1) * 3 + Math.floor((d - 1) / 10) + 1;
  return {
    ...base,
    month: m,
    monthName: MONTH_NAMES[m - 1],
    day: d,
    weekDay,
    dayName: DAY_NAMES[weekDay - 1],
    week,
    seasonIndex: Math.floor((m - 1) / 3),
    isRest: weekDay >= 8,
    isMidsummer: m === 4 && d === 1,
    isSpiral: m === 5 && d === 18
  };
}
function metricToGregorian(year, periodType, periodValue, dayOfMonth = 1) {
  const equinoxYear = year + 1970;
  const leap = isLeapYear(year + 1971);
  const yuleDayCount = leap ? 3 : 2;
  let offset;
  if (periodType === "turning") {
    if (periodValue < 0 || periodValue > 2) throw new Error("turning period value must be 0-2");
    offset = periodValue;
  } else if (periodType === "month") {
    const m = periodValue;
    const d = dayOfMonth;
    if (m < 1 || m > 12) throw new Error("month must be 1-12");
    if (d < 1 || d > 30) throw new Error("day must be 1-30");
    if (m <= 9) {
      offset = 3 + (m - 1) * 30 + (d - 1);
    } else {
      offset = 3 + 270 + yuleDayCount + (m - 10) * 30 + (d - 1);
    }
  } else {
    if (periodValue === 2 && !leap) throw new Error("Kindling only occurs in leap years");
    if (periodValue < 0 || periodValue > 2) throw new Error("yule period value must be 0-2");
    offset = 3 + 270 + periodValue;
  }
  return new Date(Date.UTC(equinoxYear, 2, 20 + offset));
}
function isRestDay(date) {
  return gregorianToMetric(date).isRest;
}
function today() {
  return gregorianToMetric(/* @__PURE__ */ new Date());
}
function isLeapYear(y) {
  return y % 4 === 0 && (y % 100 !== 0 || y % 400 === 0);
}
export {
  gregorianToMetric,
  isRestDay,
  metricToGregorian,
  today
};
