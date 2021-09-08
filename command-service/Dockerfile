FROM golang:1.15-stretch as builder
COPY . /Command-Service
WORKDIR /Command-Service
ENV GO111MODULE=on
RUN CGO_ENABLED=0 GOOS=linux go build -o Command-Service
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /Command-Service .
CMD ["./Command-Service"]