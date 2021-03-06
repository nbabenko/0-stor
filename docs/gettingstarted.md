## Installation

The server is go-getable:
```
go get -u github.com/zero-os/0-stor/cmd/zerostorserver
```


## Running the server
```
NAME:
   zerostorserver - Generic object store used by zero-os

USAGE:
   server [global options] command [command options] [arguments...]

VERSION:
   0.0.1

DESCRIPTION:
   Generic object store used by zero-os

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --debug, -d               Enable debug logging
   --bind value, -b value    Bind address (default: ":8080")
   --data value              Data directory (default: ".db/data")
   --meta value              Metadata directory (default: ".db/meta")
   --profile-addr value      Enables profiling of this server as an http service
   --auth-disable            Disable JWT authentification [$STOR_TESTING]
   --max-msg-size value      configure the maximum size of the message GRPC server can receive, in MiB (default: 32)
   --async-write             enable asynchonous writes (default: false)
   --help, -h                show help
   --version, -v             print the version

```

Start the server with grpc listening on all interfaces and port 12345
```shell
./zerostorserver --bind :12345 --data /path/to/data --meta /path/to/meta --interface grpc
```