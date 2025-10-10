package main

import _ "github.com/lib/pq"

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/ValeriiaGrebneva/BlogAggregator/internal/config"
	"github.com/ValeriiaGrebneva/BlogAggregator/internal/database"
)

func main() {
	configStruct, err := config.Read()
	if err != nil {
		fmt.Println(err)
	}

	db, err := sql.Open("postgres", configStruct.Connection)
	dbQueries := database.New(db)

	stateConfig := state{
		cfg: &configStruct,
		db:  dbQueries,
	}

	commandsStruct := commands{
		mapCommands: make(map[string]func(*state, command) error),
	}

	commandsStruct.register("login", handlerLogin)
	commandsStruct.register("register", handlerRegister)
	commandsStruct.register("reset", handlerReset)
	commandsStruct.register("users", handlerUsers)
	commandsStruct.register("agg", handlerAggregator)
	commandsStruct.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	commandsStruct.register("feeds", handlerListFeeds)
	commandsStruct.register("follow", middlewareLoggedIn(handlerFollow))
	commandsStruct.register("following", middlewareLoggedIn(handlerFollowing))
	commandsStruct.register("unfollow", middlewareLoggedIn(handlerUnfollow))

	var args = os.Args
	if len(args) < 2 {
		fmt.Printf("Supposed to have at least 2 arguments, not %d arguments\n", len(args))
		os.Exit(1)
	}

	err = commandsStruct.run(&stateConfig, command{args[1], args[2:]})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return
}
