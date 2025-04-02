package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/nccapo/kv-store/store"
)

func main() {
	// Create a new client
	client := store.NewClient(&store.ClientOptions{
		Addr: ":8080",
	})

	ctx := context.Background()

	// Test Set operation
	fmt.Println("Testing Set operation...")
	err := client.Set(ctx, "test_key", "test_value", 5*time.Second)
	if err != nil {
		log.Fatalf("Failed to set key: %v", err)
	}
	fmt.Println("✓ Set operation successful")

	// Test Get operation (before expiration)
	fmt.Println("\nTesting Get operation (before expiration)...")
	value, err := client.Get(ctx, "test_key")
	if err != nil {
		log.Fatalf("Failed to get key: %v", err)
	}
	fmt.Printf("✓ Get operation successful. Value: %s\n", value)

	// Test Delete operation
	fmt.Println("\nTesting Delete operation...")
	err = client.Delete(ctx, "test_key")
	if err != nil {
		log.Fatalf("Failed to delete key: %v", err)
	}
	fmt.Println("✓ Delete operation successful")

	// Verify deletion
	fmt.Println("\nVerifying deletion...")
	_, err = client.Get(ctx, "test_key")
	if err == nil {
		log.Fatal("Key should have been deleted")
	}
	fmt.Println("✓ Deletion verified")

	// Test TTL expiration
	fmt.Println("\nTesting TTL expiration...")
	err = client.Set(ctx, "ttl_key", "ttl_value", 2*time.Second)
	if err != nil {
		log.Fatalf("Failed to set TTL key: %v", err)
	}
	fmt.Println("✓ Set TTL key successful")

	// Wait for expiration
	fmt.Println("Waiting for TTL expiration...")
	time.Sleep(3 * time.Second)

	// Try to get expired key
	_, err = client.Get(ctx, "ttl_key")
	if err == nil {
		log.Fatal("TTL key should have expired")
	}
	fmt.Println("✓ TTL expiration verified")

	fmt.Println("\nAll tests completed successfully!")
}
