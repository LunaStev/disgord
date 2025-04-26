package discord

import (
	"net/http"

	"nhooyr.io/websocket"
)

type Client struct {
	Token           string
	GatewayURL      string
	HTTPClient      *http.Client
	WebSocket       *websocket.Conn
	BotID           string
	MessageHandlers []MessageHandler
}

func NewClient(token string) *Client {
	return &Client{
		Token:      token,
		GatewayURL: "wss://gateway.discord.gg/?v=10&encoding=json",
		HTTPClient: &http.Client{},
	}
}
