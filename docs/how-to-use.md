# How to use

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

## Configuration

When `kraken` is launched in *daemon* mode it will look for a configuration file in either the current working directory or in the persistent directory. This configuration file is mostly used to look up the hostname of the backend Kraken will have to connect to. If a configuration file does not exist, it will create one using the default parameters provided during build time (primarily `BACKEND`).

If `kraken` is launched in normal mode, it will still look for any configuration file, but it will not write one to disk in the case there isn't one. If no configuration file is found, it will use the default parameters provided provided during build time (again, `BACKEND`).

To provide it different parameters you can create a `config.yaml` file in the same directory as the `kraken` binary using the following format:

```yaml
base_domain: <value>
```

Alternatively, you can specify a custom backend from the command line using `kraken --backend example.com`.
