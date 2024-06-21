package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/liambeeton/go-grpc-over-mtls/pb/message"
	"github.com/liambeeton/go-grpc-over-mtls/pb/service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

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
	cred := newClientTLS(conf)

	// Dial the gRPC server with the given credentials
	log.Printf("Client connecting to %s:%d", conf.Host, conf.Port)
	conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", conf.Host, conf.Port), grpc.WithTransportCredentials(cred))
	if err != nil {
		log.Fatalf("Unable to connect gRPC channel %v", err)
	}

	// Close the listener when containing function terminates
	defer func() {
		err = conn.Close()
		if err != nil {
			log.Printf("Unable to close gRPC channel %v", err)
		}
	}()

	// Create the gRPC client
	c := service.NewBankServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Create account
	createResp, err := c.CreateAccount(ctx, &message.CreateAccountRequest{AccountId: "12345"})
	if err != nil {
		log.Fatalf("Could not create account %v", err)
	}
	log.Printf("Account created %v", createResp.AccountId)

	// Deposit
	depositResp, err := c.Deposit(ctx, &message.DepositRequest{AccountId: "12345", Amount: 100.0})
	if err != nil {
		log.Fatalf("Could not deposit %v", err)
	}
	log.Printf("New balance after deposit %v", depositResp.NewBalance)

	// Get balance
	balanceResp, err := c.GetBalance(ctx, &message.GetBalanceRequest{AccountId: "12345"})
	if err != nil {
		log.Fatalf("Could not get balance %v", err)
	}
	log.Printf("Balance %v", balanceResp.Balance)

	// Withdraw
	withdrawResp, err := c.Withdraw(ctx, &message.WithdrawRequest{AccountId: "12345", Amount: 50.0})
	if err != nil {
		log.Fatalf("Could not withdraw %v", err)
	}
	log.Printf("New balance after withdrawal %v", withdrawResp.NewBalance)
}

func newClientTLS(c *config) credentials.TransportCredentials {
	// Load the client certificate and its key
	clientCert, err := tls.LoadX509KeyPair(c.CertFile, c.KeyFile)
	if err != nil {
		log.Fatalf("Failed to load client certificate and key %v", err)
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
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      certPool,
		MinVersion:   tls.VersionTLS13,
		MaxVersion:   tls.VersionTLS13,
	}

	// Return new TLS credentials based on the TLS configuration
	return credentials.NewTLS(tlsConfig)
}
