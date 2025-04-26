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

	c.StartHeartbeat(ctx, hello.D.HeartbeatInterval)

	err = c.SendIdentify(ctx)
	if err != nil {
		return err
	}

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
		D struct {
			Token string `json:"token"`
		}
	}

	log.Println("Sent Identify payload")

	return nil
}