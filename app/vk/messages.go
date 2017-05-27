package vk

import "regexp"

const (
	POMODORO_DONE = "Well done!"
)

var minutesRegEx *regexp.Regexp = regexp.MustCompile("(\\d+)")