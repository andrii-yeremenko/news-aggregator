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
COPY print print

RUN go build -o /app/server/bin ./cmd/web_server/main

FROM scratch
LABEL maintainer="Andrii Yeremenko"

ENV PORT=8443
ENV TIMEOUT=12h

COPY --from=base /app/server/bin /app/bin

VOLUME ["/resources", "/config"]

EXPOSE ${PORT}

ENTRYPOINT ["app/bin"]