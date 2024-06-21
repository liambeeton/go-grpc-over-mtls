package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/liambeeton/go-grpc-over-mtls/pb/message"
	"github.com/liambeeton/go-grpc-over-mtls/pb/service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

type server struct {
	service.UnimplementedBankServiceServer
	accounts map[string]float64
}

func (s *server) CreateAccount(_ context.Context, req *message.CreateAccountRequest) (*message.CreateAccountResponse, error) {
	s.accounts[req.AccountId] = 0
	return &message.CreateAccountResponse{AccountId: req.AccountId}, nil
}

func (s *server) GetBalance(_ context.Context, req *message.GetBalanceRequest) (*message.GetBalanceResponse, error) {
	balance, exists := s.accounts[req.AccountId]
	if !exists {
		return nil, status.Error(codes.NotFound, "Account not found")
	}
	return &message.GetBalanceResponse{AccountId: req.AccountId, Balance: balance}, nil
}

func (s *server) Deposit(_ context.Context, req *message.DepositRequest) (*message.DepositResponse, error) {
	_, exists := s.accounts[req.AccountId]
	if !exists {
		return nil, status.Error(codes.NotFound, "Account not found")
	}
	s.accounts[req.AccountId] += req.Amount
	return &message.DepositResponse{NewBalance: s.accounts[req.AccountId]}, nil
}

func (s *server) Withdraw(_ context.Context, req *message.WithdrawRequest) (*message.WithdrawResponse, error) {
	balance, exists := s.accounts[req.AccountId]
	if !exists {
		return nil, status.Error(codes.NotFound, "Account not found")
	}
	if balance < req.Amount {
		return nil, status.Error(codes.FailedPrecondition, "Insufficient funds")
	}
	s.accounts[req.AccountId] -= req.Amount
	return &message.WithdrawResponse{NewBalance: s.accounts[req.AccountId]}, nil
}

func main() {
	// Load config
	conf, err := newConfig()
	if err != nil {
		log.Fatalf("Failed to load config %v", err)
	}

	// Print config
	fmt.Printf("Host: %s\n", conf.Host)
	fmt.Printf("Port: %d\n", conf.Port)
	fmt.Printf("CA File: %s\n", conf.CaFile)
	fmt.Printf("Key File: %s\n", conf.KeyFile)
	fmt.Printf("Cert File: %s\n", conf.CertFile)

	// Get TLS credentials
	cred := newServerTLS(conf)

	// Create a listener that listens to localhost port 8443
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", conf.Port))
	if err != nil {
		log.Fatalf("Failed to start listener %v", err)
	}

	// Close the listener when containing function terminates
	defer func() {
		err = lis.Close()
		if err != nil {
			log.Printf("Failed to close listener %v", err)
		}
	}()

	// Create a new gRPC server
	s := grpc.NewServer(grpc.Creds(cred))
	service.RegisterBankServiceServer(s, &server{accounts: make(map[string]float64)})

	// Start the gRPC server
	log.Printf("Server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve %v", err)
	}
}

func newServerTLS(c *config) credentials.TransportCredentials {
	// Load the server certificate and its key
	serverCert, err := tls.LoadX509KeyPair(c.CertFile, c.KeyFile)
	if err != nil {
		log.Fatalf("Failed to load server certificate and key %v", err)
	}

	// Load the CA certificate
	trustedCert, err := os.ReadFile(c.CaFile)
	if err != nil {
		log.Fatalf("Failed to load trusted certificate %v", err)
	}

	// Put the CA certificate into the certificate pool
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(trustedCert) {
		log.Fatalf("Failed to append trusted certificate to certificate pool %v", err)
	}

	// Create the TLS configuration
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		RootCAs:      certPool,
		ClientCAs:    certPool,
		MinVersion:   tls.VersionTLS13,
		MaxVersion:   tls.VersionTLS13,
	}

	// Return new TLS credentials based on the TLS configuration
	return credentials.NewTLS(tlsConfig)
}
