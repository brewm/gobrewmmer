package brewmmer

import (
	"context"
	"errors"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/brewm/gobrewmmer/cmd/brewmserver/ds18b20"
	"github.com/brewm/gobrewmmer/cmd/brewmserver/global"
	"github.com/brewm/gobrewmmer/pkg/api/brewmmer"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
)

type sessionServiceServer struct{}

func NewSessionServiceServer() brewmmer.SessionServiceServer {
	return &sessionServiceServer{}
}

//
//
//
// func Sense(c *gin.Context) {
// 	m := Measurement{Timestamp: time.Now(), Temperature: ds18b20.ReadTemperature()}
// 	c.JSON(200, m)
// }

//
// CRUD part (Reads)
//

func (s *sessionServiceServer) Get(ctx context.Context, req *brewmmer.GetSessionRequest) (*brewmmer.GetSessionResponse, error) {
	log.WithFields(log.Fields{"id": req.Id}).Info("Getting session!")

	sqlStatement := `
	SELECT
	id,
	start_time,
	stop_time,
	note
	FROM sessions
	WHERE id=$1`
	row := global.BrewmDB.QueryRow(sqlStatement, req.Id)

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

	err = fetchMeasurements(session)
	if err != nil {
		return nil, err
	}

	return &brewmmer.GetSessionResponse{
		Session: session,
	}, nil
}

func (s *sessionServiceServer) GetActive(ctx context.Context, req *brewmmer.GetActiveSessionRequest) (*brewmmer.GetSessionResponse, error) {
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
	row := global.BrewmDB.QueryRow(sqlStatement)

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

	err = fetchMeasurements(session)
	if err != nil {
		return nil, err
	}

	return &brewmmer.GetSessionResponse{
		Session: session,
	}, nil
}

func (s *sessionServiceServer) List(ctx context.Context, req *brewmmer.ListSessionRequest) (*brewmmer.ListSessionResponse, error) {
	log.Info("Getting all sessions!")
	sessions := []*brewmmer.Session{}

	rows, err := global.BrewmDB.Query(`
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

func fetchMeasurements(session *brewmmer.Session) error {
	log.WithFields(log.Fields{"session_id": session.Id}).Debug("Getting measurements!")

	sqlStatement := `
    SELECT
      timestamp,
      temperature
    FROM measurements
    WHERE session_id=$1`

	rows, err := global.BrewmDB.Query(sqlStatement, session.Id)

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

var sessionChannel chan struct{}

// in seconds
const measureInterval = 600

func (s *sessionServiceServer) Start(ctx context.Context, req *brewmmer.StartSessionRequest) (*brewmmer.StartSessionResponse, error) {
	if sessionChannel != nil {
		return nil, errors.New("session is already in progress, one session can be active at a time")
	}

	timestamp := time.Now()

	sqlStatement := `
    INSERT INTO sessions (start_time, note)
    VALUES ($1, $2)`

	result, err := global.BrewmDB.Exec(sqlStatement, timestamp, req.Note)
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
	startSessionProcess(int(sessionID))

	return &brewmmer.StartSessionResponse{
		Id: sessionID,
	}, nil
}

func startSessionProcess(id int) {
	sessionChannel = make(chan struct{})

	// Start goroutine to periodically run the insert
	go func(id int) {
		for {
			// Start goroutine to do an async insert
			go insertTemperature(id)

			time.Sleep(measureInterval * time.Second)
			select {
			case <-sessionChannel:
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

func insertTemperature(id int) {
	sqlStatement := `
    INSERT INTO measurements (session_id, timestamp, temperature)
    VALUES ($1, $2, $3);`

	_, err := global.BrewmDB.Exec(sqlStatement, id, time.Now(), ds18b20.ReadTemperature())
	if err != nil {
		log.WithFields(log.Fields{
			"session_id": id,
		}).Error("Failed to save measurement!")
	}
}

func (s *sessionServiceServer) Stop(ctx context.Context, req *brewmmer.StopSessionRequest) (*brewmmer.StopSessionResponse, error) {
	sqlStatement := `
    SELECT (CASE WHEN stop_time IS NULL THEN 1 ELSE 0 END) as is_active
    FROM sessions
    WHERE id = $1`
	row := global.BrewmDB.QueryRow(sqlStatement, req.Id)

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

	_, err = global.BrewmDB.Exec(sqlStatement, timestamp, req.Id)
	if err != nil {
		return nil, err
	}

	if sessionChannel != nil {
		log.WithFields(log.Fields{
			"session_id": req.Id,
		}).Info("Stopping session background process!")
		close(sessionChannel)
		sessionChannel = nil
	}

	return &brewmmer.StopSessionResponse{}, nil
}
