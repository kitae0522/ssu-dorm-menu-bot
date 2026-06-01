package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"testing"
)

func TestSendMessagePostsExpectedPayload(t *testing.T) {
	var gotPath string
	var gotChatID string
	var gotText string
	var gotPreview string

	httpClient := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		gotPath = r.URL.Path
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s, want POST", r.Method)
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("ReadAll: %v", err)
		}
		form, err := url.ParseQuery(string(body))
		if err != nil {
			t.Fatalf("ParseQuery: %v", err)
		}
		gotChatID = form.Get("chat_id")
		gotText = form.Get("text")
		gotPreview = form.Get("disable_web_page_preview")
		response, _ := json.Marshal(map[string]any{"ok": true})
		return &http.Response{
			StatusCode: http.StatusOK,
			Status:     "200 OK",
			Body:       io.NopCloser(bytes.NewReader(response)),
			Header:     make(http.Header),
		}, nil
	})}

	client := Client{HTTPClient: httpClient, BaseURL: "https://telegram.test"}
	err := client.SendMessage(context.Background(), "TOKEN", "12345", "오늘의 식단")
	if err != nil {
		t.Fatalf("SendMessage returned error: %v", err)
	}
	if gotPath != "/botTOKEN/sendMessage" {
		t.Fatalf("path = %q, want /botTOKEN/sendMessage", gotPath)
	}
	if gotChatID != "12345" || gotText != "오늘의 식단" || gotPreview != "true" {
		t.Fatalf("payload chat=%q text=%q preview=%q", gotChatID, gotText, gotPreview)
	}
}

func TestSendMessageReturnsTelegramError(t *testing.T) {
	httpClient := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		response, _ := json.Marshal(map[string]any{
			"ok":          false,
			"description": "chat not found",
		})
		return &http.Response{
			StatusCode: http.StatusBadRequest,
			Status:     "400 Bad Request",
			Body:       io.NopCloser(bytes.NewReader(response)),
			Header:     make(http.Header),
		}, nil
	})}

	client := Client{HTTPClient: httpClient, BaseURL: "https://telegram.test"}
	if err := client.SendMessage(context.Background(), "TOKEN", "bad", "hello"); err == nil {
		t.Fatal("expected error for Telegram ok=false response")
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}
