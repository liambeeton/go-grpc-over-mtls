# Keystone Global Banking

Keystone Global Banking is a mock banking application developed to demonstrate mutual TLS (mTLS) authentication.

This project is written in Go and serves as a practical example for developers interested in learning about secure communication between clients and servers using mTLS.

## Features

- mTLS Authentication: Ensures both client and server mutually verify each other's identities.
- Banking Operations: Simulates basic banking operations such as account creation, balance checking, and funds transfer.
- Secure API Endpoints: All interactions with the app are secured using mTLS.

## Requirements

- Go 1.16+
- OpenSSL: For generating certificates
- Docker (Optional): For containerized deployment

## Go Dependencies

Install the following dependencies using `go get` for gRPC and Protocol Buffers support.

```sh
go get google.golang.org/grpc
go get google.golang.org/protobuf
```

## Proto Files

Generate Go code from proto files.

```sh
make proto-compile
```

## Installation

### Clone the Repository

```sh
git clone https://github.com/liambeeton/go-grpc-over-mtls.git
cd go-grpc-over-mtls
```

### Generate Certificates

Use OpenSSL to generate the necessary client and server certificates.

- ECC Encryption: Use [create-certs-ecc.sh](https://github.com/liambeeton/go-grpc-over-mtls/blob/main/create-certs-ecc.sh) script for generating certificates.
- RSA Encryption: Use [create-certs-rsa.sh](https://github.com/liambeeton/go-grpc-over-mtls/blob/main/create-certs-rsa.sh) script for generating certificates.

### Build the Application

```sh
make build
```

## Usage

### Run the Server

```sh
make run
```

### Run the Client

```sh
make run-client
```

## Configuration

Set the following environment variables in `.envrc` to override the default config.

```sh
export HOST="server"
export PORT="8443"
export CA_FILE="/usr/bin/certs-ecc/ca.crt"
export CLIENT_CERT_FILE="/usr/bin/certs-ecc/client.crt"
export CLIENT_KEY_FILE="/usr/bin/certs-ecc/client.key"
export SERVER_CERT_FILE="/usr/bin/certs-ecc/server.crt"
export SERVER_KEY_FILE="/usr/bin/certs-ecc/server.key"
```

## Contributing

We welcome contributions! Please fork the repository and create a pull request with your changes.

## License

This project is licensed under the MIT License.

## Acknowledgements

Special thanks to the contributors and the Go community for their support and resources.

---

Feel free to explore the code, test the mTLS authentication, and understand the implementation details provided in this repository. Happy coding!
