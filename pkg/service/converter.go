package service

import "math"

func ConvertSecondsToHours(val int) int {
	rounded := math.Round(float64(val) / 60 / 60)
	return int(rounded)
}

func ConvertHoursToSeconds(val int) int {
	return val * 60 * 60
}
