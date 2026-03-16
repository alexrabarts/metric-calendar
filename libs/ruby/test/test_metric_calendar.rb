require 'minitest/autorun'
require 'date'
require_relative '../lib/metric_calendar'

class TestMetricCalendar < Minitest::Test
  def test_vigil
    r = MetricCalendar.gregorian_to_metric(Date.new(2026, 3, 20))
    assert_equal 56, r.year
    assert r.is_turning
    assert_equal 'Vigil', r.special_day
  end

  def test_balance
    r = MetricCalendar.gregorian_to_metric(Date.new(2026, 3, 21))
    assert r.is_turning
    assert_equal 'Balance', r.special_day
  end

  def test_dawn
    r = MetricCalendar.gregorian_to_metric(Date.new(2026, 3, 22))
    assert r.is_turning
    assert_equal 'Dawn', r.special_day
  end

  def test_unil_1
    r = MetricCalendar.gregorian_to_metric(Date.new(2026, 3, 23))
    assert_equal 56, r.year
    assert_equal 1, r.month
    assert_equal 'Unil', r.month_name
    assert_equal 1, r.day
    assert_equal 1, r.week_day
    assert_equal 'Primday', r.day_name
    refute r.is_rest
    assert_equal 1, r.week
  end

  def test_unil_10
    r = MetricCalendar.gregorian_to_metric(Date.new(2026, 4, 1))
    assert_equal 1, r.month
    assert_equal 10, r.day
    assert_equal 10, r.week_day
    assert_equal 'Decday', r.day_name
    assert r.is_rest
  end

  def test_yule_eve
    r = MetricCalendar.gregorian_to_metric(Date.new(2026, 12, 18))
    assert_equal 56, r.year
    assert r.is_yule
    assert_equal 'Yule Eve', r.special_day
  end

  def test_pre_equinox
    r = MetricCalendar.gregorian_to_metric(Date.new(2025, 1, 1))
    assert_equal 54, r.year  # before March 2025 equinox -> year 54
    assert_equal 10, r.month
    assert_equal 'Decil', r.month_name
    assert_equal 13, r.day
    assert_equal 3, r.week_day
    assert_equal 'Triday', r.day_name
    refute r.is_rest
  end

  def test_midsummer
    r = MetricCalendar.gregorian_to_metric(Date.new(2026, 6, 21))
    assert r.is_midsummer
    assert_equal 4, r.month
    assert_equal 1, r.day
  end

  def test_is_rest_day
    assert MetricCalendar.is_rest_day(Date.new(2026, 4, 1))
    refute MetricCalendar.is_rest_day(Date.new(2026, 3, 23))
  end

  def test_metric_to_gregorian_month
    d = MetricCalendar.metric_to_gregorian(56, 'month', 1, 1)
    assert_equal Date.new(2026, 3, 23), d
  end

  def test_metric_to_gregorian_turning
    d = MetricCalendar.metric_to_gregorian(56, 'turning', 0)
    assert_equal Date.new(2026, 3, 20), d
  end

  def test_metric_to_gregorian_yule
    d = MetricCalendar.metric_to_gregorian(56, 'yule', 0)
    assert_equal Date.new(2026, 12, 18), d
  end

  def test_kindling_requires_leap
    assert_raises(ArgumentError) { MetricCalendar.metric_to_gregorian(56, 'yule', 2) }
  end

  def test_season_index
    r = MetricCalendar.gregorian_to_metric(Date.new(2026, 3, 23))  # Month 1
    assert_equal 0, r.season_index
    r2 = MetricCalendar.gregorian_to_metric(Date.new(2026, 6, 21))  # Month 4
    assert_equal 1, r2.season_index
  end
end
