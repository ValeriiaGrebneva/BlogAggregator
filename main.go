package main

import _ "github.com/lib/pq"

import (
	"fmt"
	"os"

	"github.com/ValeriiaGrebneva/BlogAggregator/internal/config"
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
		db: &dbQueries
	}

	var args = os.Args
	if len(args) < 2 {
		fmt.Printf("Supposed to have at least 2 arguments, not %d arguments\n", len(args))
		os.Exit(1)
	}

	var commandLogin = command{
		nameCommand:      args[1],
		argumentsCommand: args[2:],
	}

	commandsStruct := commands{
		mapCommands: make(map[string]func(*state, command) error),
	}

	commandsStruct.register(commandLogin.nameCommand, handlerLogin)
	//commandsStruct.register(commandRegister.nameCommand, handlerRegister)

	err = commandsStruct.run(&stateConfig, commandLogin)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	configStruct, err = config.Read()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(configStruct)

	return
}
