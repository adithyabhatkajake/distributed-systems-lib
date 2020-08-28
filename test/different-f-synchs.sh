#!/bin/bash

EXP=d2-f1

./node/synchs/node -conf=testData/synchs/$EXP/nodes-0.txt &> /dev/null &
./node/synchs/node -conf=testData/synchs/$EXP/nodes-1.txt &> /dev/null &
./node/synchs/node -conf=testData/synchs/$EXP/nodes-2.txt &> /dev/null &

sleep 10
./client/synchs/client -conf=testData/synchs/$EXP/client.txt -batch=30 -metric=5 &> client-logs-$EXP.txt

killall ./node/synchs/node
killall ./client/synchs/client

EXP=d5-f1

./node/synchs/node -conf=testData/synchs/$EXP/nodes-0.txt &> /dev/null &
./node/synchs/node -conf=testData/synchs/$EXP/nodes-1.txt &> /dev/null &
./node/synchs/node -conf=testData/synchs/$EXP/nodes-2.txt &> /dev/null &

sleep 10
./client/synchs/client -conf=testData/synchs/$EXP/client.txt -batch=30 -metric=5 &> client-logs-$EXP.txt

killall ./node/synchs/node
killall ./client/synchs/client

EXP=d2-f2

./node/synchs/node -conf=testData/synchs/$EXP/nodes-0.txt &> /dev/null &
./node/synchs/node -conf=testData/synchs/$EXP/nodes-1.txt &> /dev/null &
./node/synchs/node -conf=testData/synchs/$EXP/nodes-2.txt &> /dev/null &
./node/synchs/node -conf=testData/synchs/$EXP/nodes-3.txt &> /dev/null &
./node/synchs/node -conf=testData/synchs/$EXP/nodes-4.txt &> /dev/null &

sleep 10
./client/synchs/client -conf=testData/synchs/$EXP/client.txt -batch=30 -metric=5 &> client-logs-$EXP.txt

killall ./node/synchs/node
killall ./client/synchs/client

EXP=d5-f2

./node/synchs/node -conf=testData/synchs/$EXP/nodes-0.txt &> /dev/null &
./node/synchs/node -conf=testData/synchs/$EXP/nodes-1.txt &> /dev/null &
./node/synchs/node -conf=testData/synchs/$EXP/nodes-2.txt &> /dev/null &
./node/synchs/node -conf=testData/synchs/$EXP/nodes-3.txt &> /dev/null &
./node/synchs/node -conf=testData/synchs/$EXP/nodes-4.txt &> /dev/null &

sleep 10
./client/synchs/client -conf=testData/synchs/$EXP/client.txt -batch=30 -metric=5 &> client-logs-$EXP.txt

killall ./node/synchs/node
killall ./client/synchs/client