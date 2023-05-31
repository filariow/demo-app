package main

import (
	"context"
	_ "embed"
	"eshop-catalog/pkg/config"
	"log"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

//go:embed init-db.sql
var initSql string

//go:embed seed-db.sql
var seedSql string

var checkSeedSql = `SELECT CASE WHEN EXISTS (SELECT * FROM products p LIMIT 1) THEN 1 ELSE 0 END`

func main() {
	ctx := context.Background()

	c := config.NewConfigFromServiceBinding()
	ctxDb, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()
	conn, err := pgx.Connect(ctxDb, c.Postgres.ConnectionString())
	if err != nil {
		log.Fatalf("unable to connect to database: %v", err)
	}
	defer conn.Close(ctx)

	if err := runScript(ctx, conn, initSql); err != nil {
		log.Fatalf("unable to create tables: %s", err)
	}

	r := conn.QueryRow(ctx, checkSeedSql)
	var a int
	if err := r.Scan(&a); err != nil {
		log.Fatalf("unable to scan checkSeed result: %s", err)
	}

	switch a {
	case 0:
		if err := runScript(ctx, conn, seedSql); err != nil {
			log.Fatalf("unable to seed database: %s", err)
		}
		log.Println("database initialized")
	case 1:
		log.Println("database already seeded")
	}
}

func runScript(ctx context.Context, conn *pgx.Conn, sql string) error {
	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}

	qq := strings.Split(sql, ";")
	for _, q := range qq {
		if strings.Trim(q, " \n") == "" {
			continue
		}

		if _, err := tx.Exec(ctx, q); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}
