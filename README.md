# Neoserve

The solution for my homelab on serving files / links written in Golang in an attempt to learn the platform.

So far the server supports a few functions listed below

- Upload any file type.
- Max Upload Limit (Including Headers etc)
- Discord Webhooks for Upload logs.
- Docker
- Kubernetes
- Scales perfectly
- Discord Webhook passthroughs (top level domain)

# Running

You can run Neoserve standalone or with Docker as of this release.

## Docker

Run image and store files in a named Docker volume
```shell
# Build docker image and run the stack
docker compose up --build -d

# Run and attach to container
docker compose up --build

# Run alone with docker
docker run --rm -p 8080:8080 -v $(pwd)/config.yaml:/opt/neoserve/config.yaml neoserve:latest 
```

## Standalone

Run Neoserve in standalone mode.
```shell
# Build Neoserve
go mod tidy
go get -d -v ./...
go build -v -o neoserve .

# Run (binds to port 8080)
./neoserve
```