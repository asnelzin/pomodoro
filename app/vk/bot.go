package vk

import (
	"github.com/asnelzin/pomodoro/app/models"
	"fmt"
	"context"
	"time"
	"net/http"
	"io/ioutil"
	"github.com/asnelzin/pomodoro/app/pomodoro"
	"log"
)

const (
	vkMessageSendURL = "https://api.vk.com/method/messages.send?"
)

const (
	POMODORO_DONE = "Well done!"
)


type Bot struct {
	APIToken string
	pomodoros map[int]*pomodoro.Job
}

func NewBot(token string) *Bot {
	return &Bot{
		APIToken: token,
		pomodoros: make(map[int]*pomodoro.Job),
	}
}

func (b Bot) HandleMessage(message *models.NewMessage) error {
	userID := message.UserID

	if val, ok := b.pomodoros[userID]; ok {
		val.Cancel()
	}

	ctx, cancel := context.WithCancel(context.Background())
	b.pomodoros[userID] = &pomodoro.Job{
		Cancel: cancel,
	}

	go b.waitTimer(ctx, 10*time.Second, userID)

	return nil
}

func (b Bot) sendMessageToUser(userID int, message string) error {
	req, err := http.NewRequest("GET", vkMessageSendURL, nil)
	if err != nil {
		return err
	}

	q := req.URL.Query()
	q.Add("message", message)
	q.Add("user_id", string(userID))
	q.Add("access_token", b.APIToken)
	q.Add("v", "5.0")
	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("%s", body)
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
				log.Printf("[ERROR] can't send message to user, %v", err)
			}
			return
		}
	}
}