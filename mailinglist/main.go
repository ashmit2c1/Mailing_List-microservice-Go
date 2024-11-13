package main

import (
	"log"
	"mailinglist/db"
	"mailinglist/server"
	"sync"
)

func main() {
	db, err := db.InitDB("subscribers.db")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		server.StartJSONServer(db)
	}()

	go func() {
		defer wg.Done()
		server.StartGRPCServer(db)
	}()

	wg.Wait()
}
