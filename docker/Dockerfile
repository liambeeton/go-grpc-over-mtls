FROM golang:1.22.1

ARG PLATFORM
ARG PROTOC_VERSION

RUN apt-get update && apt-get install -y unzip

# By default Intel chipset (x86_64) is assumed but if the host device is an Apple
# silicon (arm) chipset based then a relevant (aarch_64) release file is used.

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.31.0
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0

RUN export ZIP=x86_64 && \
  if [ ${PLATFORM} = "arm64" ]; then export ZIP=aarch_64; fi && \
  wget --quiet https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/protoc-${PROTOC_VERSION}-linux-${ZIP}.zip && \
  unzip -o protoc-${PROTOC_VERSION}-linux-${ZIP}.zip -d /usr/local bin/protoc && \
  unzip -o protoc-${PROTOC_VERSION}-linux-${ZIP}.zip -d /usr/local 'include/*'