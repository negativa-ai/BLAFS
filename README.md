# BAFFS
sudo apt-get install libfuse3-dev
mkdir build
cd build
cmake ../fs
cmkae --build .


SPDLOG_LEVEL=debug ./debloated_fs -d ./mount