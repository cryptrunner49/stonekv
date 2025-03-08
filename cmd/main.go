package main

import (
	"fmt"
	"log"

	"github.com/cryptrunner49/stonekvr/stone" // Adjust import path as needed
)

func main() {
	// Initialize the store (creates a new one if it doesn't exist)
	store, err := stone.NewStore("data.db")
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close()

	// Set a key/value pair
	err = store.Set([]byte("greeting"), []byte("Hello, StoneKVR! 🪨🚀"))
	if err != nil {
		log.Fatal(err)
	}

	// Retrieve a value
	value, err := store.Get([]byte("greeting"))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(value)) // Outputs: Hello, StoneKVR! 🪨🚀
}