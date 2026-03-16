import unittest
import datetime
import sys, os
sys.path.insert(0, os.path.join(os.path.dirname(__file__), '..', 'src'))
from metric_calendar import gregorian_to_metric, metric_to_gregorian, is_rest_day

class TestGregorianToMetric(unittest.TestCase):

    def test_vigil(self):
        r = gregorian_to_metric(datetime.date(2026, 3, 20))
        self.assertEqual(r.year, 56)
        self.assertTrue(r.is_turning)
        self.assertEqual(r.special_day, 'Vigil')

    def test_balance(self):
        r = gregorian_to_metric(datetime.date(2026, 3, 21))
        self.assertTrue(r.is_turning)
        self.assertEqual(r.special_day, 'Balance')

    def test_dawn(self):
        r = gregorian_to_metric(datetime.date(2026, 3, 22))
        self.assertTrue(r.is_turning)
        self.assertEqual(r.special_day, 'Dawn')

    def test_unil_1(self):
        r = gregorian_to_metric(datetime.date(2026, 3, 23))
        self.assertEqual(r.year, 56)
        self.assertEqual(r.month, 1)
        self.assertEqual(r.month_name, 'Unil')
        self.assertEqual(r.day, 1)
        self.assertEqual(r.week_day, 1)
        self.assertEqual(r.day_name, 'Primday')
        self.assertFalse(r.is_rest)
        self.assertEqual(r.week, 1)

    def test_unil_10(self):
        r = gregorian_to_metric(datetime.date(2026, 4, 1))
        self.assertEqual(r.month, 1)
        self.assertEqual(r.day, 10)
        self.assertEqual(r.week_day, 10)
        self.assertEqual(r.day_name, 'Decday')
        self.assertTrue(r.is_rest)

    def test_yule_eve(self):
        r = gregorian_to_metric(datetime.date(2026, 12, 18))
        self.assertEqual(r.year, 56)
        self.assertTrue(r.is_yule)
        self.assertEqual(r.special_day, 'Yule Eve')

    def test_pre_equinox(self):
        # Jan 1 2025 is before the March 20 2025 equinox, so it belongs to
        # Metric year 54 (which runs from March 20 2024 to March 19 2025).
        r = gregorian_to_metric(datetime.date(2025, 1, 1))
        self.assertEqual(r.year, 54)
        self.assertEqual(r.month, 10)
        self.assertEqual(r.month_name, 'Decil')
        self.assertEqual(r.day, 13)
        self.assertEqual(r.week_day, 3)
        self.assertEqual(r.day_name, 'Triday')
        self.assertFalse(r.is_rest)

    def test_midsummer(self):
        r = gregorian_to_metric(datetime.date(2026, 6, 21))
        self.assertTrue(r.is_midsummer)
        self.assertEqual(r.month, 4)
        self.assertEqual(r.day, 1)

    def test_is_rest_day(self):
        self.assertTrue(is_rest_day(datetime.date(2026, 4, 1)))
        self.assertFalse(is_rest_day(datetime.date(2026, 3, 23)))

    def test_metric_to_gregorian_month(self):
        d = metric_to_gregorian(56, 'month', 1, 1)
        self.assertEqual(d, datetime.date(2026, 3, 23))

    def test_metric_to_gregorian_turning(self):
        d = metric_to_gregorian(56, 'turning', 0)
        self.assertEqual(d, datetime.date(2026, 3, 20))

    def test_metric_to_gregorian_yule(self):
        d = metric_to_gregorian(56, 'yule', 0)
        self.assertEqual(d, datetime.date(2026, 12, 18))

    def test_kindling_requires_leap(self):
        with self.assertRaises(ValueError):
            metric_to_gregorian(56, 'yule', 2)  # Year 56 is not a leap year

if __name__ == '__main__':
    unittest.main()
