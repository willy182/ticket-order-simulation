# Stage 1
FROM golang:1.26-alpine3.23 AS dependency_builder

WORKDIR /go/src
ENV GO111MODULE=on

RUN apk update
RUN apk add --no-cache bash ca-certificates git gcc musl-dev

COPY go.mod .
COPY go.sum .

RUN go mod download

# Stage 2
FROM dependency_builder AS service_builder

WORKDIR /usr/app

COPY . .
RUN go build -o bin

# Stage 3
FROM alpine:3.23

ARG BUILD_NUMBER
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /usr/app/
ENV BUILD_NUMBER=$BUILD_NUMBER

COPY --from=service_builder /usr/app/bin bin
COPY --from=service_builder /usr/app/.env .env

ENTRYPOINT ["./bin"]
