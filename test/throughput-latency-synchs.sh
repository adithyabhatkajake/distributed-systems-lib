#!/bin/bash

./node/synchs/node -conf=testData/synchs/b5/nodes-0.txt &> /dev/null &
./node/synchs/node -conf=testData/synchs/b5/nodes-1.txt &> /dev/null &
./node/synchs/node -conf=testData/synchs/b5/nodes-2.txt &> /dev/null &

sleep 10
./client/synchs/client -conf=testData/synchs/b5/client.txt -batch=10 -metric=5 &> client-logs-1.txt

killall ./node/synchs/node
killall ./client/synchs/client

./node/synchs/node -conf=testData/synchs/b10/nodes-0.txt &> /dev/null &
./node/synchs/node -conf=testData/synchs/b10/nodes-1.txt &> /dev/null &
./node/synchs/node -conf=testData/synchs/b10/nodes-2.txt &> /dev/null &

sleep 10
./client/synchs/client -conf=testData/synchs/b10/client.txt -batch=20 -metric=5 &> client-logs-2.txt

killall ./node/synchs/node
killall ./client/synchs/client

./node/synchs/node -conf=testData/synchs/b20/nodes-0.txt &> /dev/null &
./node/synchs/node -conf=testData/synchs/b20/nodes-1.txt &> /dev/null &
./node/synchs/node -conf=testData/synchs/b20/nodes-2.txt &> /dev/null &

sleep 10
./client/synchs/client -conf=testData/synchs/b20/client.txt -batch=30 -metric=5 &> client-logs-3.txt

killall ./node/synchs/node
killall ./client/synchs/client
