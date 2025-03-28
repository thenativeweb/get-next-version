FROM golang:1.24.1-alpine3.21 as build

RUN mkdir /app
WORKDIR /app

RUN apk update && \
    apk add --no-cache build-base git make

ADD . .
RUN make build



FROM alpine:3.21.3

RUN apk update && \
    apk add --no-cache git

RUN mkdir /action
WORKDIR /action

COPY --from=build /app/build/get-next-version ./get-next-version
ADD ./action/entrypoint.sh ./entrypoint.sh

ENTRYPOINT ["/action/entrypoint.sh"]
