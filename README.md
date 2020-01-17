# gizmo
Service checker and scoreboard for CCDC-like exercises.

# Usage

```
usage: gizmo [-h|--help] -i|--input "<value>" [-o|--output "<value>"]

             Service uptime scoreboard

Arguments:

  -h  --help    Print help information
  -i  --input   Input config filename
  -o  --output  Output database filename. Default: gizmo.db

```

# Releases

Static binaries can be located at [gizmo/releases/](https://github.com/whoismissing/gizmo/releases/gizmo)

An example JSON configuration file is located at [gizmo/config/examples/gizmo_config.json](https://github.com/whoismissing/gizmo/config/examples/gizmo_config.json)

# Installation

```
go get -u github.com/whoismissing/gizmo
go build
```

# Contribution

See [CONTRIBUTING.md](CONTRIBUTING.md)
