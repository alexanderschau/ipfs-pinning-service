package main

import (
	"log"

	"github.com/alexanderschau/ipfs-pinning-service/api"
	"github.com/alexanderschau/ipfs-pinning-service/env"
)

func main() {
	env.Load()
	router := api.NewRouter()
	log.Fatal(router.Run(":3000"))
}
