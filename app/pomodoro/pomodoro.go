package pomodoro

import "context"

type Job struct {
	Cancel context.CancelFunc
}