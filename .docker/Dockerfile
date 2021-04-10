FROM golang:1.16-alpine as builder
ENV CGO=0

RUN apk --no-cache add \
    make

WORKDIR /app

COPY go.mod go.mod
COPY go.sum go.sum
RUN --mount=type=cache,target=/root/go/pkg/mod go mod download

COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build make jammer

FROM alpine:3

COPY --from=builder /app/target/ /usr/local/bin-

ENTRYPOINT [ "jammer" ]