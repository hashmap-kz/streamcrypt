package main

import (
	"log"

	"github.com/hashmap-kz/streamcrypt/v1/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatalf("error: %v", err)
	}
}
