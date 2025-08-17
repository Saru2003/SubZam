package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client
var searchCtx = context.Background()

func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
}

func SearchClosest(input string, maxDistance int) {
	cleaned := CleanText(input)

	// SimHash
	queryHash := GenerateClosestHash(cleaned)
	fmt.Println("Query hash:", queryHash)

	// Phonetic
	queryPhonetic := GeneratePhoneticHash(cleaned)
	fmt.Println("Query phonetic:", queryPhonetic)

	keys, err := rdb.Keys(searchCtx, "closest:*").Result()
	if err != nil {
		log.Fatal("Redis KEYS error:", err)
	}

	fmt.Println("Found", len(keys), "keys to check...")

	phoneticMatches := 0

	for _, key := range keys {
		parts := strings.Split(key, ":")
		if len(parts) != 2 {
			continue
		}

		storedHashStr := parts[1]
		storedHash, err := strconv.ParseUint(storedHashStr, 10, 64)
		if err != nil {
			continue
		}

		vals, err := rdb.HGetAll(searchCtx, key).Result()
		if err != nil || len(vals) == 0 {
			continue
		}

		//SimHash check
		dist := HammingDistance(queryHash, storedHash)
		if dist <= maxDistance {
			fmt.Printf("\n[SimHash Match] (distance %d):\n", dist)
			fmt.Println("Title:", vals["title"], "("+vals["year"]+")")
			fmt.Println("Text :", vals["raw"])
			fmt.Println("Phonetic in DB:", vals["phonetic"])
		}

		// Phonetic check
		if vals["phonetic"] == queryPhonetic && queryPhonetic != "" {
			fmt.Printf("\n[Phonetic Match]:\n")
			fmt.Println("Title:", vals["title"], "("+vals["year"]+")")
			fmt.Println("Text :", vals["raw"])
			fmt.Println("Phonetic in DB:", vals["phonetic"])
			phoneticMatches++
		}
	}

	if phoneticMatches == 0 {
		fmt.Println("\nNo phonetic matches found.")
	}
}

// LevenshteinDistance 
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

func SearchPhonetic(input string, maxDistance int) {
    cleaned := CleanText(input)
    queryPhonetic := GeneratePhoneticHash(cleaned)
    fmt.Println("Query phonetic:", queryPhonetic)

    keys, err := rdb.Keys(searchCtx, "phonetic:*").Result()
    if err != nil {
        log.Fatal("Redis KEYS error:", err)
    }

    fmt.Println("Found", len(keys), "phonetic keys to check...")

    matches := 0
    for _, key := range keys {
        storedPhonetic := strings.TrimPrefix(key, "phonetic:")
        dist := LevenshteinDistance(queryPhonetic, storedPhonetic)

        if dist <= maxDistance {
            vals, err := rdb.HGetAll(searchCtx, key).Result()
            if err != nil || len(vals) == 0 {
                continue
            }
            fmt.Printf("\n[Phonetic Match] (distance %d):\n", dist)
            fmt.Println("Title:", vals["title"], "("+vals["year"]+")")
            fmt.Println("Text :", vals["raw"])
            fmt.Println("Phonetic in DB:", storedPhonetic)
            matches++
        }
    }

    if matches == 0 {
        fmt.Println("\nNo phonetic matches found.")
    }
}
