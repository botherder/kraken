# Cross-compiling Windows binaries

Cross-compiling Windows binaries from a Linux development machine is a slightly more complicated process. Firstly you will need to install MingW and some other dependencies:

    $ sudo apt install gcc mingw-w64 automake libtool make go-bindata

Next you will need to download Yara sources. Use the latest available version, which at the time of this writing is 4.0.1:

    $ wget https://github.com/VirusTotal/yara/archive/v4.0.1.tar.gz

Unpack the archive and export `YARA_SRC` to the newly-created folder:

    $ export YARA_SRC=<folder>

Next you need to bootstrap Yara sources and compile them with MingW. These are the instructions to compile it for **32bit**:

    $ cd ${YARA_SRC}
    $ ./bootstrap.sh
    $ ./configure --host=i686-w64-mingw32 --disable-magic --disable-cuckoo --without-crypto --prefix=${YARA_SRC}/i686-w64-mingw32
    $ make -C ${YARA_SRC}
    $ make -C ${YARA_SRC} install

Now we can download and build `go-yara` for 32bit using the following command:

    $ go get -d -u github.com/hillu/go-yara
    $ GOOS=windows GOARCH=386 CGO_ENABLED=1 \
      CC=i686-w64-mingw32-gcc \
      PKG_CONFIG_PATH=${YARA_SRC}/i686-w64-mingw32/lib/pkgconfig \
      go install -ldflags '-extldflags "-static"' github.com/hillu/go-yara

Now you can compile Kraken using:

    $ BACKEND=example.com RULES=test/ make windows

If you get errors such as ` undefined reference to 'yr_compiler_add_file'` you might need to pass the `PKG_CONFIG_PATH` variable:

    $ PKG_CONFIG_PATH=${YARA_SRC}/i686-w64-mingw32/lib/pkgconfig BACKEND=example.com RULES=test/ make windows
