#!/bin/bash
i=0
port=80
# -lt is less than operator
docker network create myNetwork
while [ $i -lt 60 ]
    do docker run -d --name `expr $port`a -p `expr $port`:80 docker/getting-started 
    docker network connect myNetwork `expr $port`a
    i=`expr $i + 1`
    port=`expr $port + 1`
done
