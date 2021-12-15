package toolib

import "time"

const TimeFormatStr = "2006-01-02 15:04:05"

func TimeFormat(t time.Time) string {
	return t.Format(TimeFormatStr)
}
