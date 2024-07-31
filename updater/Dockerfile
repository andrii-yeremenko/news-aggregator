FROM golang:1.22-alpine AS base
LABEL maintainer="Andrii Yeremenko"
WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download
RUN apk --no-cache add ca-certificates

COPY main.go ./
COPY updater/ ./updater/
COPY storage/ ./storage/
COPY model/ ./model/

RUN go build -o news-updater main.go

FROM scratch
LABEL maintainer="Andrii Yeremenko"

COPY --from=base /app/news-updater news-updater
COPY --from=base /etc/ssl/certs /etc/ssl/certs

VOLUME ["/resources", "/config"]

ENTRYPOINT ["./news-updater"]