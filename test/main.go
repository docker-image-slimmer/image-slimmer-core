package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/pnkcaht/image-slimmer-core/pkg/slimmer"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: go run . <image-reference>")
	}

	ref := os.Args[1]

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	engine := slimmer.New()

	result, err := engine.Slim(ctx, ref)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("==== IMAGE ====")
	fmt.Println(result.Image.Reference)
	fmt.Println("Digest:", result.Image.Digest)

	fmt.Println("\n==== PLAN ====")
	fmt.Println(result.Plan.Summary())

	fmt.Println("\n==== DETERMINISTIC ====")
	fmt.Println(result.Deterministic.Summary())
}
