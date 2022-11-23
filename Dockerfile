# GO_BUILD
FROM golang:alpine AS GO_BUILD

WORKDIR /app
COPY ./ /app

ENV CGO_ENABLED 0
ENV GOARCH amd64
ENV GOOS linux
RUN go build \
  -a \
  -ldflags="-w -s" \
  -o app \
  "github.com/romnn/ldap-manager/cmd/ldap-manager"

# NODE_BUILD
FROM node:18 AS NODE_BUILD

WORKDIR /app
COPY ./ /app

RUN cd web && yarn install && yarn build

# FINAL
FROM gcr.io/distroless/static

LABEL maintainer="contact@romnn.com"

USER nonroot:nonroot

COPY --from=GO_BUILD --chown=nonroot:nonroot /app/app /app
COPY --from=NODE_BUILD --chown=nonroot:nonroot /app/web/dist /web/dist

ENV STATIC_ROOT /web/dist
ENV GRPC_PORT 9090
ENV GRPC_PORT 9090

EXPOSE 8080
EXPOSE 9090

ENTRYPOINT [ "/app", "serve" ]
