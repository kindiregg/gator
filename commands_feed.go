package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kindiregg/gator/internal/database"
)

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 2 {
		return fmt.Errorf("incorrect format, use <name> <url>")
	}

	inputFeedName := cmd.args[0]
	url := cmd.args[1]

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      inputFeedName,
		Url:       url,
		UserID:    user.ID,
	})
	if err != nil {
		return err
	}

	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
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

func handlerFollow(s *state, cmd command, user database.User) error {
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

	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})

	if err != nil {
		return fmt.Errorf("failed to create follow record: %w", err)
	}

	fmt.Println("Feed follow created:")
	fmt.Printf("* User:          %s\n", user.Name)
	fmt.Printf("* Feed:          %s\n", feed.Name)
	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
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

	err = s.db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("could not delete feed follow: %w", err)
	}

	fmt.Printf("%s successfully unfollowed %s", user.Name, feed.Name)

	return nil
}
