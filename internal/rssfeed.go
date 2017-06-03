// Copyright 2017 Keith Irwin. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package internal

import (
	"encoding/xml"
	"time"
)

type RSSGuid struct {
	XMLName     xml.Name `xml:"guid"`
	IsPermalink bool     `xml:"isPermaLink,attr"`
	Guid        string   `xml:",chardata"`
}

type RSSItem struct {
	XMLName     xml.Name `xml:"item"`
	Title       string   `xml:"title"`
	Link        string   `xml:"link"`
	Guid        RSSGuid  `xml:"guid"`
	PubDate     string   `xml:"pubDate"`
	Author      string   `xml:"author"`
	Description string   `xml:"description"`
}

type RSSChannel struct {
	Title       string    `xml:"title"`
	Link        string    `xml:"link"`
	Description string    `xml:"description"`
	Items       []RSSItem `xml:"items"`
}

type RSSFeed struct {
	XMLName xml.Name   `xml:"rss"`
	Version string     `xml:"version,attr"`
	Channel RSSChannel `xml:"channel"`
}

func NewRSSFeed(config WebConfig, posts []*LatestPost) (string, error) {

	items := make([]RSSItem, 0)
	for _, post := range posts {
		link := config.BaseURL + "/post/" + post.UUID
		items = append(items, RSSItem{
			Title:       post.Slugline,
			Link:        link,
			Guid:        RSSGuid{Guid: post.UUID, IsPermalink: false},
			Author:      post.Author,
			PubDate:     post.DateCreated.Format(time.RFC1123Z),
			Description: MarkdownToHtml(post.Text),
		})
	}

	channel := RSSChannel{
		Title:       config.Title,
		Link:        config.BaseURL,
		Description: "Last 40 bloops of " + config.Title,
		Items:       items,
	}

	feed := RSSFeed{Version: "2.0", Channel: channel}

	output, err := xml.MarshalIndent(feed, "", " ")
	if err != nil {
		return "", err
	}

	return xml.Header + string(output), nil
}
