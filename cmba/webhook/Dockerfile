# Multi-stage build for the CMBA admission webhook.
# Final image is distroless static for minimal attack surface.
FROM golang:1.22-alpine AS build

WORKDIR /src
COPY go.mod ./
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o /out/cmba-webhook .

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/cmba-webhook /usr/local/bin/cmba-webhook
EXPOSE 8443
USER nonroot:nonroot
ENTRYPOINT ["/usr/local/bin/cmba-webhook"]
