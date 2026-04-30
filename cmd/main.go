package main

import (
	"context"
	"log"
	"path/filepath"

	"github.com/bjornnorgaard/textbench"
)

var (
	moby = filepath.Join(".", "data", "moby-dick.txt")
)

func main() {
	ctx := context.Background()
	evaluate, err := textbench.Evaluate(ctx, moby)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("evaluate: %f", evaluate)
}
