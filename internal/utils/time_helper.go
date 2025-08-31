package utils

import (
	"fmt"
)

func formatDateTimeWithWeekday(date, startTime, endTime string) string {
	parsedTime, err := DateStringToTime(date)
	if err != nil {
		return fmt.Sprintf("%s %s~%s", date, startTime, endTime)
	}

	weekdays := []string{"日", "一", "二", "三", "四", "五", "六"}
	weekday := weekdays[parsedTime.Weekday()]

	return fmt.Sprintf("%s (%s) %s - %s", parsedTime.Format("2006/01/02"), weekday, startTime, endTime)
}
