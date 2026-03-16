import { test } from 'node:test'
import assert from 'node:assert/strict'
import { gregorianToMetric, metricToGregorian, isRestDay } from '../dist/index.mjs'

// Helper: create UTC date
function utcDate(y, m, d) {
  return new Date(Date.UTC(y, m - 1, d))
}

test('2026-03-20 → Year 56, Vigil', () => {
  const r = gregorianToMetric(utcDate(2026, 3, 20))
  assert.equal(r.year, 56)
  assert.equal(r.isTurning, true)
  assert.equal(r.specialDay, 'Vigil')
})

test('2026-03-21 → Balance', () => {
  const r = gregorianToMetric(utcDate(2026, 3, 21))
  assert.equal(r.specialDay, 'Balance')
})

test('2026-03-22 → Dawn', () => {
  const r = gregorianToMetric(utcDate(2026, 3, 22))
  assert.equal(r.specialDay, 'Dawn')
})

test('2026-03-23 → Year 56, Unil 1, Primday', () => {
  const r = gregorianToMetric(utcDate(2026, 3, 23))
  assert.equal(r.year, 56)
  assert.equal(r.month, 1)
  assert.equal(r.monthName, 'Unil')
  assert.equal(r.day, 1)
  assert.equal(r.weekDay, 1)
  assert.equal(r.dayName, 'Primday')
  assert.equal(r.isRest, false)
  assert.equal(r.week, 1)
})

test('2026-04-01 → Unil 10, Decday, IsRest', () => {
  const r = gregorianToMetric(utcDate(2026, 4, 1))
  assert.equal(r.month, 1)
  assert.equal(r.day, 10)
  assert.equal(r.weekDay, 10)
  assert.equal(r.dayName, 'Decday')
  assert.equal(r.isRest, true)
})

test('2026-12-18 → Yule Eve', () => {
  const r = gregorianToMetric(utcDate(2026, 12, 18))
  assert.equal(r.year, 56)
  assert.equal(r.isYule, true)
  assert.equal(r.specialDay, 'Yule Eve')
})

test('2025-01-01 → Year 54, Decil 13, Triday', () => {
  const r = gregorianToMetric(utcDate(2025, 1, 1))
  assert.equal(r.year, 54)
  assert.equal(r.month, 10)
  assert.equal(r.monthName, 'Decil')
  assert.equal(r.day, 13)
  assert.equal(r.weekDay, 3)
  assert.equal(r.dayName, 'Triday')
  assert.equal(r.isRest, false)
})

test('Midsummer: 2026-06-21', () => {
  const r = gregorianToMetric(utcDate(2026, 6, 21))
  assert.equal(r.isMidsummer, true)
  assert.equal(r.month, 4)
  assert.equal(r.day, 1)
})

test('isRestDay helper', () => {
  assert.equal(isRestDay(utcDate(2026, 4, 1)), true)
  assert.equal(isRestDay(utcDate(2026, 3, 23)), false)
})

test('metricToGregorian: Year 56, month 1, day 1 → 2026-03-23', () => {
  const d = metricToGregorian(56, 'month', 1, 1)
  assert.equal(d.getUTCFullYear(), 2026)
  assert.equal(d.getUTCMonth() + 1, 3)
  assert.equal(d.getUTCDate(), 23)
})

test('metricToGregorian: Year 56, turning 0 → 2026-03-20', () => {
  const d = metricToGregorian(56, 'turning', 0)
  assert.equal(d.getUTCFullYear(), 2026)
  assert.equal(d.getUTCMonth() + 1, 3)
  assert.equal(d.getUTCDate(), 20)
})

test('metricToGregorian: Year 56, yule 0 → 2026-12-18', () => {
  const d = metricToGregorian(56, 'yule', 0)
  assert.equal(d.getUTCFullYear(), 2026)
  assert.equal(d.getUTCMonth() + 1, 12)
  assert.equal(d.getUTCDate(), 18)
})
