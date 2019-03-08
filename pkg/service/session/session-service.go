package session

import (
	"context"
	"time"

	"github.com/brewm/gobrewmmer/cmd/brewmserver/global"
	"github.com/brewm/gobrewmmer/pkg/api/session"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
)

type sessionServiceServer struct{}

func NewSessionServiceServer() session.SessionServiceServer {
	return &sessionServiceServer{}
}

func (s *sessionServiceServer) Get(context.Context, *session.GetRequest) (*session.GetResponse, error) {
	return nil, nil
}

func (s *sessionServiceServer) List(context.Context, *session.ListRequest) (*session.ListResponse, error) {
	sessions := []*session.Session{}

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

		s := new(session.Session)

		err = rows.Scan(
			&s.Id,
			&startTime,
			&nullableStopTime,
			&s.Note,
		)

		if err != nil {
			return nil, err
		}

		if nullableStopTime != nil {
			s.StopTime, err = ptypes.TimestampProto(*nullableStopTime)
			if err != nil {
				return nil, err
			}
		} else {
			s.StopTime = new(timestamp.Timestamp)
		}

		s.StartTime, err = ptypes.TimestampProto(*startTime)
		if err != nil {
			return nil, err
		}

		sessions = append(sessions, s)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return &session.ListResponse{
		Sessions: sessions,
	}, nil
}

func (s *sessionServiceServer) Start(context.Context, *session.StartRequest) (*session.StartResponse, error) {
	return nil, nil
}

func (s *sessionServiceServer) Stop(context.Context, *session.StopRequest) (*session.StopResponse, error) {
	return nil, nil
}
