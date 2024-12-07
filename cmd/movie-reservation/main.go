package main

import (
	"context"
	"flag"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/harveytvt/movie-reservation-system/gen/go/api/movie_reservation/v1"
	"github.com/harveytvt/movie-reservation-system/internal/auth"
	"github.com/harveytvt/movie-reservation-system/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

var (
	grpcServerEndpoint = flag.String("grpc-listen-port", "0.0.0.0:50051", "gRPC server endpoint")
	done               chan struct{}
)

func listenGrpc() {
	defer close(done)
	lis, err := net.Listen("tcp", *grpcServerEndpoint)
	if err != nil {
		panic(err)
	}

	server := grpc.NewServer()
	service := service.NewService()
	movie_reservation.RegisterMovieReservationServiceServer(server, service)
	grpc_health_v1.RegisterHealthServer(server, service)
	reflection.Register(server)

	if err = server.Serve(lis); err != nil {
		panic(err)
	}
}

func listenHttp() {
	defer close(done)
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux(
		runtime.WithErrorHandler(runtime.DefaultHTTPErrorHandler),
		runtime.WithMetadata(auth.AuthAnnotator),
		runtime.WithForwardResponseOption(auth.ForwardResponseOption),
	)
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := movie_reservation.RegisterMovieReservationServiceHandlerFromEndpoint(ctx, mux, *grpcServerEndpoint, opts)

	// custom handler
	service.HandleUpload(mux)

	if err != nil {
		panic(err)
	}

	if err := http.ListenAndServe("0.0.0.0:8081", mux); err != nil {
		panic(err)
	}
}

func main() {
	flag.Parse()

	go listenGrpc()
	go listenHttp()

	select {
	case <-done:
		return
	}

}
