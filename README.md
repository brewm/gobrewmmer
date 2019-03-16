# brewmmer

## What
A semi automated brewing machine implemented in Go.


## Server
API endponts and messages are documented in the `./doc` dir

Start the server
```
$ BREWM_DB_PATH=/var/lib/brewm/brewm.db brewmserver &
```


## CLI
Check the help content to see the available commands:

```
$ brewmctl --help

NAME:
   brewmctl - command line interface to control brewmmer

USAGE:
   brewmctl [global options] command [command options] [arguments...]

VERSION:
   0.2

COMMANDS:
     get      get <temperature|sessions|recipes>
     start    start <session|tbd...>
     stop     stop <session|tbd...>
     create   create <recipe|tbd...>
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

Subcommand help is also available:
```
$ brewmctl get sessions --help
NAME:
   brewmctl get sessions -

USAGE:
   brewmctl get sessions [command options] <id>

OPTIONS:
   --active
```

Subcommand help is also available:
```
$ brewmctl get sessions --help
NAME:
   brewmctl get sessions -

USAGE:
   brewmctl get sessions [command options] <id>

OPTIONS:
   --active
```

To connect to a remote `brewmserver` set the `BREWM_ENDPOINT` env variable:
```
BREWM_ENDPOINT=http://192.168.0.22:8080 brewmctl get sessions
```
