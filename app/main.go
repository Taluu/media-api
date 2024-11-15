package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/Taluu/media-go/pkg/domain/media/adapters"
	"github.com/Taluu/media-go/pkg/domain/media/ports"
	"github.com/Taluu/media-go/pkg/domain/media/services"
)

func main() {
	// setup
	tagsRegistry := adapters.NewFakeTagRegistry()
	tagsService := services.NewTagService(tagsRegistry)

	mediasService := services.NewMediaService(
		adapters.NewFakeMediaRepository(),
		tagsRegistry,
		adapters.NewFakeUploader(),
	)

	// tags
	http.Handle("GET /tags", ports.NewHttpTagsList(tagsService))
	http.Handle("POST /tags", ports.NewHttpTagCreate(tagsService))

	// medias routes
	http.Handle("GET /medias/{tag}", ports.NewHttpMediaSeatch(mediasService))
	http.Handle("POST /medias", ports.NewHttpMediaCreate(mediasService))
	http.Handle("GET /viewer/{id}", ports.NewHttpMediaViewer(mediasService))

	// http server
	host := flag.String("host", "localhost", "Set the host")
	port := flag.Uint("port", 8080, "The port to listen to")

	flag.Parse()

	addr := fmt.Sprintf("%s:%d", *host, *port)

	fmt.Printf("Starting to listen on %s...", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
