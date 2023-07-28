package logs

import (
	"time"
)

// GetAxis 计算 x 坐标轴间隔，以及数值
func GetAxis(beginDate, endDate int64, kind string) (xAxis []int64) {
	var (
		tBegin   = time.Unix(beginDate, 0)
		tEnd     = time.Unix(endDate, 0)
		baseLine = tBegin
		endLine  = tEnd
		unit     = time.Hour
	)
	switch kind {
	case "second":
		unit = time.Second
	case "minute":
		baseLine = time.Date(tBegin.Year(), tBegin.Month(), tBegin.Day(), tBegin.Hour(), tBegin.Minute(), 0, 0, tBegin.Location())
		unit = time.Minute
	case "hour":
		baseLine = time.Date(tBegin.Year(), tBegin.Month(), tBegin.Day(), tBegin.Hour(), 0, 0, 0, tBegin.Location())
		unit = time.Hour
	case "day":
		baseLine = time.Date(tBegin.Year(), tBegin.Month(), tBegin.Day(), 0, 0, 0, 0, tBegin.Location())
		unit = time.Hour * 24
	case "week":
		baseLine = time.Date(tBegin.Year(), tBegin.Month(), tBegin.Day()-int(tBegin.Weekday())+1, 0, 0, 0, 0, tBegin.Location())
		endLine = time.Date(tEnd.Year(), tEnd.Month(), tEnd.Day()-int(tEnd.Weekday())+8, 0, 0, 0, 0, tEnd.Location())
		unit = time.Hour * 24 * 7
	case "month":
		baseLine = time.Date(tBegin.Year(), tBegin.Month(), 1, 0, 0, 0, 0, tBegin.Location())
		endLine = time.Date(tEnd.Year(), tEnd.Month(), 1, 0, 0, 0, 0, tEnd.Location())
	case "year":
		baseLine = time.Date(tBegin.Year(), 1, 1, 0, 0, 0, 0, tBegin.Location())
		unit = time.Hour * 24 * 365 // approximate year length
	default:
		baseLine = time.Date(tBegin.Year(), tBegin.Month(), tBegin.Day(), 0, 0, 0, 0, tBegin.Location())
		unit = time.Hour * 24
	}

	if kind == "month" {
		for baseLine.Before(endLine) || baseLine.Equal(endLine) {
			xAxis = append(xAxis, baseLine.Unix())
			baseLine = time.Date(baseLine.Year(), baseLine.Month()+1, 1, 0, 0, 0, 0, baseLine.Location())
		}
	} else {
		for tm := baseLine.Unix(); tm <= endLine.Unix(); tm += int64(unit.Seconds()) {
			xAxis = append(xAxis, tm)
		}
	}
	return
}
