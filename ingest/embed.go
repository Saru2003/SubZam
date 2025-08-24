package main

import (
	"context"
	"fmt"
	"log"
	"os"

	openai "github.com/sashabaranov/go-openai"
)

// OpenAI client
var openaiClient *openai.Client
var embedCtx = context.Background()

func init() {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENAI_API_KEY not set")
	}
	openaiClient = openai.NewClient(apiKey)
}

// creates an embedding vector for a given text
func GenerateEmbedding(text string) ([]float32, error) {
	resp, err := openaiClient.CreateEmbeddings(embedCtx, openai.EmbeddingRequest{
		Model: openai.AdaEmbeddingV2, 
		Input: text,
	})
	if err != nil {
		return nil, err
	}
	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("no embedding returned")
	}
	return resp.Data[0].Embedding, nil
}
