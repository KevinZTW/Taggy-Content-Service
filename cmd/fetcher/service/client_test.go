package service

import (
	"content-parser/cmd/fetcher/repository"
	"fmt"
	"testing"
)

// var client *Client

// func init() {
// 	client = NewClient()
// }

func TestClient_FeedExists(t *testing.T) {
	f := repository.Feed{}
	res := client.feedExist(f)
	fmt.Printf("feed %t", res)
}

// func TestClient_FetchRSS(t *testing.T) {

// 	fp := gofeed.NewParser()

// 	feed, _ := fp.ParseURL("https://rss.lilydjwg.me/zhihuzhuanlan/mm-fe")
// 	fmt.Printf("%+v", feed.Items[0])
// 	if feed.Items[0].Content == "" {
// 		fmt.Printf("true")
// 	}

// 	doc, err := goquery.NewDocumentFromReader(strings.NewReader(feed.Items[0].Content))
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// doc.Find("script").Each(func(i int, el *goquery.Selection) {
// 	// 	el.Remove()
// 	// })

// 	fmt.Println(doc.Text())
// 	fmt.Println("end of ")

// }
