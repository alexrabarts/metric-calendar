const DAY_NAMES: readonly string[] = [
  'Primday', 'Duoday', 'Triday', 'Quadday', 'Quintday',
  'Hexday', 'Septday', 'Octday', 'Novday', 'Decday',
]

const MONTH_NAMES: readonly string[] = [
  'Unil', 'Duil', 'Tril', 'Quadril', 'Quintil', 'Sextil',
  'Septil', 'Octil', 'Novil', 'Decil', 'Undecil', 'Duodecil',
]

const TURNING_DAY_NAMES: readonly string[] = ['Vigil', 'Balance', 'Dawn']
const YULE_DAY_NAMES: readonly string[] = ['Yule Eve', 'Midwinter', 'Kindling']

export interface MetricDate {
  year: number
  month: number        // 1-12, 0 for Turning/Yule
  monthName: string
  day: number          // 1-30, 0 for Turning/Yule
  weekDay: number      // 1-10, 0 for Turning/Yule
  dayName: string
  week: number         // 1-36, 0 for Turning/Yule
  seasonIndex: number  // 0-3, -1 for Turning/Yule
  isLeapYear: boolean
  isTurning: boolean
  isYule: boolean
  isMidsummer: boolean  // month === 4 && day === 1
  isSpiral: boolean     // month === 5 && day === 18
  isRest: boolean       // weekDay >= 8
  specialDay: string    // set for Turning and Yule days
}

export function gregorianToMetric(date: Date): MetricDate {
  const year = date.getUTCFullYear()
  const month = date.getUTCMonth() + 1
  const day = date.getUTCDate()

  // Days from equinox
  const equinoxMs = Date.UTC(year, 2, 20)  // March 20
  const dateMs = Date.UTC(year, month - 1, day)
  const daysFromEquinox = Math.floor((dateMs - equinoxMs) / 86400000)

  let metricYear: number
  let dayOfYear: number

  if (daysFromEquinox >= 0) {
    metricYear = year - 1970
    dayOfYear = daysFromEquinox + 1
  } else {
    metricYear = year - 1 - 1970
    const prevEquinoxMs = Date.UTC(year - 1, 2, 20)
    dayOfYear = Math.floor((dateMs - prevEquinoxMs) / 86400000) + 1
  }

  const leap = isLeapYear(metricYear + 1971)
  const yuleDayCount = leap ? 3 : 2

  // Base result
  const base: MetricDate = {
    year: metricYear, month: 0, monthName: '', day: 0, weekDay: 0,
    dayName: '', week: 0, seasonIndex: -1,
    isLeapYear: leap, isTurning: false, isYule: false,
    isMidsummer: false, isSpiral: false, isRest: false, specialDay: '',
  }

  // The Turning (days 1-3)
  if (dayOfYear <= 3) {
    return { ...base, isTurning: true, specialDay: TURNING_DAY_NAMES[dayOfYear - 1]! }
  }

  const adjusted = dayOfYear - 3
  let m = 0, d = 0

  if (adjusted <= 270) {
    m = Math.ceil(adjusted / 30)
    d = ((adjusted - 1) % 30) + 1
  } else if (adjusted <= 270 + yuleDayCount) {
    return { ...base, isYule: true, specialDay: YULE_DAY_NAMES[adjusted - 271]! }
  } else {
    const postYule = adjusted - 270 - yuleDayCount
    m = 9 + Math.ceil(postYule / 30)
    d = ((postYule - 1) % 30) + 1
  }

  const weekDay = ((d - 1) % 10) + 1
  const week = (m - 1) * 3 + Math.floor((d - 1) / 10) + 1
  return {
    ...base,
    month: m, monthName: MONTH_NAMES[m - 1]!,
    day: d, weekDay, dayName: DAY_NAMES[weekDay - 1]!,
    week, seasonIndex: Math.floor((m - 1) / 3),
    isRest: weekDay >= 8,
    isMidsummer: m === 4 && d === 1,
    isSpiral: m === 5 && d === 18,
  }
}

export type PeriodType = 'turning' | 'month' | 'yule'

export function metricToGregorian(
  year: number, periodType: PeriodType, periodValue: number, dayOfMonth: number = 1
): Date {
  const equinoxYear = year + 1970
  const leap = isLeapYear(year + 1971)
  const yuleDayCount = leap ? 3 : 2
  let offset: number

  if (periodType === 'turning') {
    if (periodValue < 0 || periodValue > 2) throw new Error('turning period value must be 0-2')
    offset = periodValue
  } else if (periodType === 'month') {
    const m = periodValue
    const d = dayOfMonth
    if (m < 1 || m > 12) throw new Error('month must be 1-12')
    if (d < 1 || d > 30) throw new Error('day must be 1-30')
    if (m <= 9) {
      offset = 3 + (m - 1) * 30 + (d - 1)
    } else {
      offset = 3 + 270 + yuleDayCount + (m - 10) * 30 + (d - 1)
    }
  } else {
    // yule
    if (periodValue === 2 && !leap) throw new Error('Kindling only occurs in leap years')
    if (periodValue < 0 || periodValue > 2) throw new Error('yule period value must be 0-2')
    offset = 3 + 270 + periodValue
  }

  return new Date(Date.UTC(equinoxYear, 2, 20 + offset))
}

export function isRestDay(date: Date): boolean {
  return gregorianToMetric(date).isRest
}

export function today(): MetricDate {
  return gregorianToMetric(new Date())
}

function isLeapYear(y: number): boolean {
  return y % 4 === 0 && (y % 100 !== 0 || y % 400 === 0)
}
