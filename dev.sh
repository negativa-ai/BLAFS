#!/bin/bash

docker run  -d --rm --name baffs-dev --privileged=true  -v /tmp/docker:/var/lib/docker  -v $PWD:/home/ubuntu/repos/BAFFS baffs-dev
docker exec -it baffs-dev /bin/zsh
