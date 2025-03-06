package main

import (
	"fmt"

	"github.com/kindiregg/gator/internal/config"
)

type state struct {
	*config.Config
}

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

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("username required for login command")
	}

	s.Config.CurrentUsername = cmd.args[0]

	// Save the updated config
	err := config.Write(*s.Config)
	if err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}
	fmt.Printf("User '%s' has been set\n", cmd.name)

	return nil
}
