package main

import (
	"context"
	"errors"
	"fmt"
	"goexample/pkg/applog"
	"io/fs"
	"log/slog"
	"os"
)

func main() {
	handler := applog.WithTraceIDHandler(applog.WithCustomHandler(
		slog.NewJSONHandler(os.Stdout, nil)))
	logger := slog.New(handler)
	slog.SetDefault(logger)

	foo := map[string]any{
		"aaa": "aaa",
		"bbb": false,
		"ccc": []string{"1", "2", "3"},
	}

	ctx := context.Background()
	ctx = applog.WithFields(ctx, foo)
	ctx = applog.WithTraceID(ctx)

	slog.InfoContext(ctx, "test message2")

	if _, err := os.Open("non-existing"); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			fmt.Println(err)                   // open non-existing: The system cannot find the file specified.
			fmt.Println(fs.ErrNotExist)        // file does not exist
			fmt.Println("file does not exist") // file does not exist
		} else {
			fmt.Println(err)
		}
	}
}
