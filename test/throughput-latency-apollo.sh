#!/bin/bash

./node/apollo/node -conf=testData/apollo/b5/nodes-0.txt &> /dev/null &
./node/apollo/node -conf=testData/apollo/b5/nodes-1.txt &> /dev/null &
./node/apollo/node -conf=testData/apollo/b5/nodes-2.txt &> /dev/null &

sleep 10
./client/apollo/client -conf=testData/apollo/b5/client.txt -batch=20 -metric=5 &> apollo-client-logs-1.txt

killall ./node/apollo/node
killall ./client/apollo/client

./node/apollo/node -conf=testData/apollo/b10/nodes-0.txt &> /dev/null &
./node/apollo/node -conf=testData/apollo/b10/nodes-1.txt &> /dev/null &
./node/apollo/node -conf=testData/apollo/b10/nodes-2.txt &> /dev/null &

sleep 10
./client/apollo/client -conf=testData/apollo/b10/client.txt -batch=30 -metric=5 &> apollo-client-logs-2.txt

killall ./node/apollo/node
killall ./client/apollo/client

./node/apollo/node -conf=testData/apollo/b20/nodes-0.txt &> /dev/null &
./node/apollo/node -conf=testData/apollo/b20/nodes-1.txt &> /dev/null &
./node/apollo/node -conf=testData/apollo/b20/nodes-2.txt &> /dev/null &

sleep 10
./client/apollo/client -conf=testData/apollo/b20/client.txt -batch=50 -metric=5 &> apollo-client-logs-3.txt

killall ./node/apollo/node
killall ./client/apollo/client
