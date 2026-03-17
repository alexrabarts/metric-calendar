"use strict";
var MetricCalendar = (() => {
  var __defProp = Object.defineProperty;
  var __getOwnPropDesc = Object.getOwnPropertyDescriptor;
  var __getOwnPropNames = Object.getOwnPropertyNames;
  var __hasOwnProp = Object.prototype.hasOwnProperty;
  var __export = (target, all) => {
    for (var name in all)
      __defProp(target, name, { get: all[name], enumerable: true });
  };
  var __copyProps = (to, from, except, desc) => {
    if (from && typeof from === "object" || typeof from === "function") {
      for (let key of __getOwnPropNames(from))
        if (!__hasOwnProp.call(to, key) && key !== except)
          __defProp(to, key, { get: () => from[key], enumerable: !(desc = __getOwnPropDesc(from, key)) || desc.enumerable });
    }
    return to;
  };
  var __toCommonJS = (mod) => __copyProps(__defProp({}, "__esModule", { value: true }), mod);

  // src/index.ts
  var index_exports = {};
  __export(index_exports, {
    format: () => format,
    gregorianToMetric: () => gregorianToMetric,
    isRestDay: () => isRestDay,
    metricToGregorian: () => metricToGregorian,
    today: () => today
  });
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
  var SEASON_NAMES = ["Rising", "Flourishing", "Gathering", "Stillness"];
  var TURNING_DAY_NAMES = ["Vigil", "Balance", "Dawn"];
  var YULE_DAY_NAMES = ["Yule Eve", "Midwinter", "Kindling"];
  var FORMAT_RE = /MMM|MM|M|DD|D|WW|W|Y|S/g;
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
      isSextant: false,
      isTrine: false,
      isSpiral: false,
      isConvergence: false,
      isMeridian: false,
      isMask: false,
      isHarmony: false,
      isRest: false,
      specialDay: "",
      observance: ""
    };
    if (dayOfYear <= 3) {
      const specialDay = TURNING_DAY_NAMES[dayOfYear - 1];
      return { ...base, isTurning: true, specialDay, observance: specialDay };
    }
    const adjusted = dayOfYear - 3;
    let m = 0, d = 0;
    if (adjusted <= 270) {
      m = Math.ceil(adjusted / 30);
      d = (adjusted - 1) % 30 + 1;
    } else if (adjusted <= 270 + yuleDayCount) {
      const specialDay = YULE_DAY_NAMES[adjusted - 271];
      return { ...base, isYule: true, specialDay, observance: specialDay };
    } else {
      const postYule = adjusted - 270 - yuleDayCount;
      m = 9 + Math.ceil(postYule / 30);
      d = (postYule - 1) % 30 + 1;
    }
    const weekDay = (d - 1) % 10 + 1;
    const week = (m - 1) * 3 + Math.floor((d - 1) / 10) + 1;
    const isMidsummer = m === 4 && d === 1;
    const isSextant = m === 2 && d === 30;
    const isTrine = m === 4 && d === 30;
    const isSpiral = m === 5 && d === 18;
    const isConvergence = m === 5 && d === 24;
    const isMeridian = m === 6 && d === 30;
    const isMask = m === 8 && d === 13;
    const isHarmony = m === 8 && d === 30;
    const observance = isMidsummer ? "Midsummer" : isSextant ? "The Sextant" : isTrine ? "The Trine" : isSpiral ? "The Spiral" : isConvergence ? "Convergence" : isMeridian ? "The Meridian" : isMask ? "The Mask" : isHarmony ? "Harmony" : "";
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
      isMidsummer,
      isSextant,
      isTrine,
      isSpiral,
      isConvergence,
      isMeridian,
      isMask,
      isHarmony,
      observance
    };
  }
  function format(date, pattern) {
    const seasonName = date.seasonIndex >= 0 ? SEASON_NAMES[date.seasonIndex] : "";
    const tokens = {
      "MMM": date.monthName,
      "MM": String(date.month).padStart(2, "0"),
      "M": String(date.month),
      "DD": String(date.day).padStart(2, "0"),
      "D": String(date.day),
      "WW": date.dayName,
      "W": String(date.weekDay),
      "Y": String(date.year),
      "S": seasonName
    };
    return pattern.replace(FORMAT_RE, (tok) => tokens[tok] ?? tok);
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
  return __toCommonJS(index_exports);
})();
