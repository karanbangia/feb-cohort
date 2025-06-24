package main

import (
	"context"
	"log"
	"os"

	"git.sr.ht/~jamesponddotco/bunnystorage-go"
)

func main() {

	readWriteKey, ok := os.LookupEnv("BUNNYNET_WRITE_API_KEY")
	if !ok {
		log.Fatal("missing env var: BUNNYNET_WRITE_API_KEY")
	}

	// Create new Config to be initialize a Client.
	cfg := &bunnystorage.Config{
		StorageZone: "maverick-14",
		Key:         readWriteKey,
		Endpoint:    bunnystorage.EndpointSingapore,
	}

	// Create a new Client instance with the given Config.
	client, err := bunnystorage.NewClient(cfg)
	if err != nil {
		log.Fatal(err)
	}

	// List files in the storage zone.
	files, _, err := client.List(context.Background(), "/")
	if err != nil {
		log.Fatal(err)
	}
	download, r, err := client.Download(context.Background(), files[0].ObjectName, "/tmp/"+files[0].ObjectName)
	if err != nil {
		return
	}
}
