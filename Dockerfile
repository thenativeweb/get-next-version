FROM golang:1.18.3-alpine as build

RUN mkdir /app
WORKDIR /app

RUN apk update
RUN apk add --no-cache build-base
RUN apk add --no-cache make

ADD . .
RUN make build



FROM alpine:3.18.4

RUN apk update 
RUN apk add --no-cache git

RUN mkdir /action
WORKDIR /action

COPY --from=build /app/build/get-next-version ./get-next-version
ADD ./action/entrypoint.sh ./entrypoint.sh

ENTRYPOINT ["/action/entrypoint.sh"]
