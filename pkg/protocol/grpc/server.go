package grpc

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/brewm/gobrewmmer/pkg/api/brewmmer"
	"google.golang.org/grpc"
)

// RunServer runs gRPC service to publish ToDo service
func RunServer(ctx context.Context, recepieAPI brewmmer.RecepieServiceServer, sessionAPI brewmmer.SessionServiceServer, port string) error {
	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	// register service
	server := grpc.NewServer()
	brewmmer.RegisterRecepieServiceServer(server, recepieAPI)
	brewmmer.RegisterSessionServiceServer(server, sessionAPI)

	// graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			// sig is a ^C, handle it
			log.Println("shutting down gRPC server...")

			server.GracefulStop()

			<-ctx.Done()
		}
	}()

	// start gRPC server
	log.Println("starting gRPC server...")
	return server.Serve(listen)
}
