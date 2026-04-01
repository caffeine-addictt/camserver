FROM golang:1.26.0-alpine AS builder
WORKDIR /app

COPY go* .
RUN go mod download

COPY pkg pkg
COPY internal internal
COPY cmd cmd
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o camserver ./cmd/daemon/

FROM builder AS prod
ARG UID=1000
ARG GID=1000
WORKDIR /

RUN apk add --no-cache ffmpeg ca-certificates dumb-init
RUN addgroup -g $GID cam && adduser -D -u $UID -G cam cam

COPY --from=builder /app/camserver /camserver

USER cam
ENTRYPOINT [ "dumb-init", "/camserver", "-v" ]
