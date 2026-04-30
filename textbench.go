package textbench

import (
	"bufio"
	"context"
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

	reader := bufio.NewReader(file)
	bufferSize := 1024
	buffer := make([]byte, 0, bufferSize)

	bytesRead, err := reader.Read(buffer[:cap(buffer)])
	buffer = buffer[:bytesRead]
	fmt.Println(string(buffer))

	if err != nil && err != io.EOF {
		slog.ErrorContext(ctx, "failed to read chunk", slog.Any("err", err))
	}

	return 0.0, err
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
