package service

import (
	"content-parser/cmd/fetcher/repository"
	"fmt"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/mmcdole/gofeed"
	uuid "github.com/satori/go.uuid"
)

type Client struct {
	repoClient *repository.RepositoryClient
}

var client *Client

func init() {
	client = &Client{}
	client.repoClient = repository.NewRepositoryClient()
}

func NewClient() *Client {

	return client
}

func fetchRSS(rss *repository.RSS, ch chan<- []repository.Feed) {

	feeds := []repository.Feed{}

	fp := gofeed.NewParser()

	if rss.Src != "" {
		feed, _ := fp.ParseURL(rss.Src)

		for i := 0; i < len(feed.Items); i++ {
			f := repository.Feed{}
			f.Title = feed.Items[i].Title
			f.Id = genUUID()
			f.Guid = feed.Items[i].GUID
			f.Pub = feed.Items[i].PublishedParsed
			f.RSSId = rss.Id
			f.Src = feed.Items[i].Link
			if feed.Items[i].Author != nil {
				f.Author = feed.Items[i].Author.Name
			}
			content := ""
			if f.Content != "" {
				content = feed.Items[i].Content
			} else {
				content = feed.Items[i].Description
			}

			f.Content = content
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
			if err != nil {
				log.Fatal(err)
			}
			f.ContentSnippet = doc.Text()
			fmt.Printf("fetch Feed: %s\n", f.Title)
			feeds = append(feeds, f)

		}
	}
	ch <- feeds
}

func (c *Client) feedExist(f repository.Feed) bool {
	if _, err := c.repoClient.GetMySQLFeedTitleBySrc(f.Src); err != nil {
		return false
	} else {
		return true
	}
}

func (c *Client) createFeeds(ch chan []repository.Feed) {

	feeds := <-ch
	fmt.Printf("feed amount %d", len(feeds))

	for _, feed := range feeds {
		// fmt.Printf("handle createFeeds with feed: %+v", feed)
		if c.feedExist(feed) {
			fmt.Printf("Feed: %s already exist, skip\n", feed.Title)
			continue
		} else {
			fmt.Printf("store Feed: %s\n", feed.Title)
			c.createFeed(feed)
		}

	}
}

func (c *Client) createFeed(f repository.Feed) {
	c.repoClient.CreateMySQLFeed(f)
}

func (c *Client) UpdateRSSFeeds() {

	list := c.repoClient.GetFirestoreRSSList()
	// RSSNum := len(list)
	RSSNum := 1

	ch := make(chan []repository.Feed, RSSNum)

	for i := 0; i < RSSNum; i++ {
		go fetchRSS(list[i], ch)
	}

	for i := 0; i < RSSNum; i++ {
		c.createFeeds(ch)
	}
}

func genUUID() string {
	u4 := uuid.NewV4()
	return u4.String()
}
