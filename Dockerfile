FROM golang:1.22-alpine AS base
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

FROM base AS build-server
RUN go build -o /app/server ./cmd/web_server/main

FROM scratch
COPY --from=build-server /app/server /app/server
COPY --from=build-server /app/config /config
COPY --from=build-server /app/resources /resources
COPY --from=build-server /app/certificates /certificates
EXPOSE 8443

ENTRYPOINT ["/app/server"]