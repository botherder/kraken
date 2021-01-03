# Building on Mac

While cross-compilation of Darwin binaries is not available yet, it is possible to build binaries in a native Mac environment. While Yara is generally available using [brew](https://brew.sh), we are going to compile the latest available sources. Before doing so, we need to install some dependencies:

    $ brew install automake libtool go-bindata

Now you can proceed to download and compile the Yara sources in the same way as explained in the [Building on Linux](#building-on-linux) section. Once Yara is compiled and install correctly, you can proceed building binaries for Mac using the follwing command:

    $ BACKEND=example.com RULES=test/ make darwin
