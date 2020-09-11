FROM ubuntu:20.04 as PROTO_BUILD

LABEL maintainer="contact@romnn.com"

WORKDIR /app

# Install python
RUN apt-get update 
RUN apt-get install -y wget python3-pip unzip
RUN ln -s /usr/bin/python3 /usr/bin/python
RUN ln -s /usr/bin/pip3 /usr/bin/pip
RUN pip install grpc_web_proto_compile

# Compile frontend protos
ADD ./ /app
RUN /app/frontend/gen-protos.sh echo "Compiled protobuf files"

FROM golang:alpine AS GO_BUILD

ENV GO111MODULE=on

WORKDIR /app
COPY ./ /app

# This removes debug information from the binary
# Assumes go 1.10+
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -a -ldflags="-w -s" -o app "github.com/romnnn/ldap-manager/cmd/ldap-manager"

FROM node:latest AS NODE_BUILD

WORKDIR /app
COPY --from=PROTO_BUILD /app /app

ENV SKIPPROTOCOMPILATION 1
RUN cd frontend && npm install && npm rebuild node-sass && npm run build

FROM romnn/distroless-base-grpc-health

COPY --from=GO_BUILD /app/app /app
COPY --from=NODE_BUILD /app/frontend/dist /frontend/dist

ENV HTTP_PORT 80
ENV GRPC_PORT 9090

EXPOSE 80
EXPOSE 9090

ENTRYPOINT [ "/app", "serve" ]
