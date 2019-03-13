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

CREATE TABLE recipes (
   id           INTEGER      PRIMARY KEY  NOT NULL,
   recipe       TEXT
);

CREATE TABLE brews (
   id           INTEGER      PRIMARY KEY  NOT NULL,
   recipe_id    INTEGER                   NOT NULL,
   start_time   TIMESTAMP                 NOT NULL,
   completed_time    TIMESTAMP,
   note         CHAR(255)
);

CREATE TABLE brew_steps (
   start_time        TIMESTAMP  NOT NULL,
   completed_time    TIMESTAMP  NOT NULL,
   step              TEXT,
   brew_id           INTEGER    NOT NULL,
   FOREIGN KEY (brew_id) REFERENCES brews (id)
);
