package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/kindiregg/gator/internal/config"
	"github.com/kindiregg/gator/internal/database"
)

type command struct {
	name string
	args []string
}

type commands struct {
	handlers map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(s *state, cmd command) error) {
	c.handlers[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	handler, ok := c.handlers[cmd.name]
	if !ok {
		return fmt.Errorf("unknown command: %s", cmd.name)
	}

	return handler(s, cmd)
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("username required for login command")
	}
	name := cmd.args[0]
	id := uuid.New()
	now := time.Now()

	user, err := s.db.GetUser(context.Background(), name)
	if err == nil {
		fmt.Println("A user with that name already exists. Exiting.")
		os.Exit(1)
	} else if err != sql.ErrNoRows {
		return fmt.Errorf("failed to check for existing user: %w", err)
	}

	user, err = s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        id,
		CreatedAt: now,
		UpdatedAt: now,
		Name:      name,
	})

	if err != nil {
		return fmt.Errorf("failed to create user %s: %w", name, err)
	}

	s.cfg.CurrentUsername = user.Name

	err = config.Write(*s.cfg)
	if err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Printf("User '%s' was created successfully\n", user.Name)

	return nil
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("username required for login command")
	}

	user, err := s.db.GetUser(context.Background(), cmd.args[0])
	if err != nil {
		fmt.Printf("User does not exist: '%s' \n ", user.Name)
		fmt.Println("Exiting")
		os.Exit(1)
	}

	s.cfg.CurrentUsername = cmd.args[0]

	err = config.Write(*s.cfg)
	if err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}
	fmt.Printf("User '%s' has been set\n", cmd.name)

	return nil
}

func handlerReset(s *state, cmd command) error {

	err := s.db.ResetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("could not delete all users: %w", err)
	}

	fmt.Println("sucessfully deleted all user entries")

	return nil
}

func handlerGetUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("could not get users: %w", err)
	}

	if len(users) == 0 {
		return fmt.Errorf("no users currently registered")
	}
	fmt.Println("Users:")
	for _, user := range users {

		if user == s.cfg.CurrentUsername {
			fmt.Printf("* %s (current)\n", user)
		} else {
			fmt.Printf("* %s\n", user)
		}

	}

	return nil
}

func handlerFollowing(s *state, cmd command) error {
	currentUser, err := s.db.GetUserByUsername(context.Background(), s.cfg.CurrentUsername)
	if err != nil {
		return fmt.Errorf("could not get user ID: %w", err)
	}

	following, err := s.db.GetFollowsForUser(context.Background(), currentUser.ID)
	if err != nil {
		return fmt.Errorf("could not get following list of user %s: %w", currentUser.Name, err)
	}
	// fmt.Println()
	// fmt.Printf("%s is following", currentUser.Name)
	for _, follow := range following {
		fmt.Println(follow.FeedName)
	}
	// fmt.Println()
	return nil
}
