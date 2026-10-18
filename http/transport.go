package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"

	"example.com/pharmacy-service/transport"
)

func NewService(endpoints transport.Endpoints, logger log.Logger) http.Handler {
	options := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(encodeError),
	}

	createHandler := kithttp.NewServer(endpoints.Create, decodeCreateRequest, encodeResponse, options...)
	getByIDHandler := kithttp.NewServer(endpoints.GetByID, decodeGetByIDRequest, encodeResponse, options...)
	sellHandler := kithttp.NewServer(endpoints.Sell, decodeSellRequest, encodeResponse, options...)

	r := mux.NewRouter()
	r.Handle("/medicines", createHandler).Methods("POST")
	r.Handle("/medicines/{id}", getByIDHandler).Methods("GET")
	r.Handle("/medicines/{id}/sell", sellHandler).Methods("POST")

	return r
}

func decodeCreateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req transport.CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req.Medicine); err != nil {
		return nil, err
	}
	return req, nil
}

func decodeGetByIDRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	return transport.GetByIDRequest{ID: vars["id"]}, nil
}

func decodeSellRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	var body struct {
		Count int `json:"count"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, err
	}
	return transport.SellRequest{ID: vars["id"], Count: body.Count}, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"error": err.Error()})
}
