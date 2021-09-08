FROM golang:1.15-stretch as builder
COPY . /Auth-Service
WORKDIR /Auth-Service
ENV GO111MODULE=on
RUN CGO_ENABLED=0 GOOS=linux go build -o Auth-Service
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /Auth-Service .
CMD ["./Auth-Service"]