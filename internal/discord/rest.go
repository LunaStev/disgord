package discord

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) SendMessage(ctx context.Context, channelID, content string) error {
	url := "https://discord.com/api/v10/channels/" + channelID + "/messages"

	body := map[string]string{
		"content": content,
	}
	jsonBody, _ := json.Marshal(body)

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bot "+c.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("failed to send message, status code: %d", resp.StatusCode)
	}

	return nil
}

func (c *Client) RegisterCommand(ctx context.Context, applicationID, guildID string) error {
	url := "https://discord.com/api/v10/applications/" + applicationID + "/guilds/" + guildID + "/commands"

	body := map[string]interface{}{
		"name":        "ping",
		"description": "Replies with pong",
		"type":        1, // 1 = CHAT_INPUT (slash command)
	}

	jsonBody, _ := json.Marshal(body)

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bot "+c.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("failed to register command, status code: %d", resp.StatusCode)
	}

	return nil
}

func (c *Client) SendInteractionResponse(ctx context.Context, interactionID, token, content string) error {
	url := "https://discord.com/api/v10/interactions/" + interactionID + "/" + token + "/callback"

	body := map[string]interface{}{
		"type": 4, // CHANNEL_MESSAGE_WITH_SOURCE
		"data": map[string]string{
			"content": content,
		},
	}

	jsonBody, _ := json.Marshal(body)

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bot "+c.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("failed to send interaction response, status code: %d", resp.StatusCode)
	}

	return nil
}
