package telegram

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const DefaultBaseURL = "https://api.telegram.org"

type Client struct {
	HTTPClient *http.Client
	BaseURL    string
}

func (c Client) SendMessage(ctx context.Context, token, chatID, text string) error {
	if token == "" {
		return fmt.Errorf("telegram bot token is required")
	}
	if chatID == "" {
		return fmt.Errorf("telegram chat id is required")
	}
	if text == "" {
		return fmt.Errorf("telegram message text is required")
	}

	httpClient := c.HTTPClient
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	baseURL := strings.TrimRight(c.BaseURL, "/")
	if baseURL == "" {
		baseURL = DefaultBaseURL
	}

	form := url.Values{}
	form.Set("chat_id", chatID)
	form.Set("text", text)
	form.Set("disable_web_page_preview", "true")

	endpoint := fmt.Sprintf("%s/bot%s/sendMessage", baseURL, token)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return err
	}

	var apiResp struct {
		OK          bool   `json:"ok"`
		Description string `json:"description"`
	}
	if len(body) > 0 {
		_ = json.Unmarshal(body, &apiResp)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 || !apiResp.OK {
		if apiResp.Description != "" {
			return fmt.Errorf("telegram sendMessage failed: %s", apiResp.Description)
		}
		return fmt.Errorf("telegram sendMessage failed: status %s", resp.Status)
	}
	return nil
}
