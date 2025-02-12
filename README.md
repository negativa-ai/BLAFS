![example workflow](https://github.com/jzh18/BAFFS/actions/workflows/main.yml/badge.svg)

# BAFFS
sudo apt-get install libfuse3-dev
mkdir build
cd build
cmake ../fs
cmkae --build .


SPDLOG_LEVEL=debug ./debloated_fs --realdir=/home/ubuntu/repos/BAFFS/build/real --lowerdir=/home/ubuntu/repos/BAFFS/build/lower --optimize=l0 -d /home/ubuntu/repos/BAFFS/build/mount
./debloated_fs  --realdir=/home/ubuntu/repos/BAFFS/build/real --lowerdir=/home/ubuntu/repos/BAFFS/build/lower --optimize=l0  /home/ubuntu/repos/BAFFS/build/mount

ls ./mount/aaa
Not support create file

 ./build/baffs shadow  --images=file_reader:latest
 ./build/baffs debloat  --images=file_reader:latest 
 must run with tags

./build/baffs shadow  --images=file_reader:layer_a,file_reader:layer_b
docker run -it --rm file_reader:layer_a f1.txt
docker run -it --rm file_reader:layer_b f2.txt
./build/baffs debloat  --images=file_reader:layer_a,file_reader:layer_b


./build/baffs shadow  --images=redis:7.4.1
docker run --rm -it --network host redis:7.4.1 
./build/baffs debloat  --images=redis:7.4.1