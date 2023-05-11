FROM golang:alpine AS builder
WORKDIR /app
COPY . .
RUN apk add git
RUN mkdir -p /app/bin
RUN GOOS=linux GOARCH=amd64 go build \
  -o bin \
  ./

FROM alpine:latest
RUN apk --no-cache --update add ca-certificates
COPY --from=builder /app/bin/* /app/bin/
EXPOSE 8080
CMD ["/usr/bin/false"]