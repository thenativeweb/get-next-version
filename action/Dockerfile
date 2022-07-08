FROM golang:1.18.3-bullseye as build

RUN mkdir /app
WORKDIR /app

ADD .. .
RUN make build

FROM debian:bullseye

RUN mkdir /action
WORKDIR /action

COPY --from=build /app/build/get-next-version ./get-next-version
ADD ./entrypoint.sh ./entrypoint.sh

ENTRYPOINT ["/action/entrypoint.sh"]