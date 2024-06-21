package main

import (
	"fmt"

	"github.com/jessevdk/go-flags"
)

type config struct {
	Host       string `long:"host" env:"HOST" default:"" description:"The IP/DNS of the machine that the application is running on"`
	Port       int    `long:"port" env:"PORT" default:"8080" description:"The port the service is listening on"`
	TLSEnabled bool   `long:"tls-enabled" env:"TLS_ENABLED" description:"When set, mTLS is used for establishing a connection to the server. Please set the 'keyFile', 'certFile' and 'caFile'"`
	CaFile     string `long:"ca-file" env:"CA_FILE" default:"" description:"(Optional) The path to the PEM file containing the CA certificate"`
	KeyFile    string `long:"key-file" env:"KEY_FILE" default:"" description:"(Optional) The path to the PEM file containing a private key"`
	CertFile   string `long:"cert-file" env:"CERT_FILE" default:"" description:"(Optional) The path to the PEM file containing a certificate"`
}

func newConfig() (*config, error) {
	conf := &config{}
	parser := flags.NewParser(conf, flags.Default)
	if _, err := parser.Parse(); err != nil {
		return nil, fmt.Errorf("failed to parse config: %v", err)
	}
	return conf, nil
}
