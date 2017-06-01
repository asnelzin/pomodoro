package vk

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	POMODORO_DONE_RESPONSE = "Ты справился, молодец"
	CANCEL_RESPONSE        = "Отменил"
	DEFAULT_RESPONSE       = "Не понял"
	TOO_MUCH_RESPONSE      = "Слишком много для тебя"
)

var (
	minutesRegexp *regexp.Regexp = regexp.MustCompile("^(\\d+)$")
	cancelRegexp  *regexp.Regexp = regexp.MustCompile("^(отмена|стоп|cancel|stop)$")
)

type Pomodoro struct {
	Cancel context.CancelFunc
}

type Bot struct {
	APIToken  string
	pomodoros map[int]*Pomodoro
}

func NewBot(token string) *Bot {
	return &Bot{
		APIToken:  token,
		pomodoros: make(map[int]*Pomodoro),
	}
}

func (b Bot) HandleMessage(m *NewMessage) error {
	id := m.UserID
	text := strings.ToLower(m.Body)

	// First: try to find pomodoro duration in message
	match := minutesRegexp.FindStringSubmatch(text)

	switch {
	case len(match) > 0:
		if val, ok := b.pomodoros[id]; ok {
			val.Cancel()
		}

		ctx, cancel := context.WithCancel(context.Background())
		b.pomodoros[id] = &Pomodoro{
			Cancel: cancel,
		}

		duration, err := strconv.Atoi(match[0])
		if err != nil {
			return fmt.Errorf("could not parse string to int %s: %v", match, err)
		}

		if duration > 60 {
			go b.say(id, TOO_MUCH_RESPONSE)
		}

		log.Printf("[INFO] starting new pomodoro for %v", duration)
		go b.waitTimer(ctx, time.Duration(duration)*time.Second, id)
	case cancelRegexp.MatchString(text):
		if val, ok := b.pomodoros[id]; ok {
			val.Cancel()
			delete(b.pomodoros, id)
		}
		go b.say(id, CANCEL_RESPONSE)
	default:
		log.Printf("[INFO] could not find any pattern matching string %s", text)
		go b.say(id, DEFAULT_RESPONSE)
	}

	return nil
}

func (b Bot) say(id int, message string) {
	err := b.sendMessageToUser(id, message)
	if err != nil {
		log.Printf("[ERROR] can't send message to user: %v", err)
	}
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
	select {
	case <-ctx.Done():
		return
	case <-time.After(d):
		go b.say(userID, POMODORO_DONE_RESPONSE)
	}
}
