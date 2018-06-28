BUILD_FOLDER	= $(shell pwd)/build
FLAGS_WINDOWS	= GOOS=windows GOARCH=386 CC=i686-w64-mingw32-gcc CGO_ENABLED=1

lint:
	@echo "[lint] Running linter on codebase"
	@golint ./...

deps:
	@echo "[deps] Installing dependencies..."
	$(FLAGS_WINDOWS) go get -u --ldflags '-s -w -extldflags "-static" -H windowsgui' github.com/lxn/win
	$(FLAGS_WINDOWS) go get -u --ldflags '-s -w -extldflags "-static" -H windowsgui' github.com/lxn/walk
	go get github.com/akavel/rsrc
	go get golang.org/x/sys/windows/registry
	go get gopkg.in/Knetic/govaluate.v3
	@echo "[deps] Dependencies installed."

build:
	@echo "[builder] Building Windows executable"
	@mkdir -p $(BUILD_FOLDER)/

	$(GOPATH)/bin/rsrc -manifest harden.manifest -ico harden.ico -o rsrc.syso
	$(FLAGS_WINDOWS) go build --ldflags '-s -w -extldflags "-static" -H windowsgui' -o $(BUILD_FOLDER)/hardentools.exe

	@echo "[builder] Done!"

clean:
	rm -f rsrc.syso
	rm -rf $(BUILD_FOLDER)
