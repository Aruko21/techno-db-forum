package main

import (
	"fmt"
	apiserver "github.com/soulphazed/techno-db-forum/internal/app"
)

func main() {
	fmt.Println("Techno-db-forum by A.Kosenkov started!")
	if err := apiserver.Start(); err != nil {
		fmt.Println(err)
	}
}
