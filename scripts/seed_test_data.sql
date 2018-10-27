-- Seed test data
INSERT INTO session (start_time, stop_time, note)
VALUES ('2006-01-02T15:04:05.999999999', '2006-01-02T15:10:09.999999999', 'First test session');

INSERT INTO measurement (session_id, timestamp, temperature)
VALUES ('1', '2006-01-02T15:04:08.999999999', '21.3');
INSERT INTO measurement (session_id, timestamp, temperature)
VALUES ('1', '2006-01-02T15:07:08.999999999', '21.6');
INSERT INTO measurement (session_id, timestamp, temperature)
VALUES ('1', '2006-01-02T15:10:08.999999999', '21.8');


INSERT INTO session (start_time, note)
VALUES ('2006-02-02T15:04:05.999999999', 'Second test session');

INSERT INTO measurement (session_id, timestamp, temperature)
VALUES ('1', '2006-02-02T15:04:05.999999999', '21.3');
INSERT INTO measurement (session_id, timestamp, temperature)
VALUES ('1', '2006-02-02T15:07:05.999999999', '21.6');
INSERT INTO measurement (session_id, timestamp, temperature)
VALUES ('1', '2006-02-02T15:10:05.999999999', '21.8');