debloated_fs:
	mkdir -p build
	cd build && cmake ../fs && cmake --build .

baffs:
	mkdir -p build
	cd build && go build github.com/jzh18/baffs

install: debloated_fs baffs
	cp build/debloated_fs /usr/bin/debloated_fs
	cp build/baffs /usr/bin/baffs
	
test:
	go test -v ./...