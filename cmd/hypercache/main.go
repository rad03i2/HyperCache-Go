package main

import (
	"flag"
	"log"
	"time"

	"github.com/rad03i2/HyperCache-Go/pkg/cache"
	"github.com/rad03i2/HyperCache-Go/pkg/server"
)

const version = "1.0.0"

func main() {
	addr := flag.String("addr", "127.0.0.1:6379", "TCP address to listen on")
	shards := flag.Uint("shards", 64, "number of cache shards")
	cleanup := flag.Duration("cleanup", 5*time.Second, "expired-key cleanup interval")
	flag.Parse()

	c := cache.New(cache.Options{Shards: uint32(*shards), CleanupInterval: *cleanup})
	defer c.Close()
	log.Printf("HyperCache %s by Radwan Abdulhadi Ahmed / @rad03i2", version)
	if err := server.NewTCPServer(*addr, c).Start(); err != nil {
		log.Fatal(err)
	}
}
