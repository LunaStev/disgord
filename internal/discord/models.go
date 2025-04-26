package discord

import "encoding/json"

type Payload struct {
	Op int             `json:"op"`
	T  string          `json:"t"`
	D  json.RawMessage `json:"d"`
}

type MessageCreateData struct {
	ID        string `json:"id"`
	ChannelID string `json:"channel_id"`
	Content   string `json:"content"`
	Author    struct {
		ID       string `json:"id"`
		Username string `json:"username"`
	} `json:"author"`
}

type ReadyPayload struct {
	User struct {
		ID string `json:"id"`
	} `json:"user"`
}

type InteractionCreateData struct {
	ID    string `json:"id"`
	Token string `json:"token"`
	Data  struct {
		Name string `json:"name"`
	} `json:"data"`
}
