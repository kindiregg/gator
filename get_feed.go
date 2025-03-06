package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"

	_ "github.com/lib/pq"
)

const XmlUrl = "https://www.wagslane.dev/index.xml"

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {

	client := &http.Client{}

	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("could not create new http request: %w", err)
	}

	req.Header.Set("User-Agent", "gator")

	resp, err := client.Do(req)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("could not get url: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("could not read xml body: %w", err)
	}

	var fetchedFeed RSSFeed

	if err := xml.Unmarshal(body, &fetchedFeed); err != nil {
		return &RSSFeed{}, fmt.Errorf("could not unmarshal xml body: %w", err)
	}

	return &fetchedFeed, nil
}
