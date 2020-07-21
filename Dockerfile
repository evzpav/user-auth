# --- Base ----
FROM golang:1.14-stretch AS base
WORKDIR $GOPATH/src/gitlab.com/user-auth

# ---- Dependencies ----
FROM base AS dependencies
ENV GO111MODULE=on
COPY . .
RUN ls -l

# ---- Build ----
FROM dependencies AS build
ARG VERSION
ARG BUILD
ARG DATE
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -a -installsuffix cgo -o /go/bin/user-auth -ldflags "-X main.version=${VERSION} -X main.build=${BUILD} -X main.date=${DATE}" ./cmd/server/main.go
ENTRYPOINT ["/go/bin/user-auth"]
