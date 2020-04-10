# confgen
Command-line tool to help users generate a config file for gizmo

# Usage

```
Usage:  ./confgen config_filename
```

# Build

```
cd ~/go/src/github.com/whoismissing/gizmo/tools/confgen/
go build
```

# Example Run

```
$ ./confgen test.config

Config Filename: test.config
Creating a new Game config ...
Adding team '0' ...
        Adding service '0' ...
        Enter the new value, press ENTER for default
        ServiceName [team0-service0]:
                Service name is default team0-service0
        HostIP [127.0.0.1]:
                Host ip is default 127.0.0.1
        Enter the new value, no default option
        ServiceType [www/dns/ext/ftp/ssh/ping]: www
                        Enter a URL: google.com
Add another service? [y/n] n
Add another team? [y/n] y
Adding team '1' ...
        Adding service '0' ...
        Enter the new value, press ENTER for default
        ServiceName [team1-service0]:
                Service name is default team1-service0
        HostIP [127.0.0.1]:
                Host ip is default 127.0.0.1
        Enter the new value, no default option
        ServiceType [www/dns/ext/ftp/ssh/ping]: ping
Add another service? [y/n] n
Add another team? [y/n] n
```
