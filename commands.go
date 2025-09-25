package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ValeriiaGrebneva/BlogAggregator/internal/config"
	"github.com/ValeriiaGrebneva/BlogAggregator/internal/database"
	"github.com/google/uuid"
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

	_, err := s.db.GetUser(context.Background(), user)
	if err != nil {
		fmt.Printf("The user '%s' does not exist\n", user)
		os.Exit(1)
	}

	err = s.cfg.SetUser(user)
	if err != nil {
		return err
	}

	fmt.Printf("The user '%s' has been set\n", user)

	return nil
}

func handlerRegister(s *state, cmd command) error {
	lengthCommands := len(cmd.argumentsCommand)
	if lengthCommands != 1 {
		return fmt.Errorf("Supposed to have 1 argument (username) in register command, not %d arguments", lengthCommands)
	}

	user := cmd.argumentsCommand[0]

	_, err := s.db.GetUser(context.Background(), user)
	if err == nil {
		fmt.Printf("The user '%s' already exists\n", user)
		os.Exit(1)
	}

	userData, err := s.db.CreateUser(context.Background(), database.CreateUserParams{uuid.New(), time.Now(), time.Now(), user})
	if err != nil {
		return err
	}

	err = s.cfg.SetUser(user)
	if err != nil {
		return err
	}

	fmt.Printf("The user '%v' was created\n", userData)

	return nil
}

func handlerReset(s *state, cmd command) error {
	lengthCommands := len(cmd.argumentsCommand)
	if lengthCommands != 0 {
		return fmt.Errorf("Supposed to have 0 arguments in reset command, not %d arguments", lengthCommands)
	}

	err := s.db.ResetUsers(context.Background())
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	return nil
}

func handlerUsers(s *state, cmd command) error {
	lengthCommands := len(cmd.argumentsCommand)
	if lengthCommands != 0 {
		return fmt.Errorf("Supposed to have 0 arguments in users command, not %d arguments", lengthCommands)
	}

	usersSlice, err := s.db.GetUsers(context.Background())
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	configStruct, err := config.Read()
	if err != nil {
		fmt.Println(err)
	}

	for _, u := range usersSlice {
		if u.Name == configStruct.Username {
			fmt.Printf("* %s (current)\n", u.Name)
		} else {
			fmt.Printf("* %s\n", u.Name)
		}
	}

	return nil
}
