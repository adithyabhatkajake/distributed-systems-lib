#!/bin/bash

EXP=d2-f1

./node/apollo/node -conf=testData/apollo/$EXP/nodes-0.txt &> /dev/null &
./node/apollo/node -conf=testData/apollo/$EXP/nodes-1.txt &> /dev/null &
./node/apollo/node -conf=testData/apollo/$EXP/nodes-2.txt &> /dev/null &

sleep 10
./client/apollo/client -conf=testData/apollo/$EXP/client.txt -batch=50 -metric=5 &> apollo-client-logs-$EXP.txt

killall ./node/apollo/node
killall ./client/apollo/client

EXP=d5-f1

./node/apollo/node -conf=testData/apollo/$EXP/nodes-0.txt &> /dev/null &
./node/apollo/node -conf=testData/apollo/$EXP/nodes-1.txt &> /dev/null &
./node/apollo/node -conf=testData/apollo/$EXP/nodes-2.txt &> /dev/null &

sleep 10
./client/apollo/client -conf=testData/apollo/$EXP/client.txt -batch=50 -metric=5 &> apollo-client-logs-$EXP.txt

killall ./node/apollo/node
killall ./client/apollo/client

EXP=d2-f2

./node/apollo/node -conf=testData/apollo/$EXP/nodes-0.txt &> /dev/null &
./node/apollo/node -conf=testData/apollo/$EXP/nodes-1.txt &> /dev/null &
./node/apollo/node -conf=testData/apollo/$EXP/nodes-2.txt &> /dev/null &
./node/apollo/node -conf=testData/apollo/$EXP/nodes-3.txt &> /dev/null &
./node/apollo/node -conf=testData/apollo/$EXP/nodes-4.txt &> /dev/null &

sleep 10
./client/apollo/client -conf=testData/apollo/$EXP/client.txt -batch=70 -metric=5 &> apollo-client-logs-$EXP.txt

killall ./node/apollo/node
killall ./client/apollo/client

EXP=d5-f2

./node/apollo/node -conf=testData/apollo/$EXP/nodes-0.txt &> /dev/null &
./node/apollo/node -conf=testData/apollo/$EXP/nodes-1.txt &> /dev/null &
./node/apollo/node -conf=testData/apollo/$EXP/nodes-2.txt &> /dev/null &
./node/apollo/node -conf=testData/apollo/$EXP/nodes-3.txt &> /dev/null &
./node/apollo/node -conf=testData/apollo/$EXP/nodes-4.txt &> /dev/null &

sleep 10
./client/apollo/client -conf=testData/apollo/$EXP/client.txt -batch=70 -metric=5 &> apollo-client-logs-$EXP.txt

killall ./node/apollo/node
killall ./client/apollo/client