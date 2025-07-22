debloated_fs:
	mkdir -p build
	cd build && cmake ../fs && cmake --build .

baffs:
	CGO_ENABLED=0  go build -ldflags="-extldflags=-static" -o baffs


install: debloated_fs baffs
	cp build/debloated_fs /usr/bin/debloated_fs
	cp baffs /usr/bin/baffs
	
test: debloated_fs
	./build/all_test
	go test -v ./...

integration_test: install
	./tests/test.sh
