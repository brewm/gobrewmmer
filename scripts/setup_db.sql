CREATE TABLE sessions (
   id           INTEGER      PRIMARY KEY  NOT NULL,
   start_time   TIMESTAMP                 NOT NULL,
   stop_time    TIMESTAMP,
   note         CHAR(255)
);

CREATE TABLE measurements (
   timestamp    TIMESTAMP  NOT NULL,
   temperature  REAL       NOT NULL,
   session_id   INTEGER    NOT NULL,
   FOREIGN KEY (session_id) REFERENCES sessions (id)
);

