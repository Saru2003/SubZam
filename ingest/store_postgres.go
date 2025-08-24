package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var pg *sql.DB
var pgCtx = context.Background()

func init() {
	connStr := "postgres://user:password@localhost:5432/subsdb?sslmode=disable"
	var err error
	pg, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Could not connect to Postgres:", err)
	}

	_, err = pg.ExecContext(pgCtx, `CREATE EXTENSION IF NOT EXISTS vector`)
	if err != nil {
		log.Fatal("Could not enable pgvector:", err)
	}

	_, err = pg.ExecContext(pgCtx, `
	CREATE TABLE IF NOT EXISTS subtitle_embeddings (
		id SERIAL PRIMARY KEY,
		title TEXT,
		year TEXT,
		chunk TEXT,
		raw TEXT,
		embedding VECTOR(1536) -- embedding size depends on model
	)`)
	if err != nil {
		log.Fatal("Could not create table:", err)
	}
}

// saves an embedding vector to Postgres
func StorePostgresEmbedding(embedding []float32, chunk Chunk, title, year string) error {
	_, err := pg.ExecContext(pgCtx,
		`INSERT INTO subtitle_embeddings (title, year, chunk, raw, embedding) VALUES ($1, $2, $3, $4, $5)`,
		title, year, chunk.Cleaned, chunk.Original, embedding,
	)
	return err
}
