package main

import (
	"context"
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

	// Generate a unique UUID for each feed
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
