package repository

import (
	"fmt"
	"testing"
)

var client *RepositoryClient

func init() {
	client = NewRepositoryClient()
}

func TestClient_GetFirestoreFeedsBySrc(t *testing.T) {

	tests := []struct {
		name    string
		feedSrc string
		wantNum int
	}{
		{
			name:    "case_happy_feed_exist",
			feedSrc: "https://javascript.plainenglish.io/leetcode-algorithm-challenges-minimum-maximum-depth-of-binary-tree-d075a271fbda?source=rss----4b3a1ed4f11c---4",
			wantNum: 1,
		},
		{
			name:    "case_happy_feed_not_exist",
			feedSrc: "https://not_exist_src",
			wantNum: 0,
		},
	}

	for _, tt := range tests {
		if feeds, err := client.GetFirestoreFeedsBySrc(tt.feedSrc); err != nil {
			t.Error(err)
		} else if len(feeds) != tt.wantNum {
			t.Errorf("GetFeedsBySrc() failed, got = %d, want = %d\n", len(feeds), tt.wantNum)
		} else {
			fmt.Printf("GetFeedsBySrc() success, got = %d, want = %d\n", len(feeds), tt.wantNum)
		}
	}
}

func TestClient_CreateMySQLFeed(t *testing.T) {

	tests := []struct {
		name    string
		feed    Feed
		wantErr bool
	}{
		{
			name:    "case_happy_feed_exist",
			feed:    Feed{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		if err := client.CreateMySQLFeed(tt.feed); err != nil {
			t.Error(err)
		}

	}
}
