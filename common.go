package main

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

/*
func strToDate(date string) time.Time {
	if date == "" {
		return time.Now()
	}
	s := strings.Split(date, "-")
	year, _ := strconv.Atoi(s[0])
	month, _ := strconv.Atoi(s[1])
	day, _ := strconv.Atoi(s[2])
	return time.Date(year, time.Month(month), day, 12, 0, 0, 0, time.UTC)
}
*/

func getDefautFromDate() string {
	now := time.Now()

	var dt time.Time
	switch int(now.Weekday()) {
	case 1: // MON
		dt = now.AddDate(0, 0, -3)
	default:
		dt = now.AddDate(0, 0, -1)
	}
	return formatDate(dt)
}

func getDefautToDate() string {
	now := time.Now()
	var dt time.Time
	dt = now.AddDate(0, 0, -1)
	return formatDate(dt)
}

func formatDate(date time.Time) string {
	year, month, day := date.Date()
	return fmt.Sprintf("%04d-%02d-%02d", year, month, day)
}

// make date format into yyyy-mm-dd
func normalizeDate(date string) string {
	re := regexp.MustCompile(`^[0-9]{2,4}-?[0-9]{2}-?[0-3][0-9]$`)
	if !re.MatchString(date) {
		return ""
	}

	date = strings.ReplaceAll(date, "-", "")
	if len(date) == 6 {
		date = "20" + date
	}
	year := date[0:4]
	month := date[4:6]
	day := date[6:8]
	date = fmt.Sprintf("%s-%s-%s", year, month, day)
	return date
}
