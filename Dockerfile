FROM golang:alpine AS builder

WORKDIR /src

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY cmd cmd
COPY internal internal

RUN go build ./cmd/plaza

FROM alpine:3.15

RUN apk add ffmpeg
COPY --from=builder /src/plaza .

ENTRYPOINT ["/plaza"]
