package discord

import (
	"context"
	"net/http"

	"nhooyr.io/websocket"
)

type MessageHandler func(msg MessageCreateData)

type SlashCommandHandler func(ctx context.Context, interaction InteractionCreateData)

type Client struct {
	Token                string
	GatewayURL           string
	HTTPClient           *http.Client
	WebSocket            *websocket.Conn
	BotID                string
	MessageHandlers      []MessageHandler
	SlashCommandHandlers map[string]SlashCommandHandler
}

func NewClient(token string) *Client {
	return &Client{
		Token:      token,
		GatewayURL: "wss://gateway.discord.gg/?v=10&encoding=json",
		HTTPClient: &http.Client{},
	}
}

func (c *Client) OnMessageCreate(handler MessageHandler) {
	c.MessageHandlers = append(c.MessageHandlers, handler)
}

func (c *Client) OnSlashCommand(name string, handler SlashCommandHandler) {
	if c.SlashCommandHandlers == nil {
		c.SlashCommandHandlers = make(map[string]SlashCommandHandler)
	}
	c.SlashCommandHandlers[name] = handler
}
