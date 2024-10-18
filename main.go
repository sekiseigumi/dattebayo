package main

import (
	"log"

	"github.com/sekiseigumi/dattebayo/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatalf("cmd.Execute() failed with %v", err)
	}
}
