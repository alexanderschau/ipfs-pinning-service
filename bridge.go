package pinning

import (
	"context"
	"time"

	"github.com/gorilla/mux"
	openapi "go.alxs.xyz/ipfs-pinning/go"
)

func NewPinsApiService(s Service) openapi.PinsApiServicer {
	return &PinsApiService{
		userFunctions: s,
	}
}

func NewRouter(routers ...openapi.Router) *mux.Router {
	return openapi.NewRouter(routers...)
}

func NewPinsApiController(s openapi.PinsApiServicer) openapi.Router {
	return openapi.NewPinsApiController(s)
}

// ErrorResponse - return this on failure
func ErrorResponse(statusCode int, details string) (openapi.ImplResponse, error) {
	return openapi.Response(statusCode, openapi.Failure{
		Error: openapi.FailureError{
			Details: details,
		},
	}), nil
}

type PinsGetInputData struct {
	ctx    context.Context
	cid    []string
	name   string
	match  openapi.TextMatchingStrategy
	status []openapi.Status
	before time.Time
	after  time.Time
	limit  int32
	meta   map[string]string
}

type PinsPostInputData struct {
	ctx context.Context
	pin openapi.Pin
}

type PinsRequestidDeleteInputData struct {
	ctx       context.Context
	requestid string
}

type PinsRequestidGetInputData struct {
	ctx       context.Context
	requestid string
}

type PinsRequestidPostInputData struct {
	ctx       context.Context
	requestid string
	pin       openapi.Pin
}

type Service struct {
	// PinsGet - List pin objects
	PinsGet func(PinsGetInputData) (openapi.ImplResponse, error)
	// PinsPost - Add pin object
	PinsPost func(PinsPostInputData) (openapi.ImplResponse, error)
	// PinsRequestidDelete - Remove pin object
	PinsRequestidDelete func(PinsRequestidDeleteInputData) (openapi.ImplResponse, error)
	// PinsRequestidGet - Get pin object
	PinsRequestidGet func(PinsRequestidGetInputData) (openapi.ImplResponse, error)
	// PinsRequestidPost - Replace pin object
	PinsRequestidPost func(PinsRequestidPostInputData) (openapi.ImplResponse, error)
}

type PinsApiService struct {
	userFunctions Service
}

func (s *PinsApiService) PinsGet(ctx context.Context, cid []string, name string, match openapi.TextMatchingStrategy, status []openapi.Status, before time.Time, after time.Time, limit int32, meta map[string]string) (openapi.ImplResponse, error) {
	return s.userFunctions.PinsGet(PinsGetInputData{
		ctx:    ctx,
		cid:    cid,
		name:   name,
		match:  match,
		status: status,
		before: before,
		after:  after,
		limit:  limit,
		meta:   meta,
	})
}

func (s *PinsApiService) PinsPost(ctx context.Context, pin openapi.Pin) (openapi.ImplResponse, error) {
	return s.userFunctions.PinsPost(PinsPostInputData{
		ctx: ctx,
		pin: pin,
	})
}

func (s *PinsApiService) PinsRequestidDelete(ctx context.Context, requestid string) (openapi.ImplResponse, error) {
	return s.userFunctions.PinsRequestidDelete(PinsRequestidDeleteInputData{
		ctx:       ctx,
		requestid: requestid,
	})
}

func (s *PinsApiService) PinsRequestidGet(ctx context.Context, requestid string) (openapi.ImplResponse, error) {
	return s.userFunctions.PinsRequestidGet(PinsRequestidGetInputData{
		ctx:       ctx,
		requestid: requestid,
	})
}

func (s *PinsApiService) PinsRequestidPost(ctx context.Context, requestid string, pin openapi.Pin) (openapi.ImplResponse, error) {
	return s.userFunctions.PinsRequestidPost(PinsRequestidPostInputData{
		ctx:       ctx,
		requestid: requestid,
		pin:       pin,
	})
}
