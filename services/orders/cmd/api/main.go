package main

import (
	"eshop-orders/pkg/config"
	"eshop-orders/pkg/persistence"
	"eshop-orders/pkg/rest"
	"fmt"
	"log"
	"net/http"
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
	c := config.NewConfigFromServiceBinding()
	if c.DynamoDB.Url == "" {
		c.DynamoDB.Url = fmt.Sprintf("https://dynamodb.%s.amazonaws.com", c.DynamoDB.Region)
	}

	fmt.Printf("config: %+v", c)

	r, err := persistence.NewDynamoDB(c.Aws, c.DynamoDB)
	if err != nil {
		return err
	}

	s := rest.NewHttpServer(r)

	logHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request %s:%s", r.Method, r.URL.Path)
		s.Mux.ServeHTTP(w, r)
	})
	return http.ListenAndServe(DefaultAddress, logHandler)
}
