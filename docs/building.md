# Building

In order to build Kraken you will need to have Go installed on your system. We recommend using Go >= 1.11 in order to leverage the native support for Go Modules (if it is not available in your package manager, you can use something like [gvm](https://github.com/moovweb/gvm)).

Firstly, download Kraken:

    $ git clone https://github.com/botherder/kraken.git
    $ cd kraken

Most Go libraries dependencies are available to install through:

    $ make deps
