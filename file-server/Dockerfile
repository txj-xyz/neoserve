FROM golang:1.22-alpine as builder

WORKDIR /file-server

# Deps / Download etc
COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

# Build dist
RUN go build -o neoserve .

# -------------------------------
# Production stuff things here
# -------------------------------
# Run main image
FROM alpine:latest

# Install reqs
RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /file-server/neoserve .

EXPOSE 8080

CMD ["./neoserve"]