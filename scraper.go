package main

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/squashd/blog-aggregator/internal/database"
	"log"
	"sync"
	"time"
)

func startScraping(db *database.Queries, concurrency int, interval time.Duration) {
	log.Printf("Scraping on %v go routines every %s\n duration", concurrency, interval)
	ticker := time.NewTimer(interval)

	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedsToFetch(
			context.Background(), int32(concurrency),
		)
		if err != nil {
			log.Printf("failed to get feeds to fetch: %v", err)
			continue
		}

		wg := sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)
			go scrapeFeed(db, &wg, feed)
		}
		wg.Wait()
	}
}

func scrapeFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done()

	err := db.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{
		LastFetchedAt: sql.NullTime{
			Time: time.Now(),
		},
		ID: feed.ID,
	})
	if err != nil {
		log.Printf("failed to mark feed fetched: %v", err)
		return
	}

	rssFeed, err := fetchUrlToFeed(feed.Url)
	if err != nil {
		log.Printf("failed to fetch feed: %v\n", err)
		return
	}
	for _, item := range rssFeed.Channel.Item {
		log.Printf("Found post: %v\n", item.Title)
		description := sql.NullString{}
		if item.Description != "" {
			description.String = item.Description
			description.Valid = true
		}
		pubAt, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			log.Printf("failed to parse time: %v\n", err)
			continue
		}
		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Url:         item.Link,
			Description: description,
			PublishedAt: pubAt,
			FeedID:      feed.ID,
		})
		if err != nil {
			log.Printf("failed to create post: %v\n", err)
			continue
		}

		log.Printf("Feed %s collected, %v posts found\n", feed.Url, len(rssFeed.Channel.Item))
	}
}
