package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func main() {
	host := flag.String("host", "localhost", "Set the host")
	port := flag.Uint("port", 8080, "The port to listen to")

	flag.Parse()

	addr := fmt.Sprintf("%s:%d", *host, *port)

	fmt.Printf("Starting to listen on %s...", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
