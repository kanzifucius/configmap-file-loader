FROM alpine:latest

ADD bin/configmapfileloader_unix  /configmapfileloader
RUN chmod a+rw configmapfileloader
ENTRYPOINT ["./configmapfileloader"]
