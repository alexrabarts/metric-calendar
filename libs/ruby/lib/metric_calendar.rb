require 'date'

module MetricCalendar
  DAY_NAMES = %w[Primday Duoday Triday Quadday Quintday Hexday Septday Octday Novday Decday].freeze
  MONTH_NAMES = %w[Unil Duil Tril Quadril Quintil Sextil Septil Octil Novil Decil Undecil Duodecil].freeze
  TURNING_DAY_NAMES = %w[Vigil Balance Dawn].freeze
  YULE_DAY_NAMES = ['Yule Eve', 'Midwinter', 'Kindling'].freeze

  # Struct for a Metric Calendar date. Fields:
  #   year         - Metric year (Year 0 = spring equinox 1970)
  #   month        - 1-12, 0 for Turning/Yule
  #   month_name   - e.g. "Unil", "" for Turning/Yule
  #   day          - 1-30, 0 for Turning/Yule
  #   week_day     - 1-10, 0 for Turning/Yule
  #   day_name     - e.g. "Primday", "" for Turning/Yule
  #   week         - 1-36, 0 for Turning/Yule
  #   season_index - 0-3, -1 for Turning/Yule
  #   is_leap_year - true if this metric year has 3 Yule days
  #   is_turning   - true during The Turning (3 days at spring equinox)
  #   is_yule      - true during Yule
  #   is_midsummer - true on Quadril 1 (summer solstice)
  #   is_spiral    - true on Quintil 18 (golden angle day)
  #   is_rest      - true on days 8-10 of any 10-day week
  #   special_day  - "Vigil", "Balance", "Dawn", "Yule Eve", "Midwinter", "Kindling", or ""
  MetricDate = Struct.new(
    :year, :month, :month_name, :day, :week_day, :day_name,
    :week, :season_index, :is_leap_year, :is_turning, :is_yule,
    :is_midsummer, :is_spiral, :is_rest, :special_day,
    keyword_init: true
  )

  def self.leap_year?(gregorian_year)
    gregorian_year % 4 == 0 && (gregorian_year % 100 != 0 || gregorian_year % 400 == 0)
  end
  private_class_method :leap_year?

  # Convert a Gregorian Date to a MetricDate.
  # @param date [Date] the Gregorian date to convert
  # @return [MetricDate]
  def self.gregorian_to_metric(date)
    year = date.year
    equinox = Date.new(year, 3, 20)
    days_from_equinox = (date - equinox).to_i

    if days_from_equinox >= 0
      metric_year = year - 1970
      day_of_year = days_from_equinox + 1
    else
      metric_year = year - 1 - 1970
      prev_equinox = Date.new(year - 1, 3, 20)
      day_of_year = (date - prev_equinox).to_i + 1
    end

    leap = leap_year?(metric_year + 1971)
    yule_day_count = leap ? 3 : 2

    base = {
      year: metric_year, month: 0, month_name: '', day: 0, week_day: 0,
      day_name: '', week: 0, season_index: -1,
      is_leap_year: leap, is_turning: false, is_yule: false,
      is_midsummer: false, is_spiral: false, is_rest: false, special_day: ''
    }

    # The Turning (days 1-3)
    if day_of_year <= 3
      return MetricDate.new(**base.merge(is_turning: true, special_day: TURNING_DAY_NAMES[day_of_year - 1]))
    end

    adjusted = day_of_year - 3

    if adjusted <= 270
      m = (adjusted - 1) / 30 + 1
      d = (adjusted - 1) % 30 + 1
    elsif adjusted <= 270 + yule_day_count
      return MetricDate.new(**base.merge(is_yule: true, special_day: YULE_DAY_NAMES[adjusted - 271]))
    else
      post_yule = adjusted - 270 - yule_day_count
      m = 9 + (post_yule - 1) / 30 + 1
      d = (post_yule - 1) % 30 + 1
    end

    week_day = (d - 1) % 10 + 1
    week = (m - 1) * 3 + (d - 1) / 10 + 1

    MetricDate.new(
      year: metric_year,
      month: m, month_name: MONTH_NAMES[m - 1],
      day: d, week_day: week_day, day_name: DAY_NAMES[week_day - 1],
      week: week, season_index: (m - 1) / 3,
      is_leap_year: leap,
      is_turning: false, is_yule: false,
      is_midsummer: (m == 4 && d == 1),
      is_spiral: (m == 5 && d == 18),
      is_rest: week_day >= 8,
      special_day: ''
    )
  end

  # Convert a Metric Calendar date back to a Gregorian Date.
  # @param year [Integer] Metric year
  # @param period_type [String] "turning", "month", or "yule"
  # @param period_value [Integer] 0-indexed turning/yule (0-2), or 1-12 for month
  # @param day_of_month [Integer] 1-30, used only when period_type is "month"
  # @return [Date]
  def self.metric_to_gregorian(year, period_type, period_value, day_of_month = 1)
    equinox_year = year + 1970
    leap = leap_year?(year + 1971)
    yule_day_count = leap ? 3 : 2

    offset = case period_type
    when 'turning'
      raise ArgumentError, 'turning period_value must be 0-2' unless (0..2).include?(period_value)
      period_value
    when 'month'
      m, d = period_value, day_of_month
      raise ArgumentError, 'month must be 1-12' unless (1..12).include?(m)
      raise ArgumentError, 'day must be 1-30' unless (1..30).include?(d)
      if m <= 9
        3 + (m - 1) * 30 + (d - 1)
      else
        3 + 270 + yule_day_count + (m - 10) * 30 + (d - 1)
      end
    when 'yule'
      raise ArgumentError, 'Kindling only occurs in leap years' if period_value == 2 && !leap
      raise ArgumentError, 'yule period_value must be 0-2' unless (0..2).include?(period_value)
      3 + 270 + period_value
    else
      raise ArgumentError, "Unknown period_type: #{period_type.inspect}"
    end

    Date.new(equinox_year, 3, 20) + offset
  end

  # Returns true if the given date is a rest day (days 8-10 of any 10-day week).
  # @param date [Date]
  # @return [Boolean]
  def self.is_rest_day(date)
    gregorian_to_metric(date).is_rest
  end

  # Returns the current Metric Calendar date.
  # @return [MetricDate]
  def self.today
    gregorian_to_metric(Date.today)
  end
end
