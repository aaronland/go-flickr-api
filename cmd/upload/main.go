package main

import (
	"context"
	"github.com/aaronland/go-flickr-api/application/upload"
	"log"
)

func main() {

	ctx := context.Background()

	app := &upload.UploadApplication{}
	err := app.Run(ctx)

	if err != nil {
		log.Fatalf("Failed to run upload application, %v", err)
	}
}
