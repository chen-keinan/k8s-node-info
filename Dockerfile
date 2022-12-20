FROM alpine:3.16.2

RUN adduser -u 10000 -D -g '' nodeinfo nodeinfo

COPY node-info /usr/local/bin/node-info

USER 10000

ENTRYPOINT ["node-info"]