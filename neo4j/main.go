package main

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func main() {
	ctx := context.Background()
	// URI examples: "neo4j://localhost", "neo4j+s://xxx.databases.neo4j.io"
	dbUri := "neo4j+s://7c89e1dc.databases.neo4j.io"
	dbUser := "neo4j"
	dbPassword := "iL9zhN9zCm_MPGM6XPUHBgLeS9DnB8l7zC81xlMw93Y"
	driver, err := neo4j.NewDriverWithContext(
		dbUri,
		neo4j.BasicAuth(dbUser, dbPassword, ""))
	defer driver.Close(ctx)

	err = driver.VerifyConnectivity(ctx)
	if err != nil {
		panic(err)
	}
}
