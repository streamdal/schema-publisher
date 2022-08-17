ARG ALPINE_VERSION=3.13

FROM golang:1.16-alpine$ALPINE_VERSION AS builder

RUN apk --update add make bash curl git

WORKDIR /

COPY . .

RUN make build/linux

FROM library/alpine:$ALPINE_VERSION

RUN apk --update add bash ca-certificates && update-ca-certificates

COPY entrypoint.sh /entrypoint.sh
COPY --from=builder /publish-linux /publish

RUN chmod +x /publish

ENTRYPOINT ["/entrypoint.sh"]