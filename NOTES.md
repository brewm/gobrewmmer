# RPI setup
## Install Go
```
wget https://storage.googleapis.com/golang/go1.11.1.linux-armv6l.tar.gz
sudo tar -C /usr/local -xvf go1.10.1.linux-armv6l.tar.gz
cat >> ~/.bashrc << 'EOF'
export GOPATH=$HOME/go
export PATH=/usr/local/go/bin:$PATH:$GOPATH/bin
EOF
source ~/.bashrc

```
from: https://gist.github.com/random-robbie/1f7f94beda1221b8125b62abe35f45b6

## Install go dependencies
`go get ...`
TODO: use godeps

## Setup database
```
sqlite3 /path/to/brewmmer.db < ./scripts/setup_db.sql
```

Note:
You can seed the database with some test data
```
sqlite3 /path/to/brewmmer.db < ./scripts/seed_test_data.sql
```

# GoBrewmmer

## Structure
https://golang.org/doc/code.html
https://github.com/golang-standards/project-layout

## Libs
Trying to use a minimal amount of external packages to control complexity. E.g. use the plain 'database/sql' instead of ORM.

### API
https://github.com/gin-gonic/gin

### Rpi connections
https://github.com/google/periph
-- I can't make this work with the onewire interface

# Dev process

## Local development
Everything is working except the ds18b20 package. The application is not crashing though, and ERR is printed and 0.0 temperature is returned.

## Remote development
Run the `./scripts/sync_pi.sh` script to rsync to local changes to RPI. Make sure you configure the $RPI_PORT and $RPI_HOST env variables.
