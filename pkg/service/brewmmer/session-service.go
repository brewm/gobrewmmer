package brewmmer

import (
	"context"
	"database/sql"
	"errors"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/brewm/gobrewmmer/pkg/api/brewmmer"
	"github.com/brewm/gobrewmmer/pkg/service/ds18b20"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
)

type sessionServiceServer struct {
	db             *sql.DB
	sessionChannel chan struct{}
}

func NewSessionServiceServer(db *sql.DB) brewmmer.SessionServiceServer {
	server := &sessionServiceServer{db: db}

	// Restart the session process if there is an acive
	sqlStatement := `
    SELECT MAX(id), (CASE WHEN stop_time IS NULL THEN 1 ELSE 0 END) as is_active
    FROM sessions`
	row := db.QueryRow(sqlStatement)

	var id int
	var isActive bool

	err := row.Scan(
		&id,
		&isActive,
	)

	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Checking active sessions failed!")
	}

	if isActive == true {
		log.WithFields(log.Fields{
			"session_id": id,
		}).Info("Re-starting session process after an apiserver restart!")
		startSessionProcess(db, &server.sessionChannel, id)
	}

	return server
}

//
// CRUD part (Reads)
//

func (sss *sessionServiceServer) Get(ctx context.Context, req *brewmmer.GetSessionRequest) (*brewmmer.GetSessionResponse, error) {
	log.WithFields(log.Fields{"id": req.Id}).Info("Getting session!")

	sqlStatement := `
	SELECT
	id,
	start_time,
	stop_time,
	note
	FROM sessions
	WHERE id=$1`
	row := sss.db.QueryRow(sqlStatement, req.Id)

	var nullableStopTime *time.Time
	var startTime *time.Time

	session := new(brewmmer.Session)
	err := row.Scan(
		&session.Id,
		&startTime,
		&nullableStopTime,
		&session.Note,
	)
	if err != nil {
		return nil, err
	}

	err = fillTime(session, startTime, nullableStopTime)
	if err != nil {
		return nil, err
	}

	log.WithFields(log.Fields{
		"id":         session.Id,
		"start_time": session.StartTime,
		"stop_time":  session.StopTime,
	}).Debug("Found session!")

	err = fetchMeasurements(sss.db, session)
	if err != nil {
		return nil, err
	}

	return &brewmmer.GetSessionResponse{
		Session: session,
	}, nil
}

func (sss *sessionServiceServer) GetActive(ctx context.Context, req *brewmmer.GetActiveSessionRequest) (*brewmmer.GetSessionResponse, error) {
	log.Info("Getting active session!")

	// Making sure to return only one session
	sqlStatement := `
	SELECT
	MAX(id),
	start_time,
	stop_time,
	note
	FROM sessions
	WHERE stop_time IS NULL`
	row := sss.db.QueryRow(sqlStatement)

	var nullableStopTime *time.Time
	var startTime *time.Time

	session := new(brewmmer.Session)
	err := row.Scan(
		&session.Id,
		&startTime,
		&nullableStopTime,
		&session.Note,
	)
	if err != nil {
		return nil, err
	}

	err = fillTime(session, startTime, nullableStopTime)
	if err != nil {
		return nil, err
	}

	log.WithFields(log.Fields{
		"id":         session.Id,
		"start_time": session.StartTime,
		"stop_time":  session.StopTime,
	}).Debug("Found session!")

	err = fetchMeasurements(sss.db, session)
	if err != nil {
		return nil, err
	}

	return &brewmmer.GetSessionResponse{
		Session: session,
	}, nil
}

func (sss *sessionServiceServer) List(ctx context.Context, req *brewmmer.ListSessionRequest) (*brewmmer.ListSessionResponse, error) {
	log.Info("Getting all sessions!")
	sessions := []*brewmmer.Session{}

	rows, err := sss.db.Query(`
    SELECT
      id,
      start_time,
      stop_time,
      note
    FROM sessions`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		// The sqlite driver can't handle nullable Time type so here is the workaround
		var nullableStopTime *time.Time
		var startTime *time.Time

		session := new(brewmmer.Session)

		err = rows.Scan(
			&session.Id,
			&startTime,
			&nullableStopTime,
			&session.Note,
		)

		if err != nil {
			return nil, err
		}

		err = fillTime(session, startTime, nullableStopTime)

		if err != nil {
			return nil, err
		}

		sessions = append(sessions, session)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return &brewmmer.ListSessionResponse{
		Sessions: sessions,
	}, nil
}

func fillTime(session *brewmmer.Session, startTime *time.Time, nullableStopTime *time.Time) error {
	var err error
	if nullableStopTime != nil {
		session.StopTime, err = ptypes.TimestampProto(*nullableStopTime)
		if err != nil {
			return err
		}
	} else {
		session.StopTime = new(timestamp.Timestamp)
	}

	session.StartTime, err = ptypes.TimestampProto(*startTime)
	if err != nil {
		return err
	}
	return nil
}

func fetchMeasurements(db *sql.DB, session *brewmmer.Session) error {
	log.WithFields(log.Fields{"session_id": session.Id}).Debug("Getting measurements!")

	sqlStatement := `
    SELECT
      timestamp,
      temperature
    FROM measurements
    WHERE session_id=$1`

	rows, err := db.Query(sqlStatement, session.Id)

	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var timestamp *time.Time

		m := new(brewmmer.Measurement)
		err = rows.Scan(
			&timestamp,
			&m.Temperature,
		)
		m.Timestamp, err = ptypes.TimestampProto(*timestamp)

		if err != nil {
			return err
		}

		session.Measurements = append(session.Measurements, m)
	}
	err = rows.Err()
	if err != nil {
		return err
	}

	return nil
}

//
// NOT CRUD PART
//

// in seconds
const measureInterval = 10

func (sss *sessionServiceServer) Start(ctx context.Context, req *brewmmer.StartSessionRequest) (*brewmmer.StartSessionResponse, error) {
	if sss.sessionChannel != nil {
		return nil, errors.New("session is already in progress, one session can be active at a time")
	}

	timestamp := time.Now()

	sqlStatement := `
    INSERT INTO sessions (start_time, note)
    VALUES ($1, $2)`

	result, err := sss.db.Exec(sqlStatement, timestamp, req.Note)
	if err != nil {
		return nil, err
	}

	sessionID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	log.WithFields(log.Fields{
		"session_id": int(sessionID),
	}).Info("Starting new session process!")
	startSessionProcess(sss.db, &sss.sessionChannel, int(sessionID))

	return &brewmmer.StartSessionResponse{
		Id: sessionID,
	}, nil
}

func startSessionProcess(db *sql.DB, sc *chan struct{}, id int) {
	*sc = make(chan struct{})

	// Start goroutine to periodically run the insert
	go func(id int) {
		for {
			// Start goroutine to do an async insert
			go insertTemperature(db, id)

			time.Sleep(measureInterval * time.Second)
			select {
			case <-*sc:
				log.WithFields(log.Fields{
					"session_id": id,
				}).Info("Stopping active session!")
				return
			default: // adding default will make it not block
				log.WithFields(log.Fields{
					"session_id": id,
				}).Debug("Rolling to next measurement!")
			}
		}
	}(id)
}

func insertTemperature(db *sql.DB, id int) {
	sqlStatement := `
    INSERT INTO measurements (session_id, timestamp, temperature)
    VALUES ($1, $2, $3);`

	_, err := db.Exec(sqlStatement, id, time.Now(), ds18b20.ReadTemperature())
	if err != nil {
		log.WithFields(log.Fields{
			"session_id": id,
		}).Error("Failed to save measurement!")
	}
}

func (sss *sessionServiceServer) Stop(ctx context.Context, req *brewmmer.StopSessionRequest) (*brewmmer.StopSessionResponse, error) {
	sqlStatement := `
    SELECT (CASE WHEN stop_time IS NULL THEN 1 ELSE 0 END) as is_active
    FROM sessions
    WHERE id = $1`
	row := sss.db.QueryRow(sqlStatement, req.Id)

	var isActive bool
	err := row.Scan(&isActive)

	if err != nil {
		log.WithFields(log.Fields{
			"id": req.Id,
		}).Error("Checking session with the given id failed!")
		return nil, err
	}

	if isActive == false {
		return nil, errors.New("given session is not active, can't stop")
	}

	sqlStatement = `
    UPDATE sessions
    SET stop_time = $1
    WHERE id = $2`

	timestamp := time.Now()

	_, err = sss.db.Exec(sqlStatement, timestamp, req.Id)
	if err != nil {
		return nil, err
	}

	if sss.sessionChannel != nil {
		log.WithFields(log.Fields{
			"session_id": req.Id,
		}).Info("Stopping session background process!")
		close(sss.sessionChannel)
		sss.sessionChannel = nil
	}

	return &brewmmer.StopSessionResponse{}, nil
}
