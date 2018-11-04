# brewmmer

## What
A semi automated brewing machine implemented in Go.


## Server
API server endpoints
```
GET    /ping                 - Ping
GET    /v1/sense             - Get the current temperature from the sensor
GET    /v1/sessions/         - Get all session information
GET    /v1/sessions/:id      - Get one session with all of the recorded measurements
POST   /v1/sessions/         - Create new session and start a goroutine to record temperature measurements
PUT    /v1/sessions/         - Stop session / stop goroutine
```

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
   0.1

COMMANDS:
     get      get <temperature|sessions>
     start    start <session|tbd...>
     stop     stop <session|tbd...>
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

To connect to a remote `brewmserver` set the `BREWM_ENDPOINT` env variable:
```
BREWM_ENDPOINT=http://192.168.0.22:8080 brewmctl get sessions
```