#!/bin/bash
# Author: liambeeton - https://github.com/liambeeton
# Script to generate ECC mTLS certificates

# Exit on first error
set -e

this_script_directory="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd ${this_script_directory}

mkdir -p certs-ecc
cd certs-ecc

# Create a Private Key for the CA
openssl ecparam -name prime256v1 -genkey -noout -out ca.key
cat ca.key

# Create a Self-Signed Certificate for the CA
openssl req -x509 -new -nodes -key ca.key -subj "/C=US/ST=New York/L=New York City/O=Example CA Inc./CN=Example Root CA" -days 365 -out ca.crt
openssl x509 -in ca.crt -text

# Create a Private Key for the Server
openssl ecparam -name prime256v1 -genkey -noout -out server.key
cat server.key

# Create a Certificate Signing Request (CSR) config for the Server
cat > server.conf <<EOF
[ req ]
default_bits        = 2048
default_keyfile     = server.key
default_md          = sha256
prompt              = no
distinguished_name  = req_distinguished_name
req_extensions      = req_ext

[ req_distinguished_name ]
C                   = ZA
ST                  = Western Cape
L                   = Cape Town
O                   = Keystone Global Banking
OU                  = Finance
CN                  = bank.kgb.rip

[ req_ext ]
keyUsage            = keyEncipherment, dataEncipherment
extendedKeyUsage    = serverAuth
subjectAltName      = @alt_names

[ alt_names ]
DNS.1               = localhost
DNS.2               = server
DNS.3               = bank.kgb.rip
IP.1                = 127.0.0.1
EOF
cat server.conf

# Create a Certificate Signing Request (CSR) for the Server
openssl req -new -key server.key -out server.csr -config server.conf
cat server.csr

# Sign the Server CSR with the CA Certificate to Generate the Server Certificate
openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt -days 90 -extfile server.conf -extensions req_ext
openssl x509 -in server.crt -text

# Create a Private Key for the Client
openssl ecparam -name prime256v1 -genkey -noout -out client.key
cat client.key

# Create a Certificate Signing Request (CSR) config for the Client
cat > client.conf <<EOF
[ req ]
default_bits        = 2048
default_keyfile     = client.key
default_md          = sha256
prompt              = no
distinguished_name  = req_distinguished_name
req_extensions      = req_ext

[ req_distinguished_name ]
C                   = ZA
ST                  = Western Cape
L                   = Cape Town
O                   = Keystone Global Banking
OU                  = Finance
CN                  = bank.kgb.rip

[ req_ext ]
keyUsage            = keyEncipherment, dataEncipherment
extendedKeyUsage    = clientAuth
subjectAltName      = @alt_names

[ alt_names ]
DNS.1               = localhost
DNS.2               = server
DNS.3               = bank.kgb.rip
IP.1                = 127.0.0.1
EOF
cat client.conf

# Create a Certificate Signing Request (CSR) for the Client
openssl req -new -key client.key -out client.csr -config client.conf
cat client.csr

# Sign the Client CSR with the CA Certificate to Generate the Client Certificate
openssl x509 -req -in client.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out client.crt -days 90 -extfile client.conf -extensions req_ext
openssl x509 -in client.crt -text

# Verify the server certificate
openssl verify -CAfile ca.crt server.crt

# Verify the client certificate
openssl verify -CAfile ca.crt client.crt
