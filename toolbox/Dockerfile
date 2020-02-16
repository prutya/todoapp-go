FROM ruby:2.7.0-alpine3.11 AS build

RUN apk add --update --no-cache \
  make=4.2.1-r2 \
  g++=9.2.0-r3 \
  postgresql-dev=12.1-r0 \
  postgresql-client=12.1-r0 \
  zsh=5.7.1-r0

WORKDIR /usr/src/app
