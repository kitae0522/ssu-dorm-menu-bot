package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/kitae0522/ssu-dorm-menu-bot/internal/menu"
	"github.com/kitae0522/ssu-dorm-menu-bot/internal/message"
	"github.com/kitae0522/ssu-dorm-menu-bot/internal/ssudorm"
	"github.com/kitae0522/ssu-dorm-menu-bot/internal/telegram"
)

const defaultMenuPath = "data/menus.json"

func main() {
	if err := run(os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string, out io.Writer) error {
	flags := flag.NewFlagSet("crawl-weekly", flag.ContinueOnError)
	flags.SetOutput(io.Discard)

	output := flags.String("output", envOrDefault("MENU_JSON_PATH", defaultMenuPath), "path to write menu JSON")
	sourceURL := flags.String("url", ssudorm.DefaultSourceURL, "SSU dorm cafeteria menu URL")
	notify := flags.Bool("notify", false, "send crawl status to Telegram")
	timeout := flags.Duration("timeout", 30*time.Second, "HTTP request timeout")
	if err := flags.Parse(args); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	client := &http.Client{Timeout: *timeout}
	store, err := ssudorm.FetchAndParse(ctx, client, *sourceURL, time.Now())
	if err != nil {
		if *notify {
			if notifyErr := sendNotification(message.CrawlFailure(time.Now(), *sourceURL, err.Error()), *timeout); notifyErr != nil {
				return fmt.Errorf("%w; additionally failed to send crawl failure notification: %v", err, notifyErr)
			}
		}
		return err
	}
	if err := menu.Save(*output, store); err != nil {
		if *notify {
			if notifyErr := sendNotification(message.CrawlFailure(time.Now(), *sourceURL, err.Error()), *timeout); notifyErr != nil {
				return fmt.Errorf("%w; additionally failed to send crawl failure notification: %v", err, notifyErr)
			}
		}
		return err
	}
	if *notify {
		if err := sendNotification(message.CrawlSuccess(time.Now(), store, *output), *timeout); err != nil {
			return err
		}
	}

	fmt.Fprintf(out, "wrote %d menu days to %s\n", len(store.Days), *output)
	return nil
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func sendNotification(text string, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	client := telegram.Client{HTTPClient: &http.Client{Timeout: timeout}}
	return client.SendMessage(ctx, os.Getenv("TELEGRAM_BOT_TOKEN"), os.Getenv("TELEGRAM_CHAT_ID"), text)
}
