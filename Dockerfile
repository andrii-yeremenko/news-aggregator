FROM golang:1.22-alpine AS base
LABEL maintainer="Andrii Yeremenko"
WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download
RUN apk --no-cache add ca-certificates

COPY aggregator aggregator
COPY cmd/web_server cmd/web_server
COPY storage storage
COPY manager manager
COPY certificates certificates
COPY print print

RUN go build -o /app/server/bin ./cmd/web_server/main

FROM alpine:latest
LABEL maintainer="Andrii Yeremenko"

ENV PORT=8443

COPY --from=base /app/server/bin /app/bin
COPY --from=base /app/certificates certificates

EXPOSE ${PORT}

ENTRYPOINT ["app/bin"]