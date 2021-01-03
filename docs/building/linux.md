# Building on Linux

First, install some required dependencies before continuing:

    $ sudo apt install gcc automake libtool make go-bindata

You need to install Yara development libraries and headers. You should download and compile Yara from the [official sources](https://github.com/VirusTotal/yara). It will require `dh-autoreconf` to be installed and you will need to configure some compilation flags. This is most likely the procedure you will need to follow:

    $ sudo apt install dh-autoreconf
    $ wget https://github.com/VirusTotal/yara/archive/v4.0.1.tar.gz
    $ tar -zxvf yara-v4.0.1.tar.gz
    $ cd yara-3.8.1
    $ ./bootstrap.sh
    $ ./configure --without-crypto
    $ make && sudo make install
    $ sudo ldconfig

Compiling Kraken requires you to specify a path to a file or a folder that contains the Yara rules you wish to embed with the binary. You can try for example with:

    $ BACKEND=example.com RULES=test/ make linux

You might see some warning messages like the following:

	/usr/bin/ld: /tmp/go-link-1111111/000018.o: in function `mygetgrouplist':
	$GOPATH/src/os/user/getgrouplist_unix.go:16: warning: Using 'getgrouplist' in statically linked applications requires at runtime the shared libraries from the glibc version used for linking
	/usr/bin/ld: /tmp/go-link-1111111/000017.o: in function `mygetgrgid_r':
	$GOPATH/src/os/user/cgo_lookup_unix.go:38: warning: Using 'getgrgid_r' in statically linked applications requires at runtime the shared libraries from the glibc version used for linking
	/usr/bin/ld: /tmp/go-link-1111111/000017.o: in function `mygetgrnam_r':
	$GOPATH/src/os/user/cgo_lookup_unix.go:43: warning: Using 'getgrnam_r' in statically linked applications requires at runtime the shared libraries from the glibc version used for linking
	/usr/bin/ld: /tmp/go-link-1111111/000017.o: in function `mygetpwnam_r':
	$GOPATH/src/os/user/cgo_lookup_unix.go:33: warning: Using 'getpwnam_r' in statically linked applications requires at runtime the shared libraries from the glibc version used for linking
	/usr/bin/ld: /tmp/go-link-1111111/000017.o: in function `mygetpwuid_r':
	$GOPATH/src/os/user/cgo_lookup_unix.go:28: warning: Using 'getpwuid_r' in statically linked applications requires at runtime the shared libraries from the glibc version used for linking
	/usr/bin/ld: /tmp/go-link-1111111/000015.o: in function `_cgo_18049202ccd9_C2func_getaddrinfo':
	/tmp/go-build/cgo-gcc-prolog:49: warning: Using 'getaddrinfo' in statically linked applications requires at runtime the shared libraries from the glibc version used for linking

If so, don't alarm, as they shouldn't prevent the executables from being successfully built.

Once the `make linux` command is completed, you should see Kraken binaries inside `build/linux/`.
