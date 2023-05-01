package servergrpc

import (
	"fmt"
	"log"
	"net"

	"github.com/BillyBones007/pwdm_server/internal/grpcservices"
	"github.com/BillyBones007/pwdm_server/internal/storage"
	pb "github.com/BillyBones007/pwdm_service_api/api"
	"google.golang.org/grpc"
)

// ServerGRPC - the main structure of the server gRPC.
type Server struct {
	Config     *ServerConfig
	Storage    storage.Storage
	GRPCServer *grpc.Server
}

// NewServer - returns a pointer to the Server.
func NewServer() *Server {
	server := Server{}
	server.Config = InitServerConfig()
	// TODO: server.Storage = InitStorage()
	server.GRPCServer = grpc.NewServer()
	pb.RegisterAuthServiceServer(server.GRPCServer, grpcservices.NewAuthService(&server.Storage))
	pb.RegisterGiveTakeServiceServer(server.GRPCServer, grpcservices.NewGiveTakeService(&server.Storage))
	pb.RegisterDeleteServiceServer(server.GRPCServer, grpcservices.NewDeleteService(&server.Storage))
	pb.RegisterShowInfoServiceServer(server.GRPCServer, grpcservices.NewShowInfoService(&server.Storage))

	return &server
}

// StartServer - starting the gRPC server.
func (s *Server) StartServer() {
	listen, err := net.Listen("tcp", s.Config.PortgRPC)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		fmt.Println("gRPC server is started...")
		if err := s.GRPCServer.Serve(listen); err != nil {
			log.Fatal(err)
		}
	}()
}

// Shutdown - gracefully stoped the server.
func (s *Server) Shutdown() {
	log.Println("Interrupt signal received, server shutting down...")
	s.GRPCServer.GracefulStop()
	s.Storage.Close()
}
