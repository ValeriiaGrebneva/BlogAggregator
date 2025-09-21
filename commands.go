package main

import (
	"fmt"

	"github.com/ValeriiaGrebneva/BlogAggregator/internal/config"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

type command struct {
	nameCommand      string
	argumentsCommand []string
}

type commands struct {
	mapCommands map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	commandFunc, exists := c.mapCommands[cmd.nameCommand]
	if !exists {
		return fmt.Errorf("Command '%s' doesn't exist", cmd.nameCommand)
	}

	err := commandFunc(s, cmd)
	if err != nil {
		return err
	}

	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.mapCommands[name] = f
}

func handlerLogin(s *state, cmd command) error {
	lengthCommands := len(cmd.argumentsCommand)
	if lengthCommands != 1 {
		return fmt.Errorf("Supposed to have 1 argument (username) in login command, not %d arguments", lengthCommands)
	}

	user := cmd.argumentsCommand[0]

	err := s.cfg.SetUser(user)
	if err != nil {
		return err
	}

	fmt.Printf("The user '%s' has been set\n", user)

	return nil
}

/* in progress
func handlerRegister(s *state, cmd command) error {
	lengthCommands := len(cmd.argumentsCommand)
	if lengthCommands != 1 {
		return fmt.Errorf("Supposed to have 1 argument (username) in login command, not %d arguments", lengthCommands)
	}

	user := cmd.argumentsCommand[0]

	err := s.cfg.SetUser(user)
	if err != nil {
		return err
	}

	fmt.Printf("The user '%s' has been set\n", user)

	return nil
}
*/
