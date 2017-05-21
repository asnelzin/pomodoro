package models

type Event struct {
	Type string `json:"type"`
	Object NewMessage `json:"object"`
	GroupID int `json:"group_id"`
	SecretKey string `json:"secret"`
}

type NewMessage struct {
	ID int `json:"id"`
	UserID int `json:"user_id"`
	FromID int `json:"from_id"`
	Date int `json:"date"`

	ReadState int `json:"read_state"`
	Out int `json:"out"`

	Title string `json:"title"`
	Body string `json:"body"`

	GeoData interface{} `json:"geo"`
	Attachments interface{} `json:"attachments"`
	FwdMessages interface{} `json:"fwd_messages"`

	Emoji int `json:"emoji"`
	Important int `json:"important"`
	Deleted int `json:"deleted"`
}