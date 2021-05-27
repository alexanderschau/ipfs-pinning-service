package main

import (
	"log"
	"net/http"

	pinning "go.alxs.xyz/ipfs-pinning"
)

func main() {
	log.Printf("Server started")

	PinsApiService := pinning.NewPinsApiService()
	PinsApiController := pinning.NewPinsApiController(PinsApiService)

	router := pinning.NewRouter(PinsApiController)

	log.Fatal(http.ListenAndServe(":8080", router))
}
