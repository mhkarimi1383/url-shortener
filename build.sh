#!/usr/bin/env bash

function shutdown {
    kill -s TERM $client_pid
    kill -s TERM $server_pid
}

trap "shutdown" SIGINT SIGTERM

cd client/url-shortener || exit
npm install
npm run build
client_pid=$!

cd ../.. || exit
go mod download
go build -o ./tmp/main . &
server_pid=$!

wait -n
