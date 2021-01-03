# Building on FreeBSD

While cross-compilation of FreeBSD binaries is not available yet, it is possible to build binaries in a native FreeBSD environment. In order to do so you will firstly need to install some packages:

    $ sudo pkg install git gmake pkgconf go-bindata

Then you will need to install Yara, which is normally available in [ports](https://www.freshports.org/security/yara/):

    $ cd /usr/ports/security/yara

Before installing it, it is recommended that you modify the file `Makefile` to add `--without-crypto` to `CONFIGURE_ARGS` (if you don't need the Yara modules enabled in the Makefile, feel free to remove them). Now you can proceed with installing:

    $ sudo make && sudo make install

Now you can move to the directory that contains the Kraken source code and build it with:

    $ BACKEND=example.com RULES=test/ gmake freebsd
