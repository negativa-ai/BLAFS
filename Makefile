debloated_fs:
	mkdir -p build
	cd build && cmake ../fs && cmake --build .

baffs:
	mkdir -p build
	cd build && go build -buildvcs=false github.com/jzh18/baffs

install: debloated_fs baffs
	cp build/debloated_fs /usr/bin/debloated_fs
	cp build/baffs /usr/bin/baffs
	
test: debloated_fs
	./build/all_test
	go test -v ./...

integration_test: install
	./tests/test.sh