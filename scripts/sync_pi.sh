#!/bin/sh

[ -z "$RPI_HOST" ] && echo "RPI_HOST env variable is not set. Exiting..." && exit 1;
[ -z "$RPI_PORT" ] && echo "RPI_PORT env variable is not set. Exiting..." && exit 1;

remote_path="/home/pi/go/src/github.com/brewm/gobrewmmer"

fswatch -0 -o -e .git/ . | \
xargs -0 -I {} rsync -avz -e "ssh -p $RPI_PORT " . pi@$RPI_HOST:$remote_path
