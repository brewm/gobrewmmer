CREATE TABLE session(
   id           INTEGER      PRIMARY KEY  NOT NULL,
   start_time   TIMESTAMP                 NOT NULL,
   stop_time    TIMESTAMP,
   note         CHAR(255)
);

CREATE TABLE measurement(
   timestamp    TIMESTAMP  NOT NULL,
   temperature  REAL       NOT NULL,
   session_id   INTEGER    NOT NULL,
   FOREIGN KEY (session_id) REFERENCES session (id)
);

