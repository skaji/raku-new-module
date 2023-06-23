ARG GO_VERSION
FROM golang:$GO_VERSION-alpine as builder
RUN apk add --update --no-cache git
WORKDIR /app
COPY ./ ./
RUN cd cmd/raku-new-module && go build -buildvcs=false

FROM alpine
LABEL org.opencontainers.image.source https://github.com/skaji/raku-new-module
RUN set -eux; \
  apk add --update --no-cache tzdata ca-certificates tini; \
  cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime; \
  echo Asia/Tokyo > /etc/timezone; \
  apk del tzdata; \
  :
COPY --from=builder /app/cmd/raku-new-module/raku-new-module /raku-new-module
ENTRYPOINT ["/sbin/tini", "--"]
CMD ["/raku-new-module", "-config-from-env"]
