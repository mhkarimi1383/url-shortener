#!/usr/bin/env bash

function shutdown {
    kill -s TERM $client_pid
    kill -s TERM $server_pid
}

trap "shutdown" SIGINT SIGTERM

cd client/url-shortener || exit
NUXT_APP_BASE_URL=/admin/ node .output/server/index.mjs &
client_pid=$!

cd ../.. || exit
go mod download
./tmp/main &
server_pid=$!

wait -n
