#!/bin/bash

# a script to get slot-machine flamegraph.
# it requires flamegraph scripts in the PATH or current directory.

go get github.com/uber/go-torch
go get github.com/tsenart/vegeta

killall slot-machine 2>/dev/null

sh -c "cd .. && ./slot-machine" &
sh -c "sleep 1 && bash attack.bash > /dev/null" &
sleep 1 && go-torch -u http://127.0.0.1:6060/debug/pprof -p -t=15 > atkins.svg

killall slot-machine
