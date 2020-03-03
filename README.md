# gizmo
Service checker and scoreboard for CCDC-like exercises.

# Scoreboard Example

![](scoreboard_example.png)

# Usage

```
usage: gizmo [-h|--help] -i|--input "<value>" [-o|--output "<value>"]
             [-s|--script_directory "<value>"]

             Service uptime scoreboard

Arguments:

  -h  --help              Print help information
  -i  --input             Input config filename
  -o  --output            Output database filename. Default: gizmo.db
  -s  --script_directory  Script directory. Default: 

```

# Releases

Static binaries can be located at [gizmo/releases/](https://github.com/whoismissing/gizmo/releases)

An example JSON configuration file is located at [gizmo/config/examples/gizmo_config.json](https://github.com/whoismissing/gizmo/blob/master/config/examples/gizmo_config.json)

# Installation

```
go get -u -v github.com/akamensky/argparse
go get -u -v github.com/jlaffaye/ftp
go get -u -v github.com/mattn/go-sqlite3
go get -u -v golang.org/x/crypto/ssh
go get -u -v github.com/whoismissing/gizmo
go build
```

# Contribution

See [CONTRIBUTING.md](CONTRIBUTING.md)
