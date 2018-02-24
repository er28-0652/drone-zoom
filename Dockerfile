FROM alpine
RUN apk -Uuv add ca-certificates
ADD zoom /bin/
ENTRYPOINT [ "/bin/zoom" ]