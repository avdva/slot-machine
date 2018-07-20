#!/bin/bash

for i in {1..1000}
do
    curl -X POST --data "@payload.dat" 127.0.0.1:1313/api/machines/atkins-diet/spins 1>/dev/null 2>/dev/null
done