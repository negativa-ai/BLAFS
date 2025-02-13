#!/bin/bash

docker build --no-cache -t justinzhf/baffs:latest .
docker push justinzhf/baffs:latest
