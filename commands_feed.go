package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kindiregg/gator/internal/database"
)

func handlerAgg(s *state, cmd command) error {
	url := cmd.args[0]
	rFeed, err := fetchFeed(context.Background(), url)
	if err != nil {
		return err
	}

	fmt.Println(*rFeed)
	return nil
}

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.args) < 2 {
		return fmt.Errorf("incorrect format, use <name> <url>")
	}

	name := cmd.args[0]
	url := cmd.args[1]

	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUsername)
	if err != nil {
		return fmt.Errorf("could not get feed user '%s': %w", name, err)
	}

	id := uuid.New()
	now := time.Now().UTC()

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        id,
		CreatedAt: now,
		UpdatedAt: now,
		Name:      name,
		Url:       url,
		UserID:    user.ID,
	})
	if err != nil {
		return err
	}

	followID := uuid.New()
	_, err = s.db.CreateFollowFeed(context.Background(), database.CreateFollowFeedParams{
		ID:        followID,
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    user.ID,
		FeedID:    feed.ID,
	})

	if err != nil {
		return fmt.Errorf("could not create follow feed: %w", err)
	}

	fmt.Printf("Feed added: %v!\n", feed.Name)

	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeedsWithUsernames(context.Background())
	if err != nil {
		return fmt.Errorf("error getting feeds: %w", err)
	}
	fmt.Println("--- RSS Feeds ---")
	for _, feed := range feeds {
		fmt.Printf("Feed: %s\n", feed.Name)
		fmt.Printf("URL: %s\n", feed.Url)
		fmt.Printf("User: %s\n", feed.UserName)
		fmt.Println()
	}

	return nil
}

func handlerFollow(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("please include a url argument")
	}

	url := cmd.args[0]
	feed, err := s.db.GetFeedByURL(context.Background(), url)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("no feed found with the provided URL")
		}
		return fmt.Errorf("error querying feed: %v", err)
	}

	followFeedID := uuid.New()
	now := time.Now().UTC()
	user, err := s.db.GetUserByUsername(context.Background(), s.cfg.CurrentUsername)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("user not found")
		}
		return fmt.Errorf("error querying user: %v", err)
	}

	_, err = s.db.CreateFollowFeed(context.Background(), database.CreateFollowFeedParams{
		ID:        followFeedID,
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    user.ID,
		FeedID:    feed.ID,
	})

	if err != nil {
		return fmt.Errorf("failed to create follow record: %w", err)
	}

	fmt.Printf("You are now following the feed: %s\n", feed.Name)

	return nil
}
