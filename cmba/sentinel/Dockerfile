# Multi-stage build for the cmba-sentinel sidecar.
FROM golang:1.22-alpine AS build

WORKDIR /src
COPY go.mod ./
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o /out/cmba-sentinel .

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/cmba-sentinel /usr/local/bin/cmba-sentinel
USER nonroot:nonroot
ENTRYPOINT ["/usr/local/bin/cmba-sentinel"]
