FROM golang:alpine AS GO_BUILD

LABEL maintainer="contact@romnn.com"

ENV GO111MODULE=on

WORKDIR /app
COPY ./ /app

# This removes debug information from the binary
# Assumes go 1.10+
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -a -ldflags="-w -s" -o app "github.com/romnnn/ldap-manager/cmd/ldap-manager"

FROM node:latest AS NODE_BUILD

WORKDIR /app
COPY ./ /app

RUN cd frontend && npm install && npm rebuild node-sass && npm run build

FROM gcr.io/distroless/static
COPY --from=GO_BUILD /app/app /app
COPY --from=NODE_BUILD /app/frontend/dist /frontend/dist
ENV PORT 3000
EXPOSE 3000
ENTRYPOINT [ "/app" ]
