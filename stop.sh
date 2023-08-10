#!/usr/bin/env bash

PID=`cat smpp-gateway.pid`
# echo $PID
kill -15 $PID
rm smpp-gateway.pid
