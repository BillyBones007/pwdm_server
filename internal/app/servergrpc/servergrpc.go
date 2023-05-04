package servergrpc

import (
	"fmt"
	"log"
	"net"

	"github.com/BillyBones007/pwdm_server/internal/grpcservices"
	"github.com/BillyBones007/pwdm_server/internal/storage"
	"github.com/BillyBones007/pwdm_server/internal/storage/postgres"
	"github.com/BillyBones007/pwdm_server/internal/tools/tokentools"
	pb "github.com/BillyBones007/pwdm_service_api/api"
	"google.golang.org/grpc"
)

// ServerGRPC - the main structure of the server gRPC.
type Server struct {
	Config       *ServerConfig
	Storage      storage.Storage
	TokenTools   *tokentools.JWTTools
	GRPCServer   *grpc.Server
	Interceptors *grpcservices.InterceptorsService
}

// NewServer - returns a pointer to the Server.
func NewServer() *Server {
	server := Server{}
	server.Config = InitServerConfig()
	server.TokenTools = tokentools.NewJWTTools()
	server.Storage = postgres.NewClientPostgres(server.Config.DSN)

	// Interceptors - the pointer to InterceptorsService.
	// Uses tools for working with jwt.
	server.Interceptors = grpcservices.NewInterceptorsService(server.TokenTools)
	server.GRPCServer = grpc.NewServer(grpc.UnaryInterceptor(server.Interceptors.AuthInterceptor))

	pb.RegisterAuthServiceServer(server.GRPCServer, grpcservices.NewAuthService(server.Storage, server.TokenTools))
	pb.RegisterGiveTakeServiceServer(server.GRPCServer, grpcservices.NewGiveTakeService(server.Storage, server.TokenTools))
	pb.RegisterDeleteServiceServer(server.GRPCServer, grpcservices.NewDeleteService(server.Storage, server.TokenTools))
	pb.RegisterShowInfoServiceServer(server.GRPCServer, grpcservices.NewShowInfoService(server.Storage, server.TokenTools))

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
