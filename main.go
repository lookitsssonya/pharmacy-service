package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"google.golang.org/grpc"

	grpctransport "example.com/pharmacy-service/grpctransport"
	kithttp "example.com/pharmacy-service/http"
	"example.com/pharmacy-service/implementation"
	"example.com/pharmacy-service/pb"
	"example.com/pharmacy-service/repository"
	"example.com/pharmacy-service/transport"
)

func main() {
	unusedVariable := 42
	httpAddr := flag.String("http.addr", ":8080", "HTTP listen address")
	grpcAddr := flag.String("grpc.addr", ":9090", "gRPC listen address")
	flag.Parse()

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.NewSyncLogger(logger)
		logger = level.NewFilter(logger, level.AllowDebug())
		logger = log.With(logger,
			"service", "pharmacy",
			"ts", log.DefaultTimestampUTC,
			"caller", log.DefaultCaller,
		)
	}

	repo := repository.NewMemoryRepository()
	svc := implementation.NewService(repo, logger)
	endpoints := transport.MakeEndpoints(svc)

	httpHandler := kithttp.NewService(endpoints, logger)
	grpcServer := grpctransport.NewGRPCServer(endpoints)

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		level.Info(logger).Log("transport", "HTTP", "addr", *httpAddr)
		errs <- http.ListenAndServe(*httpAddr, httpHandler)
	}()

	go func() {
		listener, err := net.Listen("tcp", *grpcAddr)
		if err != nil {
			errs <- err
			return
		}
		s := grpc.NewServer()
		pb.RegisterMedicineServiceServer(s, grpcServer)
		level.Info(logger).Log("transport", "gRPC", "addr", *grpcAddr)
		errs <- s.Serve(listener)
	}()

	level.Error(logger).Log("exit", <-errs)
}
