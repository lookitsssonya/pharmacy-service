package transport

import (
	"context"

	"github.com/go-kit/kit/endpoint"

	"example.com/pharmacy-service/medicine"
)

type CreateRequest struct {
	Medicine medicine.Medicine
}

type CreateResponse struct {
	ID  string `json:"id"`
	Err error  `json:"error,omitempty"`
}

type GetByIDRequest struct {
	ID string `json:"id"`
}

type GetByIDResponse struct {
	Medicine medicine.Medicine `json:"medicine"`
	Err      error             `json:"error,omitempty"`
}

type SellRequest struct {
	ID    string `json:"id"`
	Count int    `json:"count"`
}

type SellResponse struct {
	Err error `json:"error,omitempty"`
}

type Endpoints struct {
	Create  endpoint.Endpoint
	GetByID endpoint.Endpoint
	Sell    endpoint.Endpoint
}

func MakeEndpoints(s medicine.Service) Endpoints {
	return Endpoints{
		Create:  makeCreateEndpoint(s),
		GetByID: makeGetByIDEndpoint(s),
		Sell:    makeSellEndpoint(s),
	}
}

func makeCreateEndpoint(s medicine.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateRequest)
		id, err := s.Create(ctx, req.Medicine)
		return CreateResponse{ID: id, Err: err}, nil
	}
}

func makeGetByIDEndpoint(s medicine.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetByIDRequest)
		m, err := s.GetByID(ctx, req.ID)
		return GetByIDResponse{Medicine: m, Err: err}, nil
	}
}

func makeSellEndpoint(s medicine.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(SellRequest)
		err := s.Sell(ctx, req.ID, req.Count)
		return SellResponse{Err: err}, nil
	}
}
