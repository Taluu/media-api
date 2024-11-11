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
	host := flag.String("host", "localhost", "Set the host")
	port := flag.Uint("port", 8080, "The port to listen to")

	flag.Parse()

	addr := fmt.Sprintf("%s:%d", *host, *port)

	tagsRegistry := adapters.NewFakeTagRegistry()
	tagsService := services.NewTagService(tagsRegistry)
	http.Handle("GET /tags", ports.NewHttpTagsList(tagsService))
	http.Handle("POST /tags", ports.NewHttpTagCreate(tagsService))

	fmt.Printf("Starting to listen on %s...", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
