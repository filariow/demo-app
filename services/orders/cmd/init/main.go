package main

import (
	"context"
	"eshop-orders/pkg/config"
	"fmt"
	"log"
)

func main() {
	c := config.NewConfigFromServiceBinding()
	if c.DynamoDB.Url == "" {
		c.DynamoDB.Url = fmt.Sprintf("https://dynamodb.%s.amazonaws.com", c.DynamoDB.Region)
	}

	// connect to database
	db, err := newDynamoDB(c.Aws, c.DynamoDB)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	// initialize database
	log.Printf("Seeding database")
	if err := db.Init(ctx); err != nil {
		log.Fatal(err)
	}
	log.Printf("Database seeded")
}
