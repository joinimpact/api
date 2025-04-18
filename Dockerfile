FROM golang:1.14-buster AS builder

ENV GOROOT /usr/local/go

ENV GO111MODULE on

WORKDIR /src

COPY ./ ./

RUN set -xe \
          && make build


FROM alpine:latest

WORKDIR /app

RUN set -xe \
        && apk update && apk upgrade \
        && apk add --no-cache sudo \
	&& mkdir -p /etc/sudoers.d \
	&& echo "server ALL=(ALL) NOPASSWD: ALL" >> /etc/sudoers.d/server \
	&& addgroup -g 1001 server \
	&& adduser -D -s /bin/sh -u 1001 -G server server

EXPOSE 3000

USER 1001

COPY --from=builder --chown=1001:1001 /src/server ./
COPY --from=builder --chown=1001:1001 /src/docker/run/wait-for /wait-for
COPY --from=builder --chown=1001:1001 /src/docker/run/entrypoint.sh /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]

CMD ["./server"]
