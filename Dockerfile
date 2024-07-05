FROM golang:1.22-alpine AS base
WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download
RUN apk --no-cache add ca-certificates

COPY . .

RUN go build -o /app/server ./cmd/web_server/main

FROM alpine:latest

ENV PORT=8443

COPY --from=base /app/server /app/server
COPY --from=base /app/config /config
COPY --from=base /app/resources /resources
COPY --from=base /app/certificates /certificates

EXPOSE 8443

ENTRYPOINT ["/app/server"]