FROM alpine
ADD zoom /bin/
RUN apk -Uuv add ca-certificates
ENTRYPOINT /bin/zoom