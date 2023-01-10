package main

import (
	"context"
	"eshop-events-consumer/pkg/config"
	"eshop-events-consumer/pkg/events"
	"eshop-events-consumer/pkg/queue"
	"fmt"
	"log"
)

func main() {
	log.Println("Running orders-events-consumer")
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx := context.Background()
	cp := config.NewConfigFromServiceBinding()

	fmt.Printf("config: %+v\n", cp)

	c, err := events.NewDynamoStreams(ctx, cp.Aws, cp.DynamoStreams)
	if err != nil {
		return err
	}

	cq, err := queue.NewSQSManager(ctx, cp.Aws, cp.SQS)
	if err != nil {
		return err
	}

	cr, ce, err := c.ReadEvents(ctx)
	if err != nil {
		return err
	}

	log.Println("waiting for msg or errors")
	for {
		select {
		case msg := <-cr:
			log.Printf("sending message to sqs: '%s'", string(msg))
			if err := cq.SendMessage(ctx, msg); err != nil {
				log.Println(err)
			}
		case err := <-ce:
			log.Println(err)
		}
	}
}
