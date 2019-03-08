package grpc

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"

	"google.golang.org/grpc"

	"github.com/brewm/gobrewmmer/pkg/api/recepie"
	"github.com/brewm/gobrewmmer/pkg/api/session"
)

// RunServer runs gRPC service to publish ToDo service
func RunServer(ctx context.Context, recepieAPI recepie.RecepieServiceServer, sessionAPI session.SessionServiceServer, port string) error {
	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	// register service
	server := grpc.NewServer()
	recepie.RegisterRecepieServiceServer(server, recepieAPI)
	session.RegisterSessionServiceServer(server, sessionAPI)

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
