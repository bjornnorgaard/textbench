package textbench

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"time"
)

func Evaluate(ctx context.Context, path string) (float32, error) {
	start := time.Now()
	defer func() {
		slog.InfoContext(ctx, "evaluation completed", slog.Duration("duration", time.Since(start)))
	}()

	file, err := os.Open(path)
	if err != nil {
		return 0.0, fmt.Errorf("failed to open file: %w", err)
	}
	defer func(file *os.File) {
		if err = file.Close(); err != nil {
			slog.ErrorContext(ctx, "failed to close file", slog.Any("err", err))
		}
	}(file)

	windowSize := 1024
	// Use a buffer size larger than the window to avoid constant memory shifts during Discard
	reader := bufio.NewReaderSize(file, windowSize*2)

	for {
		window, err := reader.Peek(windowSize)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			slog.ErrorContext(ctx, "failed to peek window", slog.Any("err", err))
			return 0.0, err
		}

		// Process the window (placeholder for actual evaluation logic)
		_ = window

		// Slide the window by 1 byte
		_, err = reader.Discard(1)
		if err != nil {
			slog.ErrorContext(ctx, "failed to slide window", slog.Any("err", err))
			return 0.0, err
		}
	}

	return 0.0, nil
}

func readChunk(reader *bufio.Reader, offset, length int) (int, []byte, error) {
	if offset > 0 {
		_, err := reader.Discard(offset)
		if err != nil {
			return 0, nil, err
		}
	}
	buffer := make([]byte, length)
	n, err := reader.Read(buffer)
	return n, buffer[:n], err
}
