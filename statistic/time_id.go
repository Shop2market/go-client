package statistic

import (
	"fmt"
	"time"
)

func DailyTimeId(date time.Time) string {
	return fmt.Sprintf("%d%02d%02d", date.Year(), date.Month(), date.Day())
}
