package grpctransport

import (
	"context"

	grpctransport "github.com/go-kit/kit/transport/grpc"

	"example.com/pharmacy-service/medicine"
	"example.com/pharmacy-service/pb"
	"example.com/pharmacy-service/transport"
)

type grpcServer struct {
	pb.UnimplementedMedicineServiceServer
	create  grpctransport.Handler
	getByID grpctransport.Handler
	sell    grpctransport.Handler
}

func NewGRPCServer(endpoints transport.Endpoints) pb.MedicineServiceServer {
	return &grpcServer{
		create:  grpctransport.NewServer(endpoints.Create, decodeGRPCCreateRequest, encodeGRPCCreateResponse),
		getByID: grpctransport.NewServer(endpoints.GetByID, decodeGRPCGetByIDRequest, encodeGRPCGetByIDResponse),
		sell:    grpctransport.NewServer(endpoints.Sell, decodeGRPCSellRequest, encodeGRPCSellResponse),
	}
}

func (s *grpcServer) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	_, resp, err := s.create.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.CreateResponse), nil
}

func (s *grpcServer) GetByID(ctx context.Context, req *pb.GetByIDRequest) (*pb.GetByIDResponse, error) {
	_, resp, err := s.getByID.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.GetByIDResponse), nil
}

func (s *grpcServer) Sell(ctx context.Context, req *pb.SellRequest) (*pb.SellResponse, error) {
	_, resp, err := s.sell.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.SellResponse), nil
}

func decodeGRPCCreateRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.CreateRequest)
	return transport.CreateRequest{Medicine: medicineFromPB(req.Medicine)}, nil
}

func decodeGRPCGetByIDRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.GetByIDRequest)
	return transport.GetByIDRequest{ID: req.Id}, nil
}

func decodeGRPCSellRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.SellRequest)
	return transport.SellRequest{ID: req.Id, Count: int(req.Count)}, nil
}

func encodeGRPCCreateResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(transport.CreateResponse)
	errStr := ""
	if resp.Err != nil {
		errStr = resp.Err.Error()
	}
	return &pb.CreateResponse{Id: resp.ID, Error: errStr}, nil
}

func encodeGRPCGetByIDResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(transport.GetByIDResponse)
	errStr := ""
	if resp.Err != nil {
		errStr = resp.Err.Error()
	}
	return &pb.GetByIDResponse{Medicine: medicineToPB(resp.Medicine), Error: errStr}, nil
}

func encodeGRPCSellResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(transport.SellResponse)
	errStr := ""
	if resp.Err != nil {
		errStr = resp.Err.Error()
	}
	return &pb.SellResponse{Error: errStr}, nil
}

func medicineToPB(m medicine.Medicine) *pb.Medicine {
	return &pb.Medicine{
		Id:           m.ID,
		Name:         m.Name,
		Manufacturer: m.Manufacturer,
		Price:        m.Price,
		Prescription: m.Prescription,
		Quantity:     int32(m.Quantity),
	}
}

func medicineFromPB(m *pb.Medicine) medicine.Medicine {
	if m == nil {
		return medicine.Medicine{}
	}
	return medicine.Medicine{
		ID:           m.Id,
		Name:         m.Name,
		Manufacturer: m.Manufacturer,
		Price:        m.Price,
		Prescription: m.Prescription,
		Quantity:     int(m.Quantity),
	}
}
