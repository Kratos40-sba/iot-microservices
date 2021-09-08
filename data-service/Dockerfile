FROM golang:1.15-stretch as builder
COPY . /Data-Service
WORKDIR /Data-Service
ENV GO111MODULE=on
RUN CGO_ENABLED=0 GOOS=linux go build -o Data-Service
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /Data-Service .
CMD ["./Data-Service"]