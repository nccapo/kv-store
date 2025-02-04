package main

import (
	"fmt"
	"log"

	"github.com/nccapo/kv-store/store"
)

func main() {
	kv := store.New()

	kv.Set("key", "bunebam")

	val, _ := kv.Get("key")
	fmt.Println(val)

	exist := kv.Exists("key")
	fmt.Println(exist)

	err := kv.SaveSnapshot()
	if err != nil {
		log.Fatalln(err)
	}
}
