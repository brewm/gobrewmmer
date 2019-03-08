package brewmmer

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/brewm/gobrewmmer/cmd/brewmserver/global"
	"github.com/brewm/gobrewmmer/pkg/api/brewmmer"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
)

type sessionServiceServer struct{}

func NewSessionServiceServer() brewmmer.SessionServiceServer {
	return &sessionServiceServer{}
}

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
	}).Info("Found session!")

	err = fetchMeasurements(session)
	if err != nil {
		return nil, err
	}

	log.Info("Returning!")
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
	}).Info("Found session!")

	err = fetchMeasurements(session)
	if err != nil {
		return nil, err
	}

	log.Info("Returning!")
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

	log.Info("Returning!")
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
	log.WithFields(log.Fields{"session_id": session.Id}).Info("Getting measurements!")

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

func (s *sessionServiceServer) Start(ctx context.Context, req *brewmmer.StartSessionRequest) (*brewmmer.StartSessionResponse, error) {
	return nil, nil
}

func (s *sessionServiceServer) Stop(ctx context.Context, req *brewmmer.StopSessionRequest) (*brewmmer.StopSessionResponse, error) {
	return nil, nil
}
