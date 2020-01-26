package anp

import "time"

const baseWeek = 2611 - 1074
const offset = 3 * 24 * 60 * 60

func WeekToTime(week int) time.Time {
	return time.Unix(int64(week+baseWeek)*604800-offset, 0)
}

func TimeToWeek(t time.Time) int {
	i := (t.Unix() + offset) / 604800
	return int(i - baseWeek)
}
