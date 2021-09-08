FROM golang:1.15-stretch as builder
COPY . /Vis-Service
WORKDIR /Vis-Service
ENV GO111MODULE=on
RUN CGO_ENABLED=0 GOOS=linux go build -o Vis-Service
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /Vis-Service .
CMD ["./Vis-Service"]