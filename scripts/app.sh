#!/usr/bin/env sh

jammer "$@" \
    -addr='localhost:1986' \
    -content-dir="$HOME/blog" \
    -host="${HOST}" \
    -keyfile='/etc/jammer/key.pem' \
    -certfile='/etc/jammer/cert.pem'
