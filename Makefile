debloated_fs:
	mkdir -p build
	cd build && cmake ../fs && cmake --build .

baffs:
	mkdir -p build
	cd build && go build github.com/jzh18/baffs
