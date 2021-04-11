FROM golang:1.16-alpine as builder
ENV CGO=0
ENV HOST="example.com"

RUN apk --no-cache add \
    make

WORKDIR /app

COPY go.mod go.mod
COPY go.sum go.sum
RUN --mount=type=cache,target=/root/go/pkg/mod go mod download

COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build make jammer

FROM alpine:3

COPY --from=builder /app/target/ /usr/local/bin
COPY ./scripts/app.sh /usr/local/bin

ARG USER=app
ENV HOME=/home/${USER}

RUN adduser ${USER} -D -s ash

USER ${USER}
WORKDIR ${HOME}

ENTRYPOINT ["/usr/local/bin/app.sh"]
