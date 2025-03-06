package main

import (
	"context"
	"fmt"

	"github.com/kindiregg/gator/internal/database"
)

func handlerAgg(s *state, cmd command) error {
	rFeed, err := fetchFeed(context.Background(), XmlUrl)
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

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		Name:   name,
		Url:    url,
		UserID: user.ID,
	})
	if err != nil {
		return err
	}

	fmt.Printf("Feed added: %v!\n", feed.Name)

	return nil
}

func handlerFeeds(s *state, cmd command) error {
	return nil
}
