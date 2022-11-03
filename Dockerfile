# FROM ubuntu:latest as PROTO_BUILD
# WORKDIR /app
# # Install python
# RUN apt-get update 
# RUN apt-get install -y wget python3-pip unzip
# # RUN ln -s /usr/bin/python3 /usr/bin/python
# # RUN ln -s /usr/bin/pip3 /usr/bin/pip
# # RUN pip install grpc_web_proto_compile
# # Compile frontend protos
# ADD ./ /app
# # RUN /app/frontend/gen-protos.sh echo "Compiled protobuf files"

FROM golang:alpine AS GO_BUILD

WORKDIR /app
COPY ./ /app

ENV GOARCH amd64
ENV GOOS linux
RUN go build \
  -a \
  -ldflags="-w -s" \
  -o app \
  "github.com/romnn/ldap-manager/cmd/ldap-manager"

FROM node:latest AS NODE_BUILD

WORKDIR /app
COPY ./ /app
# COPY --from=PROTO_BUILD /app /app

# ENV SKIPPROTOCOMPILATION 1
RUN cd web && yarn install && yarn build

FROM gcr.io/distroless/static
LABEL maintainer="contact@romnn.com"

USER nonroot:nonroot

# romnn/distroless-base-grpc-health
# FROM romnn/distroless-base-grpc-health

COPY --from=GO_BUILD --chown=nonroot:nonroot /app/app /app
COPY --from=NODE_BUILD --chown=nonroot:nonroot /app/web/dist /web/dist

ENV HTTP_PORT 8080
ENV GRPC_PORT 9090

EXPOSE 8080
EXPOSE 9090

ENTRYPOINT [ "/app", "serve" ]
