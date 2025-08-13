package util

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Time is database value to parse data from database and parset time.Time to timestamp on json mashal
type Time time.Time

// TimeUnix will return time by timestamp
func TimeUnix(timestamp int64) Time {
	return Time(time.Unix(0, timestamp*1e6))
}

// TimeZero will return zero time
func TimeZero() Time {
	return Time(time.Unix(0, 0*1e6))
}

// TimeNow return current Time
func TimeNow() Time {
	return Time(Now())
}

// TimeStartOfToday return 00:00:00 of today
func TimeStartOfToday() Time {
	now := Now()
	return Time(time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()))
}

// TimeStartOfMonth return 00:00:00 of today
func TimeStartOfMonth() Time {
	now := Now()
	return Time(time.Date(now.Year(), now.Month(), 0, 0, 0, 0, 0, now.Location()))
}

// Timestamp return timestamp
func (t Time) Timestamp() int64 {
	return time.Time(t).Local().UnixNano() / 1e6
}

// MarshalJSON marshal time to string
func (t *Time) MarshalJSON() ([]byte, error) {
	raw := t.Timestamp()
	if raw < 0 {
		return []byte("null"), nil
	}
	stamp := fmt.Sprintf("%v", raw)
	return []byte(stamp), nil
}

// UnmarshalJSON unmarshal string to time
func (t *Time) UnmarshalJSON(bys []byte) (err error) {
	val := strings.TrimSpace(string(bys))
	if val == "null" {
		return
	}
	timestamp, err := strconv.ParseInt(val, 10, 64)
	if err == nil {
		*t = Time(time.Unix(0, timestamp*1e6))
	}
	return
}

// Scan is sql.Sanner
func (t *Time) Scan(src interface{}) (err error) {
	if src != nil {
		if timeSrc, ok := src.(time.Time); ok {
			*t = Time(timeSrc)
		}
	}
	return
}

// NowMilli 获得当前时间，毫秒级
func NowMilli() int64 {
	return Now().Local().UnixNano() / int64(time.Millisecond)
}

// 当天时间的0点
func DayZeroStart(add time.Duration) int64 {
	now := Now()
	startTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	startTime = startTime.Add(add)
	return startTime.UnixNano() / 1e6
}

// 当天时间的0点
func DayZeroStartTime(add time.Duration) time.Time {
	now := Now()
	startTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	startTime = startTime.Add(add)
	return startTime
}

func Now13() int64 {
	return Now().UnixNano() / 1e6
}

func Now10() int64 {
	return Now().Unix()
}

func LoadLocation() *time.Location {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	return loc
}

func Now() time.Time {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	now := time.Now().In(loc)
	return now
}

func WeekStart() time.Time {
	now := Now()
	offset := int(now.Weekday())
	if offset == 0 {
		offset = 7
	}
	weekStart := now.AddDate(0, 0, -offset+1)
	return time.Date(weekStart.Year(), weekStart.Month(), weekStart.Day(), 0, 0, 0, 0, now.Location())
}

func DayStart() time.Time {
	now := Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
}

func DayEnd() time.Time {
	now := Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, now.Location())
}

func WeekEnd() time.Time {
	now := Now()
	offset := int(now.Weekday())
	if offset == 0 {
		offset = 7
	}
	weekEnd := now.AddDate(0, 0, 7-offset)
	return time.Date(weekEnd.Year(), weekEnd.Month(), weekEnd.Day(), 23, 59, 59, 999999999, now.Location())
}

func MonthStart() time.Time {
	now := Now()
	return time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
}

func MonthEnd() time.Time {
	now := Now()
	return time.Date(now.Year(), now.Month()+1, 0, 23, 59, 59, 999999999, now.Location())
}

func HumanDurationZH(d time.Duration) string {
	if d <= 0 {
		return "0秒"
	}

	days := int(d / (24 * time.Hour))
	d -= time.Duration(days) * 24 * time.Hour

	hours := int(d / time.Hour)
	d -= time.Duration(hours) * time.Hour

	minutes := int(d / time.Minute)
	d -= time.Duration(minutes) * time.Minute

	seconds := int(d / time.Second)

	result := ""
	if days > 0 {
		result += fmt.Sprintf("%d天", days)
	}
	if hours > 0 {
		result += fmt.Sprintf("%d小时", hours)
	}
	if minutes > 0 {
		result += fmt.Sprintf("%d分钟", minutes)
	}
	if seconds > 0 {
		result += fmt.Sprintf("%d秒", seconds)
	}
	return result
}
