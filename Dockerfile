FROM golang:1.17-alpine3.14 as builder
WORKDIR /go/src/github.com/skaji/raku-new-module
COPY go.* ./
RUN apk add --update --no-cache git
RUN go mod download
COPY ./ ./
RUN cd cmd/raku-new-module && go build

FROM alpine:3.14
LABEL org.opencontainers.image.source https://github.com/skaji/raku-new-module
RUN set -eux; \
  apk add --update --no-cache tzdata ca-certificates tini; \
  cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime; \
  echo Asia/Tokyo > /etc/timezone; \
  apk del tzdata; \
  :
COPY --from=builder /go/src/github.com/skaji/raku-new-module/cmd/raku-new-module/raku-new-module /raku-new-module
ENTRYPOINT ["/sbin/tini", "--"]
