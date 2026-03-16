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
    isSpiral: boolean;
    isRest: boolean;
    specialDay: string;
}
export declare function gregorianToMetric(date: Date): MetricDate;
export type PeriodType = 'turning' | 'month' | 'yule';
export declare function metricToGregorian(year: number, periodType: PeriodType, periodValue: number, dayOfMonth?: number): Date;
export declare function isRestDay(date: Date): boolean;
export declare function today(): MetricDate;
//# sourceMappingURL=index.d.ts.map