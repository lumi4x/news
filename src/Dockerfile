FROM golang:1.26.1-alpine3.19 AS builder

ENV CGO_ENABLED=0

# Install git for fetching the dependencies
RUN apk update && apk add --no-cache git && apk add --no-cache bash && apk add build-base nodejs npm protobuf-dev

# Setup folders
RUN mkdir /app
WORKDIR /app

COPY go.sum ./
COPY go.mod ./
RUN go mod download

# Copy the source from the current directory to the working Directory inside the container
COPY ./src/common ./src/common
COPY ./scripts ./scripts

# Clients are available to all services
# HACK: RUN true is a workaround for a docker bug https://github.com/moby/moby/issues/37965

COPY ./src/services/auth/client ./src/services/auth/client
RUN true

COPY ./src/services/fileserver/client ./src/services/fileserver/client
RUN true

COPY ./src/services/app/client ./src/services/app/client
RUN true

ARG SERVICE
COPY ./src/services/${SERVICE} ./src/services/${SERVICE}

# Go generate
RUN go generate ./...

# Might resolve an infinite loop bug
RUN true

# Download all the dependencies
#RUN --mount=type=cache,target=/root/.cache/go-build go get -d -v ./...

# Install the package
#RUN --mount=type=cache,target=/root/.cache/go-build go install -v ./...

# Build the Go app
RUN --mount=type=cache,target=/root/.cache/go-build go build -o /build ./src/services/${SERVICE}

FROM alpine:latest AS goapp
WORKDIR /root/
COPY --from=builder /build .
CMD ["./build"]

FROM goapp AS goappapi
ARG PORT
EXPOSE ${PORT}
