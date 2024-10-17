ARG GO_VERSION="1.22"
ARG NODE_MAJOR="20"

FROM golang:${GO_VERSION}-bookworm AS builder

ARG NODE_MAJOR

ENV CGO_ENABLED=0

RUN apt-get update && \
  apt-get install -y gnupg2 ca-certificates curl gnupg

RUN curl -fsSL https://deb.nodesource.com/gpgkey/nodesource-repo.gpg.key | gpg --dearmor -o /etc/apt/keyrings/nodesource.gpg
RUN echo "deb [signed-by=/etc/apt/keyrings/nodesource.gpg] https://deb.nodesource.com/node_${NODE_MAJOR}.x nodistro main" | \
  tee /etc/apt/sources.list.d/nodesource.list
RUN apt-get update && apt-get install nodejs -y
RUN npm install -g pnpm

WORKDIR /go/src/github.com/mhkarimi1383/url-shortener

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go generate ./...
RUN go build -x -o /url-shortener .

FROM gcr.io/distroless/static-debian12 AS runner

ENV USH_LISTEN_ADDRESS="0.0.0.0:8080"
ENV USH_DATABASE_CONNECTION_STRING="/data/database.sqlite3"

COPY --from=builder /url-shortener /opt/app/url-shortener

VOLUME [ "/data" ]

ENTRYPOINT [ "/opt/app/url-shortener" ]
