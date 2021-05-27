package main

import (
	"net/http"

	pinning "go.alxs.xyz/ipfs-pinning"
)

func main() {
	http.ListenAndServe(":8082", http.HandlerFunc(pinning.Handler))
}
