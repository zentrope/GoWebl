//
// Copyright (c) 2017 Keith Irwin
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published
// by the Free Software Foundation, either version 3 of the License,
// or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

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

func NewRSSFeed(site *SiteConfig, posts []*LatestPost) (string, error) {

	items := make([]RSSItem, 0)
	for _, post := range posts {
		link := site.BaseURL + "/post/" + post.UUID
		items = append(items, RSSItem{
			Title:       post.Slugline,
			Link:        link,
			Guid:        RSSGuid{Guid: post.UUID, IsPermalink: false},
			Author:      "The Author",
			PubDate:     post.DatePublished.Format(time.RFC1123Z),
			Description: MarkdownToHtml(post.Text),
		})
	}

	channel := RSSChannel{
		Title:       site.Title,
		Link:        site.BaseURL,
		Description: site.Description,
		Items:       items,
	}

	feed := RSSFeed{Version: "2.0", Channel: channel}

	output, err := xml.MarshalIndent(feed, "", " ")
	if err != nil {
		return "", err
	}

	return xml.Header + string(output), nil
}
