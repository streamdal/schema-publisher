FROM alpine:3.10

COPY entrypoint.sh /entrypoint.sh
COPY publish-linux /publish

RUN apk --update add bash ca-certificates && update-ca-certificates

RUN chmod +x /publish

ENTRYPOINT ["/entrypoint.sh"]
