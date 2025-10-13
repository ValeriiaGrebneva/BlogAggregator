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
		return fmt.Errorf("Supposed to have 1 argument (username) in Register command, not %d arguments", lengthCommands)
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
		return fmt.Errorf("Supposed to have 0 arguments in Reset command, not %d arguments", lengthCommands)
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
		return fmt.Errorf("Supposed to have 0 arguments in Users command, not %d arguments", lengthCommands)
	}

	usersSlice, err := s.db.GetUsers(context.Background())
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	for _, u := range usersSlice {
		if u.Name == s.cfg.Username {
			fmt.Printf("* %s (current)\n", u.Name)
		} else {
			fmt.Printf("* %s\n", u.Name)
		}
	}

	return nil
}

func handlerAggregator(s *state, cmd command) error {
	lengthCommands := len(cmd.argumentsCommand)
	if lengthCommands != 1 {
		return fmt.Errorf("Supposed to have 1 argument (time between requests) in Aggregator command, not %d arguments", lengthCommands)
	}

	timeReq := cmd.argumentsCommand[0]
	fmt.Printf("Collecting feeds every %s\n", timeReq)
	timeBetweenRequests, _ := time.ParseDuration(timeReq)

	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		err := scrapeFeeds(s)
		if err != nil {
			return err
		}
	}
	
	return nil
}

func scrapeFeeds(s *state) error {
	feedData, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}

	err = s.db.MarkFeedFetched(context.Background(), feedData.ID)
	if err != nil {
		return err
	}

	feedRSS, err := fetchFeed(context.Background(), feedData.Url)
	if err != nil {
		return err
	}

	for _, item := range feedRSS.Channel.Item {
		post, err := s.db.CreatePost(context.Background(), database.CreatePostParams{uuid.New(), time.Now(), time.Now(), item.Title, item.Link, item.Description, item.PubDate, feedData.ID})
		//continue here
	}

	return nil
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	lengthCommands := len(cmd.argumentsCommand)
	if lengthCommands != 2 {
		return fmt.Errorf("Supposed to have 2 arguments (FeedName, URL) in AddFeeds command, not %d arguments", lengthCommands)
	}

	name := cmd.argumentsCommand[0]
	url := cmd.argumentsCommand[1]

	feedData, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{uuid.New(), time.Now(), time.Now(), name, url, user.ID})
	if err != nil {
		return err
	}

	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{uuid.New(), time.Now(), time.Now(), user.ID, feedData.ID})
	if err != nil {
		return err
	}

	fmt.Println(feedData)
	return nil
}

func handlerListFeeds(s *state, cmd command) error {
	lengthCommands := len(cmd.argumentsCommand)
	if lengthCommands != 0 {
		return fmt.Errorf("Supposed to have 0 arguments in ListFeeds command, not %d arguments", lengthCommands)
	}

	feedsSlice, err := s.db.GetFeeds(context.Background())
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	for _, f := range feedsSlice {
		fmt.Printf("Name: %s\n", f.Name)
		fmt.Printf("URL: %s\n", f.Url)
		name, err := s.db.GetUserName(context.Background(), f.UserID)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("User: %s\n", name)
	}

	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	lengthCommands := len(cmd.argumentsCommand)
	if lengthCommands != 1 {
		return fmt.Errorf("Supposed to have 1 argument (URL) in Follow command, not %d arguments", lengthCommands)
	}

	url := cmd.argumentsCommand[0]
	feed, err := s.db.GetFeed(context.Background(), url)
	if err != nil {
		fmt.Printf("The feed '%s' doesn't exist\n", url)
		os.Exit(1)
	}

	feedFollowData, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{uuid.New(), time.Now(), time.Now(), user.ID, feed.ID})
	if err != nil {
		return err
	}

	fmt.Printf("Feed: %s\n", feedFollowData.FeedName)
	fmt.Printf("Name: %s\n", feedFollowData.UserName)
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	lengthCommands := len(cmd.argumentsCommand)
	if lengthCommands != 0 {
		return fmt.Errorf("Supposed to have 0 arguments in Following command, not %d arguments", lengthCommands)
	}

	feedFollowDataSlice, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}

	for _, f := range feedFollowDataSlice {
		fmt.Printf("* %s\n", f.FeedName)
	}

	return nil
}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.cfg.Username)
		if err != nil {
			fmt.Printf("The user '%s' doesn't exist\n", s.cfg.Username)
			os.Exit(1)
		}
		return handler(s, cmd, user)
	}
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	lengthCommands := len(cmd.argumentsCommand)
	if lengthCommands != 1 {
		return fmt.Errorf("Supposed to have 1 argument (URL) in Unfollow command, not %d arguments", lengthCommands)
	}

	url := cmd.argumentsCommand[0]
	feed, err := s.db.GetFeed(context.Background(), url)
	if err != nil {
		fmt.Printf("The feed '%s' doesn't exist\n", url)
		os.Exit(1)
	}

	err = s.db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{user.ID, feed.ID})
	if err != nil {
		return err
	}

	return nil
}
