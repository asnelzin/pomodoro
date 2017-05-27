package vk

import (
	"fmt"
	"context"
	"time"
	"net/http"
	"log"
	"strconv"
)

type Pomodoro struct {
	Cancel context.CancelFunc
}

type Bot struct {
	APIToken string
	pomodoros map[int]*Pomodoro
}

func NewBot(token string) *Bot {
	return &Bot{
		APIToken: token,
		pomodoros: make(map[int]*Pomodoro),
	}
}

func (b Bot) HandleMessage(m *NewMessage) error {
	id := m.UserID

	if val, ok := b.pomodoros[id]; ok {
		val.Cancel()
	}

	ctx, cancel := context.WithCancel(context.Background())
	b.pomodoros[id] = &Pomodoro{
		Cancel: cancel,
	}

	duration := 2 * time.Second
	log.Printf("[INFO] starting new pomodoro for %v", duration)
	go b.waitTimer(ctx, duration, id)

	return nil
}

const MessageSendURL = "https://api.vk.com/method/messages.send?"
func (b Bot) sendMessageToUser(userID int, message string) error {
	req, err := http.NewRequest("GET", MessageSendURL, nil)
	if err != nil {
		return fmt.Errorf("could not create a request: %v", err)
	}

	q := req.URL.Query()
	q.Add("message", message)
	q.Add("user_id", strconv.Itoa(userID))
	q.Add("access_token", b.APIToken)
	q.Add("v", "5.0")
	req.URL.RawQuery = q.Encode()

	log.Printf("[INFO] calling an API with request: %v", req)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("could not call an API: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad response from a server: %s", resp.Status)
	}
	return nil
}

func (b Bot) waitTimer(ctx context.Context, d time.Duration, userID int) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(d):
			err := b.sendMessageToUser(userID, POMODORO_DONE)
			if err != nil {
				log.Printf("[ERROR] can't send message to user: %v", err)
			}
			return
		}
	}
}