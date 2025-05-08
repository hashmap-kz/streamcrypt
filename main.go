package main

import (
	"log"

	"github.com/hashmap-kz/streaming-compress-encrypt/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatalf("error: %v", err)
	}
}
