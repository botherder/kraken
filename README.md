<img src="graphics/kraken.png" width="450" />

Kraken is a simple cross-platform Yara scanner that can be built for Windows, Mac, FreeBSD and Linux. It is primarily intended for incident response, research and ad-hoc detections (*not* for endpoint protection). Following are the core features:

- Scan running executables and memory of running processes with provided Yara rules (leveraging [go-yara](https://github.com/hillu/go-yara)).
- Scan executables installed for autorun (leveraging [go-autoruns](https://github.com/botherder/go-autoruns)).
- Scan the filesystem with the provided Yara rules.
- Report any detection to a remote server provided with a Django-based web interface.
- Run continuously and periodically check for new autoruns and scan any newly-executed processes. Kraken will store events in a local SQLite3 database and will keep copies of autorun and detected executables.

Some features are still under work or almost completed:

* Installer and launcher to automatically start Kraken at startup.
* Download updated Yara rules from the server.

## Screenshots

![](graphics/cmd.png)

![](graphics/linux.png)

## How to use

Launch Kraken with any of the available options:

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

## User Guide

For details on install, use and build Kraken you should refer to the [User Guide](https://kraken.gitbook.io/user-guide/). The original source files for the documentation are available [here](https://github.com/botherder/kraken-docs), please open any issue or pull request pertinent to documentation there.

## License

Kraken is released under the [GNU General Public License v3.0](LICENSE) and is copyrighted to [Claudio Guarnieri](https://nex.sx).
