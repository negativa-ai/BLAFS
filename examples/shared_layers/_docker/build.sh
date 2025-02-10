#!/bin/bash

mkdir -p build
cd build
cmake ../..
make

cd ..
cp build/file_reader ./base
cd base

docker build -t file_reader:base .

cd ..
rm -rf build
rm file_reader

cd layer_a
docker build -t file_reader:layer_a .
cd ..

cd layer_b
docker build -t file_reader:layer_b .