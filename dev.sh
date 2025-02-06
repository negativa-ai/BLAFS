#!/bin/bash

docker run  -d --rm --name baffs-dev --privileged=true  -v /tmp/docker:/var/lib/docker   baffs-dev
docker exec -it baffs-dev /bin/zsh
