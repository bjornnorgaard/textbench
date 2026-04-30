package main

import (
	"log"
	"os"
	"path/filepath"
)

var (
	moby = filepath.Join(".", "data", "moby-dick.txt")
)

func main() {
	bytes, err := os.ReadFile(moby)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(string(bytes))
}
