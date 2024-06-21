FROM alpine:3 AS user-and-certs

RUN apk add -U --no-cache ca-certificates

RUN addgroup -S kgb -g 2000 && adduser -S kgb -G kgb -u 2000

FROM golang:1.22.4 AS base

WORKDIR /workspace/
COPY . .

RUN make linux

FROM scratch

COPY --from=user-and-certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=user-and-certs /etc/group /etc/group
COPY --from=user-and-certs /etc/passwd /etc/passwd

USER kgb:kgb

WORKDIR /usr/bin/
COPY --from=base /workspace/client .

ENTRYPOINT [ "client" ]
