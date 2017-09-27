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
	"bytes"
	"encoding/json"
	"time"
)

type ItemAuthor struct {
	Name string `json:"name"`
}

type FeedItem struct {
	Id            string     `json:"id"`
	Url           string     `json:"url"`
	Title         string     `json:"title"`
	DatePublished string     `json:"date_published"`
	DateModified  string     `json:"date_modified,omitempty"`
	Author        ItemAuthor `json:"author"`
	ContentHtml   string     `json:"content_html"`
	ContentText   string     `json:"content_text,omitempty"`
}

type JSONFeed struct {
	Version     string     `json:"version"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	HomePageURL string     `json:"home_page_url"`
	FeedURL     string     `json:"feed_url"`
	Author      ItemAuthor `json:"author"`
	Icon        string     `json:"icon"`
	Favicon     string     `json:"favicon"`
	Items       []FeedItem `json:"items"`
}

func NewJSONFeed(site *SiteConfig, posts []*LatestPost) (string, error) {

	items := make([]FeedItem, 0)
	for _, p := range posts {
		items = append(items, FeedItem{
			Id:            site.BaseURL + "/post/" + p.UUID,
			Url:           site.BaseURL + "/post/" + p.UUID,
			Title:         p.Slugline,
			DatePublished: p.DatePublished.Format(time.RFC3339),
			DateModified:  p.DateUpdated.Format(time.RFC3339),
			Author:        ItemAuthor{"The Author"},
			ContentHtml:   MarkdownToHtml(p.Text),
		})
	}

	feed := JSONFeed{
		Version:     "https://jsonfeed.org/version/1",
		Title:       site.Title,
		Description: site.Description,
		HomePageURL: site.BaseURL,
		FeedURL:     site.BaseURL + "/feeds/json",
		Author:      ItemAuthor{"The Author"},
		Icon:        site.BaseURL + "/images/apple-touch-icon-180.png",
		Favicon:     site.BaseURL + "/images/apple-touch-icon-60.png",
		Items:       items,
	}

	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")

	if err := enc.Encode(feed); err != nil {
		return "", err
	}

	return buf.String(), nil
}
