package main

import (
    "path/filepath"
    "regexp"
    "strings"
)

func ParseFilename(path string) (string, string) {
    base := filepath.Base(path)
    base = strings.TrimSuffix(base, filepath.Ext(base))

    cleaned := strings.ReplaceAll(base, ".", " ")
    cleaned = strings.ReplaceAll(cleaned, "_", " ")
    cleaned = strings.ReplaceAll(cleaned, "-", " ")

    yearRegex := regexp.MustCompile(`\b(19[0-9]{2}|20[0-9]{2})\b`)
    year := ""
    yearMatch := yearRegex.FindString(cleaned)
    if yearMatch != "" {
        year = yearMatch
        cleaned = strings.Replace(cleaned, yearMatch, "", 1)
    }

    commonTags := []string{
        "1080p", "720p", "480p", "x264", "x265", "bluray", "brrip", "yify", "yts", "webdl", "webrip", "hdrip",
        "dvdrip", "bdrip", "h264", "aac", "ac3", "sdh", "eng", "english", "subs", "sub",
    }
    for _, tag := range commonTags {
        re := regexp.MustCompile(`(?i)\b` + tag + `\b`)
        cleaned = re.ReplaceAllString(cleaned, "")
    }

    cleaned = strings.Join(strings.Fields(cleaned), " ")

    return strings.Title(cleaned), year
}
