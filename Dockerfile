FROM golang:1.20 as builder
ARG OS
ARG ARCH
ARG CGO_ENABLED

WORKDIR /go/src/app
COPY . .
RUN go get
RUN make build ARCH=$ARCH OS=$OS CGO_ENABLED=$CGO_ENABLED

FROM scratch
WORKDIR /
COPY --from=builder /go/src/app/kbot .
COPY --from=alpine:latest /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["./kbot", "start"]
