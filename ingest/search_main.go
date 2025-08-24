package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run . \"your search query\"")
		return
	}

	query := os.Args[1]
	SearchWithFallback(query)
}
