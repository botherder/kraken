BUILD_FOLDER  = $(CURDIR)/build

FLAGS_LINUX   = GOOS=linux GOARCH=amd64 CGO_ENABLED=1
FLAGS_DARWIN  = GOOS=darwin GOARCH=amd64 CGO_ENABLED=1
FLAGS_FREEBSD = GOOS=freebsd GOARCH=amd64 CGO_ENABLED=1
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
	@mkdir -p $(BUILD_FOLDER)

ifndef RULES
	@echo "[check-env] You have not specified any RULES env, kraken will not have any default Yara rules."
endif

ifndef BACKEND
	@echo "[check-env] You have not specified any BACKEND env, kraken will not have any default backend server configured."
endif


rules-compiler:
ifdef RULES
	@echo "[rules-compiler] Building rules compiler..."
	@cd compiler; go build -o $(BUILD_FOLDER)/compiler; cd ..

	@echo "[rules-compiler] Compiling Yara rules..."
	@$(BUILD_FOLDER)/compiler $(RULES)

	@echo "[rules-compiler] Launching binary resource builder..."
	@go-bindata rules
endif


linux: check-env rules-compiler
	@mkdir -p $(BUILD_FOLDER)/linux

	@echo "[builder] Building Linux executable..."
	@$(FLAGS_LINUX) go build --ldflags '-s -w -extldflags "-lm -static" -X main.DefaultBaseDomain=$(BACKEND)' \
		-tags yara_static -o $(BUILD_FOLDER)/linux/kraken

	# @echo "[builder] Building launcher..."
	# @cd launcher; $(FLAGS_LINUX) go build --ldflags '-s -w' \
	# 	-o $(BUILD_FOLDER)/linux/kraken-launcher; cd ..

	@echo "[builder] Done!"


darwin: check-env rules-compiler
	@mkdir -p $(BUILD_FOLDER)/darwin

	@echo "[builder] Building Darwin executable..."
	@$(FLAGS_DARWIN) go build --ldflags '-s -w -extldflags "-lm" -X main.DefaultBaseDomain=$(BACKEND)' \
		-tags yara_static -o $(BUILD_FOLDER)/darwin/kraken

	# @echo "[builder] Building launcher..."
	# @cd launcher; $(FLAGS_DARWIN) go build --ldflags '-s -w' \
	# 	-o $(BUILD_FOLDER)/darwin/kraken-launcher; cd ..

	@echo "[builder] Done!"


freebsd: check-env rules-compiler
	@mkdir -p $(BUILD_FOLDER)/freebsd

	@echo "[builder] Building FreeBSD executable..."
	@$(FLAGS_FREEBSD) go build --ldflags '-s -w -extldflags "-lm -static" -X main.DefaultBaseDomain=$(BACKEND)' \
		-tags yara_static -o $(BUILD_FOLDER)/freebsd/kraken

	# @echo "[builder] Building launcher..."
	# @cd launcher; $(FLAGS_FREEBSD) go build --ldflags '-s -w' -o $(BUILD_FOLDER)/freebsd/kraken-launcher; cd ..

	@echo "[builder] Done!"


windows: check-env rules-compiler
	@mkdir -p $(BUILD_FOLDER)/windows

	#@rsrc -manifest kraken.manifest -ico kraken.ico -o rsrc.syso
	@rsrc -manifest kraken.manifest -o rsrc.syso

	@echo "[builder] Building Windows executable..."
	@$(FLAGS_WINDOWS) go build --ldflags '-s -w -extldflags "-static" -X main.DefaultBaseDomain=$(BACKEND)' \
		-tags yara_static -o $(BUILD_FOLDER)/windows/kraken.exe

	# @echo "[builder] Building launcher..."
	# @cd launcher; $(FLAGS_WINDOWS) go build --ldflags '-s -w -extldflags "-static" -H=windowsgui' \
	# 	-o $(BUILD_FOLDER)/windows/kraken-launcher.exe; cd ..

	@echo "[builder] Done!"


clean:
	rm -f rules
	rm -f bindata.go
	rm -f rsrc.syso
	rm -rf $(BUILD_FOLDER)
