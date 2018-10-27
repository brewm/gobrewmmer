CREATE TABLE session(
   id           INTEGER      PRIMARY KEY  NOT NULL,
   start_time   TIMESTAMP                 NOT NULL,
   stop_time    INTEGER,
   note         CHAR(255)
);

CREATE TABLE measurement(
   timestamp    TIMESTAMP  NOT NULL,
   temperature  REAL       NOT NULL,
   session_id   INTEGER    NOT NULL,
   FOREIGN KEY (session_id) REFERENCES session (id)
);


-- Seed test data
INSERT INTO session (start_time, note)
VALUES ('2006-01-02T15:04:05.999999999', 'Initial test data');

INSERT INTO measurement (session_id, timestamp, temperature)
VALUES ('1', '2006-01-02T15:04:05.999999999', '21.3');
INSERT INTO measurement (session_id, timestamp, temperature)
VALUES ('1', '2006-01-02T15:07:05.999999999', '21.6');
INSERT INTO measurement (session_id, timestamp, temperature)
VALUES ('1', '2006-01-02T15:10:05.999999999', '21.8');
