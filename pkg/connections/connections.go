package connections

import "database/sql"


// BrewmmerDB is the connection handle for the database
// https://stackoverflow.com/questions/31218008/sharing-a-globally-defined-db-conn-with-multiple-packages-in-golang
var BrewmmerDB *sql.DB
