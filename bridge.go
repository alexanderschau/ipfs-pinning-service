package pinning

import (
	"github.com/gorilla/mux"
	openapi "go.alxs.xyz/ipfs-pinning/go"
)

func NewPinsApiService() openapi.PinsApiServicer {
	return &openapi.PinsApiService{}
}

func NewRouter(routers ...openapi.Router) *mux.Router {
	return openapi.NewRouter(routers...)
}

func NewPinsApiController(s openapi.PinsApiServicer) openapi.Router {
	return openapi.NewPinsApiController(s)
}
