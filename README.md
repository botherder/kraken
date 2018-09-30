<img src="graphics/kraken.png" width="450" />

Kraken is a simple cross-platform Yara scanner that can be built for Windows, Mac, FreeBSD and Linux. It is primarily intended for incident response, research and ad-hoc detections (*not* for endpoint protection). Following are the core features:

- Scan running executables and memory of running processes with provided Yara rules (leveraging [go-yara](https://github.com/hillu/go-yara)).
- Scan executables installed for autorun (leveraging [go-autoruns](https://github.com/botherder/go-autoruns)).
- Report any detection to a remote server provided with a Django-based web interface.
- Run continuously and periodically check for new autoruns and scan any newly-executed processes. Kraken will store events in a local SQLite3 database and will keep copies of autorun and detected executables.

Some features are still under work or almost completed:

- Installer and launcher to automatically start Kraken at startup.
- Download updated Yara rules from the server.

## Table of Contents

- [Screenshots](#screenshots)
- [How to Use](#how-to-use)
  - [Configuration](#configuration)
- [Installing the Web Interface](#installing-the-web-interface)
- [Building](#building)
  - [Building on Linux](#building-on-linux)
  - [Building on FreeBSD](#building-on-freebsd)
  - [Building on Mac](#building-on-mac)
  - [Cross-Compiling Windows Binaries](#cross-compiling-windows-binaries)
- [License](#license)

## Screenshots

![](graphics/cmd.png)

![](graphics/linux.png)

![](graphics/windows.png)

## How to use

Once the binaries are compiled you will have a `kraken-launcher` and a `kraken` in the appropriate platform build folder.

`kraken` can be launched without any arguments and it will perform a scan of detected autorun entries and running processes and terminate. It will not communicate any results to any remote server.

Alternatively, `kraken` can also be launched using the following arguments:

    Usage of kraken:
          --backend string   Specify a particular hostname to the backend to connect to (overrides the default)
          --daemon           Enable daemon mode (this will also enable the report flag)
          --debug            Enable debug logs
          --folder string    Specify a particular folder to be scanned (overrides the default full filesystem)
          --no-autoruns      Disable scanning of autoruns
          --no-filesystem    Disable scanning of filesystem
          --no-process       Disable scanning of running processes
          --report           Enable reporting of events to the backend
          --rules            Specify a particular path to a file or folder containing the Yara rules to use

Using `kraken --backend example.com` will override the default `BACKEND` that was provided during build time.

Using `kraken --report` will make Kraken report any autoruns or detections to the configured backend server.

Launching `kraken --daemon` will execute a first scan and then run continuously. In *daemon* mode Kraken will monitor any new process creation and scan its binary and memory, as well as check regularly for any new entries registered for autorun. Enabling `--daemon` will automatically enable `--report` as well, even when not explicitly specified.

Enabling `--debug` will only display all debug log messages, mostly including details on files and processes being scanned.

Using `--no-autoruns`, `--no-filesystem` or `--no-process` will disable the scanning of autoruns, files stored on disk and running processes, respectively. Note: these flags do not impact the behavior of kraken when running in daemon mode.

If filesystem scanning is enabled, Kraken will recursively scan the entire root folder (`/` on \*nix systems and any fixed drive mounted on Windows systems). Using `--folder` you can specify a particular folder you want to scan instead.

The `--rules` option allows you to specify a path to a file or folder containing the Yara rules you want to use for your scanning. If the compilation of any of these rules fails (for example, because they include modules that are not enabled in the default Yara library), the execution will be aborted. If no `--rules` option is specified, Kraken will attempt to load a compiled rules file using the following order:

1. It will look for a compiled `rules` file in the current working directory.
2. It will look for a ocmpiled `rules` file in the local Kraken storage folder, in case it is running in *daemon* mode.
3. It will attempt to extract the compiled `rules` file from the embedded assets generated at build time (as explained in the [Building](#building) section).

If no compiled `rules` file is found, Kraken's Yara scanner will be disabled and execution will continue without it.

### Configuration

When `kraken` is launched in *daemon* mode it will look for a configuration file in either the current working directory or in the persistent directory. This configuration file is mostly used to look up the hostname of the backend Kraken will have to connect to. If a configuration file does not exist, it will create one using the default parameters provided during build time (primarily `BACKEND`).

If `kraken` is launched in normal mode, it will still look for any configuration file, but it will not write one to disk in the case there isn't one. If no configuration file is found, it will use the default parameters provided provided during build time (again, `BACKEND`).

To provide it different parameters you can create a `config.yaml` file in the same directory as the `kraken` binary using the following format:

```yaml
base_domain: <value>
```

Alternatively, you can specify a custom backend from the command line using `kraken --backend example.com`.

## Installing the Web Interface

The web interface is built using Django. You can run it using Python 3, which will require the following dependencies:

    $ sudo apt install python3 python3-dev python3-pip python3-mysqldb
    $ sudo pip3 install Django python-decouple django-geoip2-extras

To configure your Krakan Django app, instead of modifying `server/settings.py` you can create a file named `.env` inside the `server/` folder with the following content:

```shell
SECRET_KEY=your_secret_key
DEBUG=True
DB_NAME=kraken
DB_USER=user
DB_PASSWORD=pass
STATIC_ROOT=/home/user/kraken/server/static/
GEOIP_PATH=/home/user/geoip/
```

Change those values appropriately. The `GEOIP_PATH` variable should point to a folder containing your [MaxMind GeoLite2 City](https://dev.maxmind.com/geoip/geoip2/geolite2/) database.

After having configured the settings in the `.env` file, you will need to initialize the database with:

    $ python3 manage.py makemigrations autoruns detections downloads hosts
    $ python3 manage.py migrate

If you want to run the server using Gunicorn, you can install it with:

    $ sudo pip3 install gunicorn

You can create a Gunicorn systemd service by creating a `kraken.service` file in `/etc/systemd/system` like the following:

    Description=Gunicorn Application Server handling Kraken Servers
    After=network.target

    [Service]
    User=user
    Group=www-data
    WorkingDirectory=/home/user/kraken/server/
    ExecStart=/usr/local/bin/gunicorn --workers 3 --bind unix:/home/user/kraken-server.sock server.wsgi:application
    Restart=always

    [Install]
    WantedBy=multi-user.target

You can then configure your webserver to proxy requests to the unix socket at `/home/user/kraken-server.sock`.

## Building

In order to build Kraken you will need to have Go installed on your system. We recommend using Go >= 1.11 in order to leverage the native support for Go Modules (if it is not available in your package manager, you can use something like [gvm](https://github.com/moovweb/gvm)).

Firstly, download Kraken:

    $ git clone https://github.com/botherder/kraken.git
    $ cd kraken

Most Go libraries dependencies are available to install through:

    $ make deps

### Building on Linux

You need to install Yara development libraries and headers. You should download and compile Yara from the [official sources](https://github.com/VirusTotal/yara). It will require `dh-autoreconf` to be installed and you will need to configure some compilation flags. This is most likely the procedure you will need to follow:

    $ sudo apt install dh-autoreconf
    $ wget https://github.com/VirusTotal/yara/archive/v3.8.1.tar.gz
    $ tar -zxvf yara-v3.8.1.tar.gz
    $ cd yara-3.8.1
    $ ./bootstrap.sh
    $ ./configure --without-crypto
    $ make && sudo make install
    $ sudo ldconfig

Compiling Kraken requires you to specify a path to a file or a folder that contains the Yara rules you wish to embed with the binary. You can try for example with:

    $ BACKEND=example.com RULES=test/ make linux

### Building on FreeBSD

While cross-compilation of FreeBSD binaries is not available yet, it is possible to build binaries in a native FreeBSD environment. In order to do so you will firstly need to install some packages:

    $ sudo pkg install git gmake pkgconf go-bindata

Then you will need to install Yara, which is normally available in [ports](https://www.freshports.org/security/yara/):

    $ cd /usr/ports/security/yara

Before installing it, it is recommended that you modify the file `Makefile` to add `--without-crypto` to `CONFIGURE_ARGS` (if you don't need the Yara modules enabled in the Makefile, feel free to remove them). Now you can proceed with installing:

    $ sudo make && sudo make install

Now you can move to the directory that contains the Kraken source code and build it with:

    $ BACKEND=example.com RULES=test/ gmake freebsd

### Building on Mac

While cross-compilation of Darwin binaries is not available yet, it is possible to build binaries in a native Mac environment. While Yara is generally available using [brew](https://brew.sh), we are going to compile the latest available sources. Before doing so, we need to install some dependencies:

    $ brew install automake libtool go-bindata

Now you can proceed to download and compile the Yara sources in the same way as explained in the [Building on Linux](#building-on-linux) section. Once Yara is compiled and install correctly, you can proceed building binaries for Mac using the follwing command:

    $ BACKEND=example.com RULES=test/ make darwin

### Cross-compiling Windows binaries

Cross-compiling Windows binaries from a Linux development machine is a slightly more complicated process. Firstly you will need to install MingW and some other dependencies:

    $ sudo apt install gcc mingw-w64 automate libtool make

Next you will need to download Yara sources. Use the latest available version, which at the time of this writing is 3.8.1:

    $ wget https://github.com/VirusTotal/yara/archive/v3.8.1.tar.gz

Unpack the archive and export `YARA_SRC` to the newly-created folder:

    $ export YARA_SRC=<folder>

Next you need to bootstrap Yara sources and compile them with MingW. These are the instructions to compile it for **32bit**:

    $ cd ${YARA_SRC}
    $ ./bootstrap.sh
    $ ./configure --host=i686-w64-mingw32 --without-crypto --prefix=${YARA_SRC}/i686-w64-mingw32
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

## License

Kraken is released under the [GNU General Public License v3.0](LICENSE) and is copyrighted to [Claudio Guarnieri](https://nex.sx).
