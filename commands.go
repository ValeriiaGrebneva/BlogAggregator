package main

import (
	"fmt"

	"github.com/ValeriiaGrebneva/BlogAggregator/internal/config"
)

type state struct {
	statePointer *config.Config
}

type command struct {
	nameCommand string
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
	//remove check for the name of the command
	lengthCommands := len(cmd.argumentsCommand)
	if lengthCommands == 0 {
		return fmt.Errorf("No arguments in command: %s", cmd.nameCommand)
	}

	if cmd.nameCommand == "login" {
		if lengthCommands != 1 {
			return fmt.Errorf("Supposed to have 1 argument in login command, not %d", lengthCommands)
		}

		user := cmd.argumentsCommand[0]

		err := s.statePointer.SetUser(user)
		if err != nil {
			return err
		}

		fmt.Printf("The user '%s' has been set", user)
		return nil
	}

	return nil
}