package pinning

import (
	"context"
	"net/http"
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
func ErrorResponse(statusCode int, details string) Response {
	resp := openapi.Response(statusCode, openapi.Failure{
		Error: openapi.FailureError{
			Reason:  http.StatusText(statusCode),
			Details: details,
		},
	})
	return Response{resp, nil}
}

// PinCreationStatusResponse - return this, if a new pin was added
func PinCreationStatusResponse(pinStatus openapi.PinStatus) Response {
	resp := openapi.Response(http.StatusCreated, pinStatus)
	return Response{resp, nil}
}

// PinStatusResponse - return this on PinsRequestidGet
func PinStatusResponse(pinStatus openapi.PinStatus) Response {
	resp := openapi.Response(http.StatusOK, pinStatus)
	return Response{resp, nil}
}

// DeleteResponse - return this on PinsRequestidDelete
func DeleteResponse() Response {
	resp := openapi.Response(http.StatusOK, nil)
	return Response{resp, nil}
}

// PinListResponse - return a list of pins (used in PinsGet)
func PinListResponse(pinResults openapi.PinResults) Response {
	resp := openapi.Response(http.StatusOK, pinResults)
	return Response{resp, nil}
}

type Response struct {
	Resp openapi.ImplResponse
	Err  error
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
	PinsGet func(PinsGetInputData) Response
	// PinsPost - Add pin object
	PinsPost func(PinsPostInputData) Response
	// PinsRequestidDelete - Remove pin object
	PinsRequestidDelete func(PinsRequestidDeleteInputData) Response
	// PinsRequestidGet - Get pin object
	PinsRequestidGet func(PinsRequestidGetInputData) Response
	// PinsRequestidPost - Replace pin object
	PinsRequestidPost func(PinsRequestidPostInputData) Response
}

type PinsApiService struct {
	userFunctions Service
}

func (s *PinsApiService) PinsGet(ctx context.Context, cid []string, name string, match openapi.TextMatchingStrategy, status []openapi.Status, before time.Time, after time.Time, limit int32, meta map[string]string) (openapi.ImplResponse, error) {
	r := s.userFunctions.PinsGet(PinsGetInputData{
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
	return r.Resp, r.Err
}

func (s *PinsApiService) PinsPost(ctx context.Context, pin openapi.Pin) (openapi.ImplResponse, error) {
	r := s.userFunctions.PinsPost(PinsPostInputData{
		ctx: ctx,
		pin: pin,
	})
	return r.Resp, r.Err
}

func (s *PinsApiService) PinsRequestidDelete(ctx context.Context, requestid string) (openapi.ImplResponse, error) {
	r := s.userFunctions.PinsRequestidDelete(PinsRequestidDeleteInputData{
		ctx:       ctx,
		requestid: requestid,
	})
	return r.Resp, r.Err
}

func (s *PinsApiService) PinsRequestidGet(ctx context.Context, requestid string) (openapi.ImplResponse, error) {
	r := s.userFunctions.PinsRequestidGet(PinsRequestidGetInputData{
		ctx:       ctx,
		requestid: requestid,
	})
	return r.Resp, r.Err
}

func (s *PinsApiService) PinsRequestidPost(ctx context.Context, requestid string, pin openapi.Pin) (openapi.ImplResponse, error) {
	r := s.userFunctions.PinsRequestidPost(PinsRequestidPostInputData{
		ctx:       ctx,
		requestid: requestid,
		pin:       pin,
	})
	return r.Resp, r.Err
}
