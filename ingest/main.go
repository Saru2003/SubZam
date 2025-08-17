package main
import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)
func main() {
	subsDir := "../subs"
	err := filepath.Walk(subsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
            return err
        }
        if filepath.Ext(path) == ".srt" {
            fmt.Println("Processing:", path)

            title, year := ParseFilename(path)

            chunks, err := PreprocessSRT(path)
            if err != nil {
                log.Println("Error preprocessing:", path, err)
                return nil
            }

            for _, chunk := range chunks {
                closestHash := GenerateClosestHash(chunk.Cleaned)
                phoneticHash := GeneratePhoneticHash(chunk.Cleaned)
                chunk.Phonetic = phoneticHash
                // embedding := GenerateEmbedding(chunk.Cleaned)

                StoreRedis(closestHash, phoneticHash, chunk, title, year)
                // StoreRedisPhonetic(phoneticHash, chunk, title, year)
                // StorePostgresEmbedding(embedding, chunk, title, year)
            }
        }
        return nil
    })

    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Ingestion complete")
}