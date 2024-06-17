FROM golang:latest as build

RUN apt-get update
RUN apt-get install ca-certificates make npm -y
RUN apt-get upgrade -y
RUN npm i -g pnpm

WORKDIR /app
# Set context to ../
COPY . .

#!/bin/sh
ARG GOPATH=/var/lib/jenkins/go
ENV GOPATH=$GOPATH
ENV PATH=$PATH:/usr/local/bin
ENV PATH=$PATH:$GOPATH/bin
ENV PATH=$PATH:/usr/local/go/bin

ARG BASE_URL=localhost:8000/api
ENV VITE_API_BASE_URL=http://$BASE_URL
ENV VITE_WS_BASE_URL=ws://$BASE_URL/ws

RUN rm -r cahfrontend/node_modules
RUN make frontend -j

FROM nginx:latest

WORKDIR /app
COPY --from=build /app/cahfrontend/dist dist/
COPY cahfrontend/nginx.conf /etc/nginx/nginx.conf
COPY devProxy/mime.types .

EXPOSE 80
