#!/bin/bash
trap "rm server;kill 0" EXIT

go build -o server
./server -port=8001 thq_f=thq8 thq_s=thq9 &
./server -port=8002 thq_f=thq6 thq_s=thq7 &
./server -port=8003 -api=1 &

sleep 2
echo ">>> start test"
curl "http://localhost:9999/api?key=thq1" &
curl "http://localhost:9999/api?key=thq2" &
curl "http://localhost:9999/api?key=thq3" &
curl "http://localhost:9999/api?key=thq4" &
curl "http://localhost:9999/api?key=thq5" &
curl "http://localhost:9999/api?key=thq6" &
curl "http://localhost:9999/api?key=thq7" &
curl "http://localhost:9999/api?key=thq8" &
curl "http://localhost:9999/api?key=thq9" &

wait