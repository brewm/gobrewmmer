-- Seed test data
INSERT INTO sessions (start_time, stop_time, note)
VALUES ('2006-01-02T15:04:05.999999999', '2006-01-02T15:10:09.999999999', 'First test session');

INSERT INTO measurements (session_id, timestamp, temperature)
VALUES ('1', '2006-01-02T15:04:08.999999999', '21.3');
INSERT INTO measurements (session_id, timestamp, temperature)
VALUES ('1', '2006-01-02T15:07:08.999999999', '21.6');
INSERT INTO measurements (session_id, timestamp, temperature)
VALUES ('1', '2006-01-02T15:10:08.999999999', '21.8');


INSERT INTO sessions (start_time, note)
VALUES ('2006-02-02T15:04:05.999999999', 'Second test session');

INSERT INTO measurements (session_id, timestamp, temperature)
VALUES ('2', '2006-02-02T15:04:05.999999999', '21.3');
INSERT INTO measurements (session_id, timestamp, temperature)
VALUES ('2', '2006-02-02T15:07:05.999999999', '21.6');
INSERT INTO measurements (session_id, timestamp, temperature)
VALUES ('2', '2006-02-02T15:10:05.999999999', '21.8');

INSERT INTO recipes (id, recipe)
VALUES (1, '{"id":1,"name":"Little Cub IPA - Bronze Cub","description":"Brown IPA","ingredients":[{"type":"WATER","name":"Init water","quantity":{"volume":32,"unit":"L"}},{"type":"WATER","name":"Sparge water","quantity":{"volume":25,"unit":"L"}},{"type":"MALT","name":"Pale Ale","quantity":{"volume":8,"unit":"KG"}},{"type":"MALT","name":"Crystal","quantity":{"volume":1.5,"unit":"KG"}},{"type":"HOPS","name":"Tomahawk","quantity":{"volume":75,"unit":"G"}},{"type":"HOPS","name":"Cascade","quantity":{"volume":80,"unit":"G"}},{"type":"HOPS","name":"Citra","quantity":{"volume":170,"unit":"G"}},{"type":"YIEST","name":"Safale US-05","quantity":{"volume":20,"unit":"G"}}],"steps":[{"temperature":55,"ingredients":[{"type":"WATER","name":"Init water","quantity":{"volume":32,"unit":"L"}},{"type":"MALT","name":"Pale Ale","quantity":{"volume":8,"unit":"KG"}},{"type":"MALT","name":"Crystal","quantity":{"volume":1.5,"unit":"KG"}}]},{"phase":"MASHING","temperature":55,"duration":{"volume":30,"unit":"MIN"}},{"phase":"MASHING","temperature":65,"duration":{"volume":45,"unit":"MIN"}},{"phase":"MASHING","temperature":78,"duration":{"volume":45,"unit":"MIN"}},{"phase":"SPARGING","temperature":55,"ingredients":[{"type":"WATER","name":"Sparge water","quantity":{"volume":25,"unit":"L"}}]},{"phase":"BOILING","duration":{"volume":60,"unit":"MIN"},"ingredients":[{"type":"HOPS","name":"Tomahawk","quantity":{"volume":75,"unit":"G"}}]},{"phase":"BOILING","duration":{"volume":10,"unit":"MIN"},"ingredients":[{"type":"HOPS","name":"Cascade","quantity":{"volume":80,"unit":"G"}}]},{"phase":"BOILING","ingredients":[{"type":"HOPS","name":"Citra","quantity":{"volume":75,"unit":"G"}}]},{"phase":"FERMENTATION","duration":{"volume":7,"unit":"DAY"}},{"phase":"FERMENTATION","duration":{"volume":7,"unit":"DAY"},"ingredients":[{"type":"HOPS","name":"Citra","quantity":{"volume":95,"unit":"G"}}]}]}');
