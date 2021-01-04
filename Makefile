BUILD_FOLDER  = $(CURDIR)/build

FLAGS_LINUX   = GOOS=linux GOARCH=amd64 CGO_ENABLED=1 CGO_CFLAGS="-g -O2 -Wno-return-local-addr"
FLAGS_DARWIN  = GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 CGO_CFLAGS="-g -O2 -Wno-return-local-addr"
FLAGS_FREEBSD = GOOS=freebsd GOARCH=amd64 CGO_ENABLED=1 CGO_CFLAGS="-g -O2 -Wno-return-local-addr"
FLAGS_WINDOWS_386 = GOOS=windows GOARCH=386 CC=i686-w64-mingw32-gcc CGO_ENABLED=1 PKG_CONFIG_PATH=$(CURDIR)/_non-golang/prefix-windows-386/lib/pkgconfig/ CGO_CFLAGS="-g -O2 -Wno-return-local-addr"
FLAGS_WINDOWS_AMD64 = GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc CGO_ENABLED=1 PKG_CONFIG_PATH=$(CURDIR)/_non-golang/prefix-windows-amd64/lib/pkgconfig/ CGO_CFLAGS="-g -O2 -Wno-return-local-addr"

AGENT_VERSION = $(shell git describe --tags)

KRAKEN_SRC   = $(wildcard cmd/kraken/*.go)
LAUNCHER_SRC = $(wildcard cmd/kraken-launcher/*.go)
COMPILER_SRC = $(wildcard cmd/rules-compiler/*.go)

YARA_VERSION ?= 4.0.2
YARA_URL = https://github.com/VirusTotal/yara/archive/v$(YARA_VERSION).tar.gz
YARA_SRC = $(CURDIR)/_non-golang/yara-$(YARA_VERSION)

.PHONY: all
all: $(shell go env GOOS)

#==============================================================================
# Yara-related builds
#==============================================================================
$(YARA_SRC).tar.gz:
	mkdir -p $(@D)
	wget -O$@ $(YARA_URL)

.PHONY: yara-src
yara-src: $(YARA_SRC)/configure
$(YARA_SRC)/configure: $(YARA_SRC).tar.gz
	tar -C $(dir $(@D)) -xzf $^
	( cd $(@D) && patch -p1 < $(CURDIR)/data/yara/yara-pr-1416-backport.patch )
	cd $(@D) && ./bootstrap.sh

.PHONY: yara-windows-386
yara-windows-386: $(YARA_SRC)-windows-386/done
$(YARA_SRC)-windows-386/done: $(YARA_SRC)/configure
	mkdir -p $(@D)
	cd $(@D) && \
		$^ --prefix=$(CURDIR)/_non-golang/prefix-windows-386 --host=i686-w64-mingw32 \
			 --disable-magic --disable-cuckoo --without-crypto
	$(MAKE) -C $(@D) install
	touch $@

.PHONY: yara-windows-amd64
yara-windows-amd64: $(YARA_SRC)-windows-amd64/done
$(YARA_SRC)-windows-amd64/done: $(YARA_SRC)/configure
	mkdir -p $(@D)
	cd $(@D) && \
		$^ --prefix=$(CURDIR)/_non-golang/prefix-windows-amd64 --host=x86_64-w64-mingw32 \
			 --disable-magic --disable-cuckoo --without-crypto
	$(MAKE) -C $(@D) install
	touch $@


#==============================================================================
# Environment
#==============================================================================
.PHONY: check-env
check-env:
	@mkdir -p $(BUILD_FOLDER)

ifndef RULES
	@echo "[check-env] You have not specified any RULES env, kraken will not have any default Yara rules."
endif

ifndef BACKEND
	@echo "[check-env] You have not specified any BACKEND env, kraken will not have any default backend server configured."
endif


#==============================================================================
# Rules Compiler
#==============================================================================
.PHONY: rules-compiler
rules-compiler: $(BUILD_FOLDER)/rules-compiler $(BUILD_FOLDER)/rules-compiler $(CURDIR)/cmd/kraken/bindata.go

$(BUILD_FOLDER)/rules-compiler: $(COMPILER_SRC)
	@mkdir -p $(@D)
	@echo "[rules-compiler] Building rules compiler..."
	@cd cmd/rules-compiler; go build -o $@

$(BUILD_FOLDER)/rules: $(BUILD_FOLDER)/rules-compiler
ifdef RULES
	@echo "[rules-compiler] Compiling Yara rules..."
	cd $(BUILD_FOLDER); ./rules-compiler $(RULES)
else
	$(error "RULES has not been specified")
endif

$(CURDIR)/cmd/kraken/bindata.go: $(BUILD_FOLDER)/rules
	@echo "[rules-compiler] Launching binary resource builder..."
	go-bindata -o $@ $^


#==============================================================================
# Linux
#==============================================================================
.PHONY: linux
linux: check-env rules-compiler $(BUILD_FOLDER)/linux/kraken $(BUILD_FOLDER)/linux/kraken-launcher

$(BUILD_FOLDER)/linux/kraken: $(KRAKEN_SRC)
	@mkdir -p $(@D)
	@echo "[builder] Building Linux executable..."
	@cd cmd/kraken; $(FLAGS_LINUX) go build --ldflags '-s -w -extldflags "-lm -static" -X main.DefaultBaseDomain=$(BACKEND) -X main.AgentVersion=$(AGENT_VERSION)' \
		-tags yara_static -o $@
	@echo "[builder] Done!"

$(BUILD_FOLDER)/linux/kraken-launcher: $(LAUNCHER_SRC)
	@mkdir -p $(@D)
	@echo "[builder] Building Linux launcher..."
	@cd cmd/kraken-launcher; $(FLAGS_LINUX) go build --ldflags '-s -w' -o $@
	@echo "[builder] Done!"


#==============================================================================
# Darwin
#==============================================================================
.PHONY: darwin
darwin: check-env rules-compiler  $(BUILD_FOLDER)/darwin/kraken $(BUILD_FOLDER)/darwin/kraken-launcher

$(BUILD_FOLDER)/darwin/kraken: $(KRAKEN_SRC)
	@mkdir -p $(@D)
	@echo "[builder] Building Darwin executable..."
	@cd cmd/kraken; $(FLAGS_DARWIN) go build --ldflags '-s -w -extldflags "-lm" -X main.DefaultBaseDomain=$(BACKEND) -X main.AgentVersion=$(AGENT_VERSION)' \
		-tags yara_static -o $@
	@echo "[builder] Done!"

$(BUILD_FOLDER)/darwin/kraken-launcher: $(LAUNCHER_SRC)
	@mkdir -p $(@D)
	@echo "[builder] Building Darwin launcher..."
	@cd cmd/kraken-launcher; $(FLAGS_DARWIN) go build --ldflags '-s -w' -o $@
	@echo "[builder] Done!"


#==============================================================================
# FreeBSD
#==============================================================================
.PHONY: freebsd
freebsd: check-env rules-compiler  $(BUILD_FOLDER)/freebsd/kraken $(BUILD_FOLDER)/freebsd/kraken-launcher

$(BUILD_FOLDER)/freebsd/kraken: $(KRAKEN_SRC)
	@mkdir -p $(@D)
	@echo "[builder] Building FreeBSD executable..."
	@cd cmd/kraken; $(FLAGS_FREEBSD) go build --ldflags '-s -w -extldflags "-lm -static" -X main.DefaultBaseDomain=$(BACKEND) -X main.AgentVersion=$(AGENT_VERSION)' \
		-tags yara_static -o $@
	@echo "[builder] Done!"

$(BUILD_FOLDER)/freebsd/kraken-launcher: $(LAUNCHER_SRC)
	@mkdir -p $(@D)
	@echo "[builder] Building FreeBSD launcher..."
	@cd cmd/kraken-launcher; $(FLAGS_FREEBSD) go build --ldflags '-s -w' -o $@
	@echo "[builder] Done!"


#==============================================================================
# Windows
#==============================================================================
.PHONY: windows
windows: windows-386 windows-amd64

.PHONY: windows-386
windows-386: check-env rules-compiler  $(BUILD_FOLDER)/windows-386/kraken.exe $(BUILD_FOLDER)/windows-386/kraken-launcher.exe

$(BUILD_FOLDER)/windows-386/kraken.exe: $(YARA_SRC)-windows-386/done $(KRAKEN_SRC)
	@mkdir -p $(@D)

	#@rsrc -manifest $(CURDIR)/data/windows/kraken.manifest -ico kraken.ico -o rsrc.syso
	@rsrc -arch 386 -manifest $(CURDIR)/data/windows/kraken.manifest -o $(BUILD_FOLDER)/rsrc_windows_386.syso

	@echo "[builder] Building Windows 32bit executable..."
	@cd cmd/kraken; $(FLAGS_WINDOWS_386) go build --ldflags '-s -w -extldflags "-static" -X main.DefaultBaseDomain=$(BACKEND) -X main.AgentVersion=$(AGENT_VERSION)' \
		-tags yara_static -o $@
	@echo "[builder] Done!"

$(BUILD_FOLDER)/windows-386/kraken-launcher.exe: $(LAUNCHER_SRC)
	@mkdir -p $(@D)
	@echo "[builder] Building Windows 32bit launcher..."
	@cd cmd/kraken-launcher; $(FLAGS_WINDOWS_386) go build --ldflags '-s -w -extldflags "-static" -H=windowsgui' -o $@
	@echo "[builder] Done!"

.PHONY: windows-amd64
windows-amd64: check-env rules-compiler  $(BUILD_FOLDER)/windows-amd64/kraken.exe $(BUILD_FOLDER)/windows-amd64/kraken-launcher.exe

$(BUILD_FOLDER)/windows-amd64/kraken.exe: $(YARA_SRC)-windows-amd64/done $(KRAKEN_SRC)
	@mkdir -p $(@D)

	#@rsrc -manifest $(CURDIR)/data/windows/kraken.manifest -ico kraken.ico -o rsrc.syso
	@rsrc -arch amd64 -manifest $(CURDIR)/data/windows/kraken.manifest -o $(BUILD_FOLDER)/rsrc_windows_amd64.syso

	@echo "[builder] Building Windows 64bit executable..."
	@cd cmd/kraken; $(FLAGS_WINDOWS_AMD64) go build --ldflags '-s -w -extldflags "-static" -X main.DefaultBaseDomain=$(BACKEND) -X main.AgentVersion=$(AGENT_VERSION)' \
		-tags yara_static -o $@
	@echo "[builder] Done!"

$(BUILD_FOLDER)/windows-amd64/kraken-launcher.exe: $(LAUNCHER_SRC)
	@mkdir -p $(@D)
	@echo "[builder] Building Windows 64bit launcher..."
	@cd cmd/kraken-launcher; $(FLAGS_WINDOWS_AMD64) go build --ldflags '-s -w -extldflags "-static" -H=windowsgui' -o $@
	@echo "[builder] Done!"


#==============================================================================
# Misc
#==============================================================================
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

.PHONY: clean
clean:
	rm -f $(CURDIR)/cmd/kraken/bindata.go
	rm -rf $(BUILD_FOLDER)
