# Introduction

![](.gitbook/assets/kraken.png)

Kraken is a simple cross-platform Yara scanner that can be built for Windows, Mac, FreeBSD and Linux. It is primarily intended for incident response, research and ad-hoc detections \(_not_ for endpoint protection\). Following are the core features:

- Scan running executables and memory of running processes with provided Yara rules (leveraging [go-yara](https://github.com/hillu/go-yara)).
- Scan executables installed for autorun (leveraging [go-autoruns](https://github.com/botherder/go-autoruns)).
- Scan the filesystem with the provided Yara rules.
- Report any detection to a remote server provided with a Django-based web interface.
- Run continuously and periodically check for new autoruns and scan any newly-executed processes. Kraken will store events in a local SQLite3 database and will keep copies of autorun and detected executables.

Some features are still under work or almost completed:

* Installer and launcher to automatically start Kraken at startup.
* Download updated Yara rules from the server.

## Screenshots

![](.gitbook/assets/cmd.png)

![](.gitbook/assets/linux.png)

![](.gitbook/assets/windows.png)

## License

Kraken is released under the [GNU General Public License v3.0](https://github.com/botherder/kraken/tree/636a4f467228ea84df7162ed748c47469263b60b/LICENSE/README.md) and is copyrighted to [Claudio Guarnieri](https://nex.sx).
