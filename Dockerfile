# ---- build stage ----
FROM golang:1.26-alpine AS build
WORKDIR /src
RUN apk add --no-cache git ca-certificates
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath \
    -ldflags="-s -w" \
    -o /out/server ./apps/server
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath \
    -ldflags="-s -w" \
    -o /out/worker ./apps/server/worker

# ---- runtime: API + worker binaries in one image ----
# Alpine (not distroless): Heroku Container wraps process commands as
# `/bin/sh -c …`, which distroless cannot run.
FROM alpine:3.21
RUN apk add --no-cache ca-certificates
COPY --from=build /out/server /server
COPY --from=build /out/worker /worker
EXPOSE 8080
USER nobody
# Default process is the API; Heroku overrides via heroku.yml for worker
CMD ["/server"]
