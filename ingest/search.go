package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/redis/go-redis/v9"
)

var (
	rdb       *redis.Client
	searchCtx = context.Background()
)

func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
}

// SimHash search
func SearchClosest(input string, maxDistance int) bool {
	cleaned := CleanText(input)
	queryHash := GenerateClosestHash(cleaned)
	queryPhonetic := GeneratePhoneticHash(cleaned)

	fmt.Println("Query hash:", queryHash)
	fmt.Println("Query phonetic:", queryPhonetic)

	keys, err := rdb.Keys(searchCtx, "closest:*").Result()
	if err != nil {
		log.Fatal("Redis KEYS error:", err)
	}

	found := false
	for _, key := range keys {
		parts := strings.Split(key, ":")
		if len(parts) != 2 {
			continue
		}

		storedHash, err := strconv.ParseUint(parts[1], 10, 64)
		if err != nil {
			continue
		}

		vals, err := rdb.HGetAll(searchCtx, key).Result()
		if err != nil || len(vals) == 0 {
			continue
		}

		// SimHash check
		dist := HammingDistance(queryHash, storedHash)
		if dist <= maxDistance {
			fmt.Printf("\n[SimHash Match] (distance %d):\n", dist)
			fmt.Printf("Title: %s (%s)\n", vals["title"], vals["year"])
			fmt.Println("Text :", vals["raw"])
			found = true
		}

		if vals["phonetic"] == queryPhonetic && queryPhonetic != "" {
			fmt.Println("\n[Phonetic Match]:")
			fmt.Printf("Title: %s (%s)\n", vals["title"], vals["year"])
			fmt.Println("Text :", vals["raw"])
			found = true
		}
	}

	return found
}

// Phonetic search with edit distance
func SearchPhonetic(input string, maxDistance int) bool {
	cleaned := CleanText(input)
	queryPhonetic := GeneratePhoneticHash(cleaned)
	fmt.Println("Query phonetic:", queryPhonetic)

	keys, err := rdb.Keys(searchCtx, "phonetic:*").Result()
	if err != nil {
		log.Fatal("Redis KEYS error:", err)
	}

	found := false
	for _, key := range keys {
		storedPhonetic := strings.TrimPrefix(key, "phonetic:")
		dist := LevenshteinDistance(queryPhonetic, storedPhonetic)

		if dist <= maxDistance {
			vals, err := rdb.HGetAll(searchCtx, key).Result()
			if err != nil || len(vals) == 0 {
				continue
			}
			fmt.Printf("\n[Phonetic Match] (distance %d):\n", dist)
			fmt.Printf("Title: %s (%s)\n", vals["title"], vals["year"])
			fmt.Println("Text :", vals["raw"])
			found = true
		}
	}
	return found
}

// Embedding-based semantic search (Postgres)
func SearchEmbedding(input string, topK int) bool {
	queryEmbedding, err := GenerateEmbedding(input)
	if err != nil {
		log.Println("Embedding generation error:", err)
		return false
	}

	rows, err := pg.QueryContext(pgCtx, `
        SELECT title, year, raw,
               1 - (embedding <=> $1) AS similarity
        FROM subtitle_embeddings
        ORDER BY embedding <=> $1
        LIMIT $2;
    `, queryEmbedding, topK)
	if err != nil {
		log.Println("Postgres search error:", err)
		return false
	}
	defer rows.Close()

	found := false
	fmt.Println("\n[Embedding Matches]:")
	for rows.Next() {
		var title, year, raw string
		var similarity float64
		if err := rows.Scan(&title, &year, &raw, &similarity); err == nil {
			fmt.Printf("Title: %s (%s) [score=%.4f]\n", title, year, similarity)
			fmt.Println("Text :", raw)
			found = true
		}
	}
	return found
}


func SearchWithFallback(input string) {
	fmt.Println("Searching for:", input)

	fmt.Println("\n--- SimHash Search ---")
	if SearchClosest(input, 17) {
		return
	}

	fmt.Println("\n--- Phonetic Search ---")
	if SearchPhonetic(input, 3) {
		return
	}

	fmt.Println("\n--- Embedding Search (Fallback) ---")
	SearchEmbedding(input, 5)
}

func LevenshteinDistance(a, b string) int {
	m, n := len(a), len(b)
	if m == 0 {
		return n
	}
	if n == 0 {
		return m
	}

	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}

	for i := 0; i <= m; i++ {
		dp[i][0] = i
	}
	for j := 0; j <= n; j++ {
		dp[0][j] = j
	}

	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			cost := 0
			if a[i-1] != b[j-1] {
				cost = 1
			}
			dp[i][j] = min(
				dp[i-1][j]+1,   // deletion
				dp[i][j-1]+1,   // insertion
				dp[i-1][j-1]+cost, // substitution
			)
		}
	}
	return dp[m][n]
}

func min(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}
