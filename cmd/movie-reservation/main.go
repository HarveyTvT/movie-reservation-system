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

	srv := grpc.NewServer()
	movie_reservation.RegisterMovieReservationServiceServer(srv, service.NewService())
	reflection.Register(srv)

	if err = srv.Serve(lis); err != nil {
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
