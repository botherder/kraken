BUILD_FOLDER  = $(CURDIR)/build

FLAGS_LINUX   = GOOS=linux GOARCH=amd64 CGO_ENABLED=1
FLAGS_DARWIN  = GOOS=darwin GOARCH=amd64 CGO_ENABLED=1
FLAGS_FREEBSD = GOOS=freebsd GOARCH=amd64 CGO_ENABLED=1
FLAGS_WINDOWS_386 = GOOS=windows GOARCH=386 CC=i686-w64-mingw32-gcc CGO_ENABLED=1
FLAGS_WINDOWS_AMD64 = GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc CGO_ENABLED=1

KRAKEN_SRC   = $(wildcard *.go)
LAUNCHER_SRC = $(wildcard launcher/*.go)
COMPILER_SRC = $(wildcard compiler/*.go)

.PHONY: all
all: $(shell go env GOOS)

.PHONY: lint
lint:
	@echo "[lint] Running linter on codebase"
	@golint ./...

.PHONY: deps
deps:
	@echo "[deps] Installing dependencies..."
	go mod download
	go get github.com/akavel/rsrc
	@echo "[deps] Dependencies installed."

.PHONY: check-env
check-env:
	@mkdir -p $(BUILD_FOLDER)

ifndef RULES
	@echo "[check-env] You have not specified any RULES env, kraken will not have any default Yara rules."
endif

ifndef BACKEND
	@echo "[check-env] You have not specified any BACKEND env, kraken will not have any default backend server configured."
endif


.PHONY: rules-compiler
rules-compiler: $(BUILD_FOLDER)/compiler
$(BUILD_FOLDER)/compiler: $(COMPILER_SRC)
ifdef RULES
	@mkdir -p $(@D)

	@echo "[rules-compiler] Building rules compiler..."
	@cd compiler; go build -o $@

	@echo "[rules-compiler] Compiling Yara rules..."
	@$@ $(RULES)

	@echo "[rules-compiler] Launching binary resource builder..."
	@go-bindata rules
endif


.PHONY: linux
linux: $(BUILD_FOLDER)/linux/kraken $(BUILD_FOLDER)/linux/kraken-launcher

$(BUILD_FOLDER)/linux/kraken: $(BUILD_FOLDER)/compiler $(KRAKEN_SRC)
	@mkdir -p $(@D)
	@echo "[builder] Building Linux executable..."
	@$(FLAGS_LINUX) go build --ldflags '-s -w -extldflags "-lm -static" -X main.DefaultBaseDomain=$(BACKEND)' \
		-tags yara_static -o $@
	@echo "[builder] Done!"

$(BUILD_FOLDER)/linux/kraken-launcher: $(LAUNCHER_SRC)
	@mkdir -p $(@D)
	@echo "[builder] Building Linux launcher..."
	@cd launcher; $(FLAGS_LINUX) go build --ldflags '-s -w' \
		-o $@
	@echo "[builder] Done!"

.PHONY: darwin
darwin: $(BUILD_FOLDER)/darwin/kraken $(BUILD_FOLDER)/darwin/kraken-launcher

$(BUILD_FOLDER)/darwin/kraken: $(BUILD_FOLDER)/compiler $(KRAKEN_SRC)
	@mkdir -p $(@D)
	@echo "[builder] Building Darwin executable..."
	@$(FLAGS_DARWIN) go build --ldflags '-s -w -extldflags "-lm" -X main.DefaultBaseDomain=$(BACKEND)' \
		-tags yara_static -o $@
	@echo "[builder] Done!"

$(BUILD_FOLDER)/darwin/kraken-launcher: $(LAUNCHER_SRC)
	@mkdir -p $(@D)
	@echo "[builder] Building Darwin launcher..."
	@cd launcher; $(FLAGS_DARWIN) go build --ldflags '-s -w' \
		-o $@
	@echo "[builder] Done!"

.PHONY: freebsd
freebsd: $(BUILD_FOLDER)/freebsd/kraken $(BUILD_FOLDER)/freebsd/kraken-launcher

$(BUILD_FOLDER)/freebsd/kraken: $(BUILD_FOLDER)/compiler $(KRAKEN_SRC)
	@mkdir -p $(@D)
	@echo "[builder] Building FreeBSD executable..."
	@$(FLAGS_FREEBSD) go build --ldflags '-s -w -extldflags "-lm -static" -X main.DefaultBaseDomain=$(BACKEND)' \
		-tags yara_static -o $@
	@echo "[builder] Done!"

$(BUILD_FOLDER)/freebsd/kraken-launcher: $(LAUNCHER_SRC)
	@mkdir -p $(@D)
	@echo "[builder] Building FreeBSD launcher..."
	@cd launcher; $(FLAGS_FREEBSD) go build --ldflags '-s -w' }
		-o $@
	@echo "[builder] Done!"

.PHONY: windows
windows: windows-386 windows-amd64

.PHONY: windows-386
windows-386: $(BUILD_FOLDER)/windows-386/kraken.exe $(BUILD_FOLDER)/windows-386/kraken-launcher.exe

$(BUILD_FOLDER)/windows-386/kraken.exe: $(BUILD_FOLDER)/compiler $(YARA_SRC)-windows-386/done $(KRAKEN_SRC)
	@mkdir -p $(@D)

	#@rsrc -manifest kraken.manifest -ico kraken.ico -o rsrc.syso
	@rsrc -arch 386 -manifest kraken.manifest -o rsrc_windows_386.syso

	@echo "[builder] Building Windows 32bit executable..."
	@$(FLAGS_WINDOWS_386) go build --ldflags '-s -w -extldflags "-static" -X main.DefaultBaseDomain=$(BACKEND)' \
		-tags yara_static -o $@
	@echo "[builder] Done!"

$(BUILD_FOLDER)/windows-386/kraken-launcher.exe: $(LAUNCHER_SRC)
	@mkdir -p $(@D)
	@echo "[builder] Building Windows 32bit launcher..."
	@cd launcher; $(FLAGS_WINDOWS_386) go build --ldflags '-s -w -extldflags "-static" -H=windowsgui' \
		-o $@
	@echo "[builder] Done!"

.PHONY: windows-amd64
windows-amd64: $(BUILD_FOLDER)/windows-amd64/kraken.exe $(BUILD_FOLDER)/windows-amd64/kraken-launcher.exe

$(BUILD_FOLDER)/windows-amd64/kraken.exe: $(BUILD_FOLDER)/compiler $(YARA_SRC)-windows-amd64/done $(KRAKEN_SRC)
	@mkdir -p $(@D)

	#@rsrc -manifest kraken.manifest -ico kraken.ico -o rsrc.syso
	@rsrc -arch amd64 -manifest kraken.manifest -o rsrc_windows_amd64.syso

	@echo "[builder] Building Windows 64bit executable..."
	@$(FLAGS_WINDOWS_AMD64) go build --ldflags '-s -w -extldflags "-static" -X main.DefaultBaseDomain=$(BACKEND)' \
		-tags yara_static -o $@
	@echo "[builder] Done!"

$(BUILD_FOLDER)/windows-amd64/kraken-launcher.exe: $(LAUNCHER_SRC)
	@mkdir -p $(@D)
	@echo "[builder] Building Windows 64bit launcher..."
	@cd launcher; $(FLAGS_WINDOWS_AMD64) go build --ldflags '-s -w -extldflags "-static" -H=windowsgui' \
		-o $@
	@echo "[builder] Done!"

.PHONY: clean
clean:
	rm -f rules
	rm -f bindata.go
	rm -f rsrc.syso
	rm -rf $(BUILD_FOLDER)
