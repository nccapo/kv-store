package main

import (
	"context"
	"fmt"
	"github.com/nccapo/kv-store/store"
	"time"
)

func main() {
	client := store.NewClient(&store.Options{
		Addr: "rand",
	})

	ctx := context.Background()

	client.Set(ctx, "mykey", "myvalue", 5*time.Second)

	// Get before expiration
	val, err := client.Get(context.Background(), "mykey")
	if err != nil {
		fmt.Println("Value:", val) // Output: "myvalue"
	}

	fmt.Println(val)
	// Wait for expiration
	time.Sleep(6 * time.Second)

	// Get after expiration
	if _, err := client.Get(context.Background(), "mykey"); err != nil {
		fmt.Println("Key expired") // Output: "Key expired"
	}

}
