package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"google.golang.org/api/iterator"
)

type RSS struct {
	Description string    `firestore:"description,omitempty"`
	Id          string    `firestore:"id,omitempty"`
	ImgSrc      string    `firestore:"img,omitempty"`
	LastUpdate  time.Time `firestore:"omitempty"`
	Src         string    `firestore:"url,omitempty"`
	Title       string    `firestore:"title,omitempty"`
	Website     string    `firestore:"link,omitempty"`
}

func (c *RepositoryClient) GetFirestoreRSSList() []*RSS {
	ctx := context.Background()
	iter := c.FirestoreClient.Collection("RSSFetchList").Documents(ctx)
	list := []*RSS{}
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			fmt.Println(err)
		}
		rss := &RSS{}

		err = doc.DataTo(rss)
		if err != nil {
			log.Fatal(err)
		}

		list = append(list, rss)

	}
	return list
}

type Feed struct {
	RSSId          string `firestore:"RSSId,omitempty"`
	Id             string `firestore:"id,omitempty"`
	Title          string `firestore:"title,omitempty"`
	Author         string `firestore:"author,omitempty"`
	Content        string `firestore:"content,omitempty"`
	ContentSnippet string `firestore:"contentSnippet,omitempty"`
	Guid           string `firestore:"guid,omitempty"`
	Pub            *time.Time
	Src            string `firestore:"link,omitempty"`
}

func (c *RepositoryClient) GetFirestoreFeedsBySrc(src string) ([]Feed, error) {
	ctx := context.Background()
	feeds := []Feed{}
	iter := c.FirestoreClient.Collection("RSSItem").Where("link", "==", src).Documents(ctx)

	for {
		doc, err := iter.Next()

		if err == iterator.Done {
			return feeds, nil
		} else if err != nil {
			return nil, err
		}

		f := Feed{}
		if err = doc.DataTo(&f); err != nil {
			return nil, err
		}
		feeds = append(feeds, f)
	}
}

func (c *RepositoryClient) GetMySQLFeedById(id string) ([]Feed, error) {

	f := Feed{}
	fmt.Println("start query")

	stmtOut, err := c.MySQLDB.Prepare("SELECT * FROM Feed WHERE FeedId = ?")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtOut.Close()

	row := stmtOut.QueryRow(id)

	err = row.Scan(&f.Id, &f.RSSId, &f.Title, &f.Content, &f.ContentSnippet, &f.Guid, &f.Pub, &f.Src)
	if err == sql.ErrNoRows {
		fmt.Println("no rows!")
	} else if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n", f)

	return nil, nil
}

func (c *RepositoryClient) CreateMySQLFeed(f Feed) error {

	stmtIns, err := c.MySQLDB.Prepare("INSERT INTO Feed VALUES( ?, ?, ?, ?, ?, ?, ?, ? )")
	if err != nil {
		return err
	}
	defer stmtIns.Close()

	_, err = stmtIns.Exec(f.Id, f.RSSId, f.Title, f.Content, f.ContentSnippet, f.Guid, f.Pub, f.Src)
	if err != nil {
		return err
	}

	return nil
}

func (c *RepositoryClient) GetMySQLFeedTitleBySrc(src string) (string, error) {
	stmtOut, err := c.MySQLDB.Prepare("SELECT FeedTitle FROM Feed WHERE FeedLink = ?")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtOut.Close()

	row := stmtOut.QueryRow(src)
	var title string

	err = row.Scan(&title)
	if err == sql.ErrNoRows {
		return "", err
	} else if err != nil {
		return "", err
	}

	return title, nil

}

func (c *RepositoryClient) GetMySQLRSS() ([]Feed, error) {
	r := &RSS{}
	fmt.Println("start query")
	row := c.MySQLDB.QueryRow("SELECT * FROM RSS")
	fmt.Println("start scan")

	err := row.Scan(&r.Id, &r.Title, &r.Description, &r.Src, &r.LastUpdate, &r.ImgSrc)
	if err == sql.ErrNoRows {
		fmt.Println("no rows!")
	} else if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v", r)

	return nil, nil
}
