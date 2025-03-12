package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"io"
	"log"
	"math/rand"
	"strings"
	"unicode"
)

func parseTemplate(template string) string {
	parts := strings.Split(template, "\n")

	for i, part := range parts {
		if len(part) > 0 {
			parts[i] = "RUN " + part
		}
	}

	return strings.Join(parts, "\n")
}

func randomString(length int) string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	output := make([]byte, length)

	for i := 0; i < len(output); i++ {
		output[i] = charset[rand.Intn(len(charset))]
	}

	return string(output)
}

func parseBashString(line string) []string {
	var tokens []string
	i := 0
	n := len(line)

	for i < n {
		for i < n && unicode.IsSpace(rune(line[i])) {
			i++
		}

		if i >= n {
			break
		}

		var token string
		if line[i] == '"' || line[i] == '\'' {
			quote := line[i]
			i++
			start := i

			for i < n && line[i] != quote {
				i++
			}

			token = line[start:i]

			if i < n && line[i] == quote {
				i++
			}
		} else {
			start := i
			for i < n && !unicode.IsSpace(rune(line[i])) {
				i++
			}
			token = line[start:i]
		}

		tokens = append(tokens, token)
	}

	return tokens
}

func handleLogStream(ctx context.Context, stream io.ReadCloser) string {
	runtime.EventsEmit(ctx, "openLog")
	defer stream.Close()
	scanner := bufio.NewScanner(stream)

	for scanner.Scan() {
		logLine := scanner.Text()
		var logJSON map[string]interface{}

		if err := json.Unmarshal([]byte(logLine), &logJSON); err == nil {
			if val, ok := logJSON["stream"]; ok {
				fmt.Println("docker output:", val)
				runtime.EventsEmit(ctx, "syslog", val)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Println("Error reading Docker build stream", err)
		return err.Error()
	}

	runtime.EventsEmit(ctx, "closeLog")
	return ""
}
