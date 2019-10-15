FROM golang as builder
WORKDIR /go/src/github.com/skaji/perl6-cpan-new
COPY go.* ./
RUN go mod download
COPY ./ ./
RUN cd cmd/perl6-cpan-new && go build

FROM alpine
RUN set -eux; \
  apk add --update --no-cache tzdata ca-certificates tini; \
  cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime; \
  echo Asia/Tokyo > /etc/timezone; \
  apk del tzdata; \
  :
COPY --from=builder /go/src/github.com/skaji/perl6-cpan-new/cmd/perl6-cpan-new/perl6-cpan-new /perl6-cpan-new
ENTRYPOINT ["/sbin/tini", "--"]
CMD ["/perl6-cpan-new"]
