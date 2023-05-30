package main

import (
	"context"
	"eshop-catalog/pkg/config"
	"eshop-catalog/pkg/persistence"
	"eshop-catalog/pkg/queue"
	"eshop-catalog/pkg/rest"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

const (
	DefaultAddress = ":8080"
)

func main() {
	log.Printf("Starting server at '%s'", DefaultAddress)
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx := context.Background()
	cp := config.NewConfigFromServiceBinding()

	fmt.Printf("config: %+v\n", cp)

	sm, err := queue.NewSQSManager(ctx, cp.Aws, cp.SQS)
	if err != nil {
		return fmt.Errorf("error creating SQS Manager: %w", err)
	}

	r, err := persistence.NewPostgresRepo(ctx, cp.Postgres.ConnectionString())
	if err != nil {
		return fmt.Errorf("error connecting to psql: %w", err)
	}
	defer r.Close(ctx)

	s := rest.NewHttpServer(r)
	logHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		s.Mux.ServeHTTP(w, r)
	})

	go func() {
		cmsg, cerr, _ := sm.IncomingMessages(ctx)
		for {
			select {
			case msg := <-cmsg:
				processReceivedMessage(ctx, r, sm, msg)
			case err := <-cerr:
				log.Printf("error retrieving message: %s", err)
			}
		}
	}()

	return http.ListenAndServe(DefaultAddress, logHandler)
}

func processReceivedMessage(ctx context.Context, repo persistence.Repository, sm *queue.SQSManager, msg queue.OrderCreatedSQSMessage) {
	// TODO: this must be a transaction if more than one product in order is allowed
	for _, op := range msg.Value.Dynamodb.NewImage.OrderedProducts.Value {
		uos := op.Value.UnitsOrdered.Value
		uo, err := strconv.ParseInt(uos, 10, 64)
		if err != nil {
			log.Printf("error processing message: invalid quantity of units ordered: %s", uos)
		}
		if err := repo.AddOrderedUnits(ctx, msg.Value.Dynamodb.NewImage.ID.Value, op.Value.ID.Value, uo); err != nil {
			log.Printf("error processing message: error updating ordered units for product '%s': %s", op.Value.ID.Value, err)
		}
	}

	sm.CompleteMessage(ctx, msg)
	log.Printf("processed message: %+v", msg)
}
