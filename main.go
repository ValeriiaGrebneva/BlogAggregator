package main

import (
	"fmt"
	"github.com/ValeriiaGrebneva/BlogAggregator/internal/config"
)

func main() {
	configStruct, err := config.Read()
	if err != nil {
		fmt.Println(err)
	}

	err = configStruct.SetUser("Valeriia")
	if err != nil {
		fmt.Println(err)
	}

	configStruct, err = config.Read()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(configStruct)

	return
}
