package main

import (
	"encoding/xml"
	"html"
	"io"
	"net/http"
)

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
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "gator")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	readResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var feedRSS RSSFeed
	itemsRSS := []RSSItem{}
	feedRSS.Channel.Item = itemsRSS

	err = xml.Unmarshal(readResp, &feedRSS)
	if err != nil {
		return nil, err
	}

	feedRSS.Channel.Title = html.UnescapeStrings(feedRSS.Channel.Title)
	feedRSS.Channel.Description = html.UnescapeStrings(feedRSS.Channel.Description)
	for _, item := range feedRSS.Channel.Item {
		item.Title = html.UnescapeStrings(item.Title)
		item.Description = html.UnescapeStrings(item.Title)
	}

	return feedRSS, nil
}
