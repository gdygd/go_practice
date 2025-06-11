package general

import "time"

func CheckElapsedTime(timer *time.Time, msDuration int) bool {
	currentTm := time.Now()

	if timer.IsZero() || timer == nil {
		*timer = time.Now()
	}

	//elapsed = timer - currentTm
	elapsed := currentTm.Sub(*timer)
	// nano sec to milli sec
	msTime := elapsed / 1000000

	if int(msTime) < msDuration {
		return false
	}

	*timer = time.Now()
	return true
}

func MinInt(a int, b int) int {
	if a > b {
		return b
	}
	return a
}

func MaxInt(a int, b int) int {
	if a > b {
		return a
	}
	return b
}
