#!/bin/bash

# check if baffs-dev already running
if [ "$(docker ps -q -f name=baffs-dev)" ]; then
    echo "baffs-dev already exists, connecting to it"
    docker exec -it baffs-dev /bin/zsh
    exit 0
fi

# check if baffs-dev already exists but stopped
if [ "$(docker ps -aq -f status=exited -f name=baffs-dev)" ]; then
    echo "baffs-dev already exists but stopped, starting it"
    docker start baffs-dev
    docker exec -it baffs-dev /bin/zsh
    exit 0
fi

# no need to mount .ssh if you don't need to push/pull from the repo
docker run  -d --name baffs-dev --privileged=true  -v /tmp/docker:/var/lib/docker  -v $PWD:/home/ubuntu/repos/BAFFS  baffs-dev
docker exec -it baffs-dev /bin/zsh
