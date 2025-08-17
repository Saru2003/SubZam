package main

import (
	"hash/fnv"
	"strings"
	"github.com/vividvilla/metaphone"

)

// GenerateClosestHash = compute 64-bit SimHash for fuzzy matching
func GenerateClosestHash(text string) uint64 {
	normalized := strings.ToLower(strings.Join(strings.Fields(text), " "))
	words := strings.Fields(normalized)

	var bitCounts [64]int

	for _, word := range words {
		h := fnv.New64a()
		h.Write([]byte(word))
		hashVal := h.Sum64()

		for i := 0; i < 64; i++ {
			if (hashVal>>i)&1 == 1 {
				bitCounts[i] += 1
			} else {
				bitCounts[i] -= 1
			}
		}
	}

	//final SimHash
	var simhash uint64
	for i := 0; i < 64; i++ {
		if bitCounts[i] > 0 {
			simhash |= (1 << i)
		}
	}

	return simhash
}

// GeneratePhoneticHash using Double Metaphone
func GeneratePhoneticHash(text string) string {
	words := strings.Fields(strings.ToLower(text))
	if len(words) == 0 {
		return ""
	}

	var codes []string
	for _, w := range words {
		p, _ := metaphone.DoubleMetaphone(w)
		if p != "" {
			codes = append(codes, p)
		}
	}
	return strings.Join(codes, "-")
}