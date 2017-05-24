package vk

import (
	"github.com/asnelzin/pomodoro/app/models"
	"fmt"
	"context"
	"time"
	"log"
	"net/http"
	"io/ioutil"
	"github.com/asnelzin/pomodoro/app/pomodoro"
)

const sendURL = "https://api.vk.com/method/messages.send?"

const (
	POMODORO_DONE = "Work is done!"
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
	userID := message.FromID

	if message.Body == "cancel" {
		log.Println("Canceling...")
		b.pomodoros[userID].Cancel()
		return nil
	}

	log.Printf("Starting timer for %d\n", string(userID))

	ctx, cancel := context.WithCancel(context.Background())
	b.pomodoros[userID] = &pomodoro.Job{
		Cancel: cancel,
	}

	go b.waitTimer(ctx, 10*time.Second, userID)

	return nil
}

func (b Bot) sendMessageToUser(userID int, message string) error {
	url := fmt.Sprintf(
		sendURL + "message=%s&user_id=%d&access_token=%s&v=5.0",
		message,
		userID,
		b.APIToken,
	)

	resp, err := http.Get(url)
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
			log.Println(b.sendMessageToUser(userID, POMODORO_DONE))
			return
		}
	}
}