export interface MetricDate {
    year: number;
    month: number;
    monthName: string;
    day: number;
    weekDay: number;
    dayName: string;
    week: number;
    seasonIndex: number;
    isLeapYear: boolean;
    isTurning: boolean;
    isYule: boolean;
    isMidsummer: boolean;
    isSextant: boolean;
    isTrine: boolean;
    isSpiral: boolean;
    isConvergence: boolean;
    isMeridian: boolean;
    isMask: boolean;
    isHarmony: boolean;
    isRest: boolean;
    specialDay: string;
    observance: string;
}
export declare function gregorianToMetric(date: Date): MetricDate;
/**
 * Format a MetricDate using a pattern string.
 *
 * Tokens:
 *   MMM  month name (e.g. "Unil")
 *   MM   month zero-padded (e.g. "01")
 *   M    month number (e.g. "1")
 *   DD   day zero-padded (e.g. "04")
 *   D    day number (e.g. "4")
 *   WW   weekday name (e.g. "Quintday")
 *   W    weekday number (e.g. "5")
 *   Y    year number (e.g. "56")
 *   S    season name (e.g. "Rising")
 *
 * Example: format(d, "WW, MMM D, Year Y") → "Quintday, Unil 4, Year 56"
 */
export declare function format(date: MetricDate, pattern: string): string;
export type PeriodType = 'turning' | 'month' | 'yule';
export declare function metricToGregorian(year: number, periodType: PeriodType, periodValue: number, dayOfMonth?: number): Date;
export declare function isRestDay(date: Date): boolean;
export declare function today(): MetricDate;
//# sourceMappingURL=index.d.ts.map