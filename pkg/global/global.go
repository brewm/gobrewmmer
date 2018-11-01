package connections

import (
  "database/sql"
)

// BrewmDB is the connection handle for the database
// https://stackoverflow.com/questions/31218008/sharing-a-globally-defined-db-conn-with-multiple-packages-in-golang
var BrewmDB *sql.DB