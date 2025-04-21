FROM golang:1.23-alpine AS builder
RUN apk add --no-cache ca-certificates
WORKDIR /opt/neoserve
COPY go.mod go.sum ./
RUN go mod tidy
COPY . .
RUN go get -d -v ./...
RUN go build -v -o neoserve .

# Run main image
FROM alpine:3.19
RUN apk add --no-cache ca-certificates
WORKDIR /opt/neoserve
COPY --from=builder /opt/neoserve/static/* /opt/neoserve/static/
COPY --from=builder /opt/neoserve/neoserve /usr/bin/neoserve
COPY --from=builder /opt/neoserve/uploads/* /opt/neoserve/uploads/
COPY --from=builder /opt/neoserve/public/* /opt/neoserve/public/
EXPOSE 8080
CMD ["neoserve"]
