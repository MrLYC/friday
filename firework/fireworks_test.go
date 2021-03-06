package firework_test

import (
	"friday/firework"
	"testing"
	"time"
)

func TestCronFireworkUpdateTime(t *testing.T) {
	f := firework.NewCronFirework(
		"0,1 2-6/2 28-31 */1 * 2017",
		time.Date(2017, 2, 28, 2, 0, 0, 0, time.Local),
		&firework.Firework{},
	)

	timeValues := []string{
		"2017-02-28 02:01:00",
		"2017-02-28 04:00:00",
		"2017-02-28 04:01:00",
		"2017-02-28 06:00:00",
		"2017-02-28 06:01:00",

		"2017-03-28 02:00:00",
		"2017-03-28 02:01:00",
		"2017-03-28 04:00:00",
		"2017-03-28 04:01:00",
		"2017-03-28 06:00:00",
		"2017-03-28 06:01:00",

		"2017-03-31 02:00:00",
		"2017-03-31 02:01:00",
		"2017-03-31 04:00:00",
		"2017-03-31 04:01:00",
		"2017-03-31 06:00:00",
		"2017-03-31 06:01:00",

		"2017-04-28 02:00:00",
		"2017-04-28 02:01:00",
		"2017-04-28 04:00:00",
		"2017-04-28 04:01:00",
		"2017-04-28 06:00:00",
		"2017-04-28 06:01:00",

		"2017-04-30 02:00:00",
		"2017-04-30 02:01:00",
		"2017-04-30 04:00:00",
		"2017-04-30 04:01:00",
		"2017-04-30 06:00:00",
		"2017-04-30 06:01:00",

		"2017-05-28 02:00:00",
		"2017-05-28 02:01:00",
		"2017-05-28 04:00:00",
		"2017-05-28 04:01:00",
		"2017-05-28 06:00:00",
		"2017-05-28 06:01:00",

		"2017-05-31 02:00:00",
		"2017-05-31 02:01:00",
		"2017-05-31 04:00:00",
		"2017-05-31 04:01:00",
		"2017-05-31 06:00:00",
		"2017-05-31 06:01:00",

		"2017-06-28 02:00:00",
		"2017-06-28 02:01:00",
		"2017-06-28 04:00:00",
		"2017-06-28 04:01:00",
		"2017-06-28 06:00:00",
		"2017-06-28 06:01:00",
	}
	timeChangeTables := map[string]time.Time{
		"2017-03-28 06:01:00": time.Date(2017, 3, 30, 6, 1, 0, 0, time.Local),
		"2017-04-28 06:01:00": time.Date(2017, 4, 29, 6, 1, 0, 0, time.Local),
		"2017-05-28 06:01:00": time.Date(2017, 5, 30, 6, 1, 0, 0, time.Local),
	}
	for _, v := range timeValues {
		if !f.UpdateTime() || f.Time.Format("2006-01-02 15:04:05") != v {
			t.Errorf("update time failed: %s -> %v", v, f.Time)
		}
		t, ok := timeChangeTables[v]
		if ok {
			f.Time = t
		}
	}
}
