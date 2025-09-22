package main

import (
	"log"
	"os"
)

func main() {
	args := os.Args[1:]
	if err := run(args); err != nil {
		log.Fatal(err)
	}
}

func run(args []string) error {
	panic("not implemented")
}
