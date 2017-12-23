package util

import (
	"time"
	"strings"
	"fmt"
)

// DATE-TIME BEGIN
func NowInS() int64 {
	return time.Now().Unix()
}

func NowInNs() int64 {
	return time.Now().UnixNano()
}

func NowInMs() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func IsSameDay(date1, date2 time.Time) bool {
	y1, m1, d1 := date1.Date()
	y2, m2, d2 := date2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

func DateTimeToString(dateTime time.Time) (dateTimeStr string) {
	return dateTime.Format("2006-01-02 15:04:05")
}

func DateToString(dateTime time.Time) (dateStr string) {
	return dateTime.Format("2006-01-02")
}

func DateTimeToStringMs(dateTime time.Time) (dateTimeStr string) {
	return dateTime.Format("2006-01-02 15:04:05.999")
}

func ToDateTime(data interface{}) (res time.Time, err error) {
	switch data.(type) {
	case []byte:
		res, err = time.ParseInLocation("2006-01-02 15:04:05", strings.TrimSpace(string(data.([]byte))), time.Local)
	case string:
		res, err = time.ParseInLocation("2006-01-02 15:04:05", strings.TrimSpace(data.(string)), time.Local)
	default:
		res, err = time.ParseInLocation("2006-01-02 15:04:05", fmt.Sprintf("%v", data), time.Local)
	}
	return
}

func ToDate(data interface{}) (res time.Time, err error) {
	switch data.(type) {
	case []byte:
		res, err = time.ParseInLocation("2006-01-02", strings.TrimSpace(string(data.([]byte))), time.Local)
	case string:
		res, err = time.ParseInLocation("2006-01-02", strings.TrimSpace(data.(string)), time.Local)
	default:
		res, err = time.ParseInLocation("2006-01-02", fmt.Sprintf("%v", data), time.Local)
	}
	return
}

func MorningDateTime(now time.Time) time.Time {
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
}

func MiddayDateTime(now time.Time) time.Time {
	return time.Date(now.Year(), now.Month(), now.Day(), 12, 0, 0, 0, now.Location())
}

func NightDateTime(now time.Time) time.Time {
	return time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())
}

func GetYearAgoTime() int64 {
	romoteDay := time.Now().AddDate(-1, 0, 0)
	return romoteDay.UnixNano() / int64(time.Millisecond)
}
// DATE-TIME END