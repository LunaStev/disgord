package discord

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"nhooyr.io/websocket"
)

type helloPayload struct {
	Op int `json:"op"`
	D  struct {
		HeartbeatInterval int `json:"heartbeat"`
	} `json:"d"`
}

func (c *Client) ConnectGateway(ctx context.Context) error {
	conn, _, err := websocket.Dial(ctx, c.GatewayURL, nil)
	if err != nil {
		return err
	}

	c.WebSocket = conn
	log.Println("Connected to gateway")

	_, data, err := c.WebSocket.Read(ctx)
	if err != nil {
		return err
	}

	var hello helloPayload
	if err := json.Unmarshal(data, &hello); err != nil {
		return err
	}

	log.Printf("Received hello event. Heartbeat interval: %d ms", hello.D.HeartbeatInterval)

	// A temporary measure
	interval := hello.D.HeartbeatInterval
	if interval == 0 {
		interval = 40000 // fallback to 40 seconds
	}

	c.StartHeartbeat(ctx, interval)

	err = c.SendIdentify(ctx)
	if err != nil {
		return err
	}

	c.ListenEvents(ctx)

	return nil
}

func (c *Client) StartHeartbeat(ctx context.Context, intervalMs int) {
	ticker := time.NewTicker(time.Duration(intervalMs) * time.Millisecond)
	go func() {
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				heartbeat := struct {
					Op int         `json:"op"`
					D  interface{} `json:"d"`
				}{
					Op: 1,
					D:  nil,
				}
				data, _ := json.Marshal(heartbeat)
				err := c.WebSocket.Write(ctx, websocket.MessageText, data)
				if err != nil {
					log.Println("Failed to send heartbeat:", err)
					return
				}
				log.Println("Sent heartbeat")
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (c *Client) SendIdentify(ctx context.Context) error {
	identify := struct {
		Op int `json:"op"`
		D  struct {
			Token      string `json:"token"`
			Intents    int    `json:"intents"`
			Properties struct {
				OS      string `json:"$os"`
				Browser string `json:"$browser"`
				Device  string `json:"$device"`
			} `json:"properties"`
		} `json:"d"`
	}{
		Op: 2,
		D: struct {
			Token      string `json:"token"`
			Intents    int    `json:"intents"`
			Properties struct {
				OS      string `json:"$os"`
				Browser string `json:"$browser"`
				Device  string `json:"$device"`
			} `json:"properties"`
		}{
			Token:   c.Token,
			Intents: 513,
			Properties: struct {
				OS      string `json:"$os"`
				Browser string `json:"$browser"`
				Device  string `json:"$device"`
			}{
				OS:      "linux",
				Browser: "disgord",
				Device:  "disgord",
			},
		},
	}

	data, _ := json.Marshal(identify)
	err := c.WebSocket.Write(ctx, websocket.MessageText, data)
	if err != nil {
		return err
	}

	log.Println("Sent Identify payload")

	return nil
}

func (c *Client) ListenEvents(ctx context.Context) {
	go func() {
		for {
			_, data, err := c.WebSocket.Read(ctx)
			if err != nil {
				log.Println("Failed to read from WebSocket:", err)
				return
			}

			var payload Payload
			if err := json.Unmarshal(data, &payload); err != nil {
				log.Println("Failed to decode event payload:", err)
				continue
			}

			if payload.T == "MESSAGE_CREATE" {
				var msg MessageCreateData
				if err := json.Unmarshal(payload.D, &msg); err != nil {
					log.Println("Failed to decode message create payload:", err)
					continue
				}

				if msg.Author.ID == c.BotID {
					continue
				}

				log.Printf("Received message from %s: %s", msg.Author.Username, msg.Content)

				for _, handler := range c.MessageHandlers {
					handler(msg)
				}
				log.Printf("Payload raw: %s", string(payload.D))
			}

			if payload.T == "READY" {
				var ready ReadyPayload
				if err := json.Unmarshal(payload.D, &ready); err != nil {
					log.Println("Failed to decode ready payload:", err)
				} else {
					c.BotID = ready.User.ID
					log.Printf("Bot ID set to: %s", c.BotID)
				}
			}

			if payload.T == "INTERACTION_CREATE" {
				var interaction InteractionCreateData
				if err := json.Unmarshal(payload.D, &interaction); err != nil {
					log.Println("Failed to decode interaction payload:", err)
					continue
				}

				log.Println("Received slash command interaction")

				if handler, ok := c.SlashCommandHandlers[interaction.Data.Name]; ok {
					handler(ctx, interaction)
				} else {
					log.Println("No handler for slash command:", interaction.Data.Name)
				}
			}

			log.Printf("Received event type: %s", payload.T)
		}
	}()
}
