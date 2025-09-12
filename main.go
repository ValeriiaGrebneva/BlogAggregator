package main

import (
	"fmt"
	"github.com/ValeriiaGrebneva/BlogAggregator/internal/config"
)

func main() {
	configStruct, err := config.Read()
	if err != nil {
		fmt.Println("Reading didn't work")
	}

	err = configStruct.SetUser("Valeriia")
	if err != nil {
		fmt.Println("User wasn't set")
	}

	configStruct, err = config.Read()
	if err != nil {
		fmt.Println("Reading again didn't work")
	}

	return
}
