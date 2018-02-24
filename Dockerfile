FROM alpine:3.7 as alpine
ADD zoom /bin/
RUN apk -Uuv add ca-certificates
ENTRYPOINT /bin/zoom