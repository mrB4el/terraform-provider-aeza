package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/scinfra-pro/terraform-provider-aeza/internal/catalog"
	"github.com/scinfra-pro/terraform-provider-aeza/internal/client"
)

func main() {
	outputPath := flag.String("output", "docs/product-catalog.md", "path to the generated Markdown catalog")
	baseURL := flag.String("base-url", envOrDefault("AEZA_BASE_URL", "https://my.aeza.net/api"), "Aeza API base URL")
	flag.Parse()

	apiKey := os.Getenv("AEZA_API_KEY")
	if apiKey == "" {
		fatal("AEZA_API_KEY is required")
	}
	// The shared client has verbose request/response logging intended for
	// provider debugging. The catalog job must not copy API responses into CI
	// logs; the generated Markdown is the only intended output.
	log.SetOutput(io.Discard)

	apiClient, err := client.NewClient(*baseURL, apiKey)
	if err != nil {
		fatal("create Aeza API client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	products, err := apiClient.ListProducts(ctx)
	if err != nil {
		fatal("list Aeza products: %v", err)
	}
	if len(products) == 0 {
		fatal("Aeza API returned an empty product catalog; refusing to overwrite documentation")
	}

	if err := writeAtomic(*outputPath, []byte(catalog.RenderProductsMarkdown(products))); err != nil {
		fatal("write product catalog: %v", err)
	}
	fmt.Printf("updated %s with %d products\n", *outputPath, len(products))
}

func writeAtomic(path string, contents []byte) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	temporary, err := os.CreateTemp(dir, ".product-catalog-*.tmp")
	if err != nil {
		return err
	}
	temporaryPath := temporary.Name()
	defer os.Remove(temporaryPath)

	if _, err := temporary.Write(contents); err != nil {
		temporary.Close()
		return err
	}
	if err := temporary.Close(); err != nil {
		return err
	}
	return os.Rename(temporaryPath, path)
}

func envOrDefault(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}

func fatal(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
