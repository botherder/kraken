BUILD_FOLDER  = $(shell pwd)/build

FLAGS_LINUX   = GOOS=linux GOARCH=amd64
FLAGS_DARWIN  = GOOS=darwin
FLAGS_WINDOWS = GOOS=windows GOARCH=386 CC=i686-w64-mingw32-gcc CGO_ENABLED=1


lint:
	@echo "[lint] Running linter on codebase"
	@golint ./...


deps:
	@echo "[deps] Installing dependencies..."
	go mod download
	go get github.com/akavel/rsrc
	@echo "[deps] Dependencies installed."


check-env:
ifndef RULES
	$(error You need to specify a RULES file or folder)
endif

ifndef BACKEND
	$(error You need to specify a BACKEND domain name)
endif


linux: check-env
	@mkdir -p $(BUILD_FOLDER)/linux

	@echo "[builder] Building rules compiler"
	@cd compiler; go build -o $(BUILD_FOLDER)/compiler; cd ..

	@echo "[builder] Compiling Yara rules..."
	@$(BUILD_FOLDER)/compiler $(RULES)

	@echo "[builder] Launching binary resource builder..."
	@go-bindata rules

	@echo "[builder] Building Linux executable..."
	@go build --ldflags '-s -w -X main.DefaultBaseDomain=$(BACKEND)' -tags yara_static -o $(BUILD_FOLDER)/linux/kraken

	@echo "[builder] Building launcher..."
	@cd launcher; go build --ldflags '-s -w' -o $(BUILD_FOLDER)/linux/kraken-launcher; cd ..

	@echo "[builder] Done!"


darwin: check-env
	@mkdir -p $(BUILD_FOLDER)/darwin

	@echo "[builder] Building rules compiler"
	@cd compiler; go build -o $(BUILD_FOLDER)/compiler; cd ..

	@echo "[builder] Compiling Yara rules..."
	@$(BUILD_FOLDER)/compiler $(RULES)

	@echo "[builder] Launching binary resource builder..."
	@go-bindata rules

	@echo "[builder] Building Darwin executable..."
	@go build --ldflags '-s -w -X main.DefaultBaseDomain=$(BACKEND)' -o $(BUILD_FOLDER)/darwin/kraken

	@echo "[builder] Building launcher..."
	@cd launcher; go build --ldflags '-s -w' -o $(BUILD_FOLDER)/darwin/kraken-launcher; cd ..

	@echo "[builder] Done!"


windows: check-env
	@mkdir -p $(BUILD_FOLDER)/windows

	@echo "[builder] Building rules compiler"
	@cd compiler; go build -o $(BUILD_FOLDER)/compiler; cd ..

	@echo "[builder] Compiling Yara rules..."
	@$(BUILD_FOLDER)/compiler $(RULES)

	@echo "[builder] Launching binary resource builder..."
	@go-bindata rules

	#@rsrc -manifest kraken.manifest -ico kraken.ico -o rsrc.syso
	@rsrc -manifest kraken.manifest -o rsrc.syso

	@echo "[builder] Building Windows executable..."
	@$(FLAGS_WINDOWS) go build --ldflags '-s -w -extldflags "-static" -X main.DefaultBaseDomain=$(BACKEND)' -o $(BUILD_FOLDER)/windows/kraken.exe

	@echo "[builder] Building launcher..."
	@cd launcher; $(FLAGS_WINDOWS) go build --ldflags '-s -w -extldflags "-static" -H=windowsgui' -o $(BUILD_FOLDER)/windows/kraken-launcher.exe; cd ..

	@echo "[builder] Done!"


clean:
	rm -f rules
	rm -f bindata.go
	rm -f rsrc.syso
	rm -rf $(BUILD_FOLDER)
