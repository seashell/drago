FROM alpine:3.10

RUN apk add -U wireguard-tools

COPY ./build/linux_amd64/drago ./drago

ENTRYPOINT ./drago agent --config=./