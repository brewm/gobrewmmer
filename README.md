# brewmmer

## What
A semi automated brewing machine implemented in Go.


## Server
API server endpoints
```
GET    /ping
GET    /v1/sense
GET    /v1/sessions/
GET    /v1/sessions/:id
POST   /v1/sessions/
PUT    /v1/sessions/
```

Start the server
```
$ BREWM_DB_PATH=/var/lib/brewm/brewm.db brewmserver &
```


## CLI

```
$ brewmctl

NAME:
   brewmctl - command line interface to control brewmmer

USAGE:
   brewmctl [global options] command [command options] [arguments...]

VERSION:
   0.1

COMMANDS:
     get      get <resource>
     start    start <process>
     stop     stop <process>
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

`brewmctl` can work remotely, e.g.
```
BREWM_ENDPOINT=http://192.168.0.22:8080 brewmctl get sessions
```