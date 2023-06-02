package servergrpc

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"os"

	"github.com/BillyBones007/pwdm_server/internal/grpcservices"
	"github.com/BillyBones007/pwdm_server/internal/logger"
	"github.com/BillyBones007/pwdm_server/internal/storage"
	"github.com/BillyBones007/pwdm_server/internal/storage/postgres"
	"github.com/BillyBones007/pwdm_server/internal/tools/tokentools"
	pb "github.com/BillyBones007/pwdm_service_api/api"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Server - the main structure of the server gRPC.
type Server struct {
	Config       *ServerConfig
	Storage      storage.Storage
	TokenTools   *tokentools.JWTTools
	GRPCServer   *grpc.Server
	Interceptors *grpcservices.InterceptorsService
	Logger       *logrus.Logger
}

// NewServer - returns a pointer to the Server.
func NewServer() *Server {
	server := Server{}
	server.Config = InitServerConfig()
	server.Logger = logger.NewLogger()
	server.TokenTools = tokentools.NewJWTTools()
	stor, err := postgres.NewClientPostgres(server.Config.DSN)
	if err != nil {
		server.Logger.WithField("err", err).Fatalf("Failed database: %s", err)
	}
	server.Storage = stor

	// Interceptors - the pointer to InterceptorsService.
	// Uses tools for working with jwt.
	server.Interceptors = grpcservices.NewInterceptorsService(server.TokenTools, server.Logger)

	cert, err := tls.LoadX509KeyPair("cert/server.crt", "cert/server.key")
	if err != nil {
		server.Logger.WithField("err", err).Fatal("Failed to load server certificates")
	}
	caCert, err := os.ReadFile("cert/ca.crt")
	if err != nil {
		server.Logger.WithField("err", err).Fatal("Failed to load server certificates")
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.NoClientCert,
		ClientCAs:    caCertPool,
	}
	creds := credentials.NewTLS(tlsConfig)
	opts := []grpc.ServerOption{grpc.Creds(creds), grpc.UnaryInterceptor(server.Interceptors.AuthInterceptor)}
	server.GRPCServer = grpc.NewServer(opts...)

	pb.RegisterAuthServiceServer(server.GRPCServer, grpcservices.NewAuthService(server.Storage, server.TokenTools, server.Logger))
	pb.RegisterGiveTakeServiceServer(server.GRPCServer, grpcservices.NewGiveTakeService(server.Storage, server.TokenTools, server.Logger))
	pb.RegisterUpdateServiceServer(server.GRPCServer, grpcservices.NewUpdateService(server.Storage, server.TokenTools, server.Logger))
	pb.RegisterDeleteServiceServer(server.GRPCServer, grpcservices.NewDeleteService(server.Storage, server.TokenTools, server.Logger))
	pb.RegisterShowInfoServiceServer(server.GRPCServer, grpcservices.NewShowInfoService(server.Storage, server.TokenTools, server.Logger))

	return &server
}

// StartServer - starting the gRPC server.
func (s *Server) StartServer() {
	listen, err := net.Listen("tcp", s.Config.PortgRPC)
	if err != nil {
		s.Logger.WithField("err", err).Fatal("The server crashed")
	}

	go func() {
		s.Logger.WithFields(logrus.Fields{
			"grpc_port": s.Config.PortgRPC,
			"dsn":       s.Config.DSN,
		}).Info("Server gRPC is started")
		fmt.Println("Server gRPC is started...")
		if err := s.GRPCServer.Serve(listen); err != nil {
			s.Logger.WithField("err", err).Fatal("The server crashed")
		}
	}()
}

// Shutdown - gracefully stoped the server.
func (s *Server) Shutdown() {
	s.Logger.Info("Interrupt signal received, server shutting down")
	s.GRPCServer.GracefulStop()
	s.Storage.Close()
}
