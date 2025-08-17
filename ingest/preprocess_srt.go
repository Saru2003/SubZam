package main

import (
    "regexp"
    "strings"
    "os"
    "log"
    "gopkg.in/yaml.v3"
    "github.com/vividvilla/metaphone"
)

type Config struct {
    Preprocess struct {
        BlockWindowSize int `yaml:"block_window_size"`
    } `yaml:"preprocess"`
}

var cfg Config

func init() {
    data, err := os.ReadFile("../config/config.yaml")
    if err != nil {
        log.Fatal("Could not read config:", err)
    }
    if err := yaml.Unmarshal(data, &cfg); err != nil {
        log.Fatal("Could not parse config:", err)
    }
}

type Chunk struct {
    Original string
    Raw      string 
    Cleaned  string 
    Phonetic string 
}
func PreprocessSRT(path string) ([]Chunk, error) {
    raw, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }

    htmlTagRegex := regexp.MustCompile(`</?[^>]+>`)
    text := htmlTagRegex.ReplaceAllString(string(raw), "")

    text = strings.ReplaceAll(text, "\r\n", "\n")

    blocks := regexp.MustCompile(`\n\s*\n`).Split(text, -1)

    var originalBlocks []string
    var cleanedBlocks []string

    for _, block := range blocks {
        lines := strings.Split(block, "\n")

        if len(lines) > 2 {
            lines = lines[2:]
        }

        originalJoined := strings.Join(lines, " ")
        originalJoined = strings.TrimSpace(originalJoined)

        cleanedJoined := strings.ToLower(originalJoined)
        cleanedJoined = regexp.MustCompile(`[^\w\s]`).ReplaceAllString(cleanedJoined, "")
        cleanedJoined = strings.Join(strings.Fields(cleanedJoined), " ")

        if cleanedJoined != "" {
            originalBlocks = append(originalBlocks, originalJoined)
            cleanedBlocks = append(cleanedBlocks, cleanedJoined)
        }
    }

    // Apply sliding window
    var chunks []Chunk
    win := cfg.Preprocess.BlockWindowSize
    for i := 0; i < len(cleanedBlocks); i++ {
        end := i + win
        if end > len(cleanedBlocks) {
            end = len(cleanedBlocks)
        }
        originalText := strings.Join(blocks[i:end], " ")
        cleanedText := strings.Join(cleanedBlocks[i:end], " ")
        phoneticText, _ := metaphone.DoubleMetaphone(cleanedText)

        chunks = append(chunks, Chunk{
            Original: originalText,
            Raw:      cleanedText,
            Cleaned:  cleanedText,
            Phonetic: phoneticText,
        })
    }

    return chunks, nil
}

func CleanText(input string) string {
    htmlTagRegex := regexp.MustCompile(`</?[^>]+>`)
    cleaned := htmlTagRegex.ReplaceAllString(input, "")

    cleaned = strings.ToLower(cleaned)

    cleaned = regexp.MustCompile(`[^\w\s]`).ReplaceAllString(cleaned, "")

    cleaned = strings.Join(strings.Fields(cleaned), " ")

    return cleaned
}
