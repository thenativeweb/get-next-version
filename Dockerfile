FROM golang:1.18-alpine as build
RUN apk add --no-cache make gcc musl-dev
RUN mkdir /app
WORKDIR /app

ADD . .
RUN make build

FROM alpine

RUN apk add --no-cache git
RUN mkdir /action
WORKDIR /action

COPY --from=build /app/build/get-next-version /action/get-next-version
ADD ./action/entrypoint.sh /action/entrypoint.sh

ENTRYPOINT ["/action/entrypoint.sh"]
