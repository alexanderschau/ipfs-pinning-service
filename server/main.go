package main

import (
	"log"

	"github.com/alexanderschau/ipfs-pinning-service/api"
)

func main() {
	router := api.NewRouter()
	log.Fatal(router.Run(":3000"))
}
