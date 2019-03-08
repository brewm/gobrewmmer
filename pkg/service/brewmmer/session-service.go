package brewmmer

import (
	"context"
	"time"

	"github.com/brewm/gobrewmmer/cmd/brewmserver/global"
	"github.com/brewm/gobrewmmer/pkg/api/brewmmer"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
)

type sessionServiceServer struct{}

func NewSessionServiceServer() brewmmer.SessionServiceServer {
	return &sessionServiceServer{}
}

func (s *sessionServiceServer) Get(context.Context, *brewmmer.GetSessionRequest) (*brewmmer.GetSessionResponse, error) {
	return nil, nil
}

func (s *sessionServiceServer) GetActive(context.Context, *brewmmer.GetSessionRequest) (*brewmmer.GetSessionResponse, error) {
	return nil, nil
}

func (s *sessionServiceServer) List(context.Context, *brewmmer.ListSessionRequest) (*brewmmer.ListSessionResponse, error) {
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

		s := new(brewmmer.Session)

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

	return &brewmmer.ListSessionResponse{
		Sessions: sessions,
	}, nil
}

func (s *sessionServiceServer) Start(context.Context, *brewmmer.StartSessionRequest) (*brewmmer.StartSessionResponse, error) {
	return nil, nil
}

func (s *sessionServiceServer) Stop(context.Context, *brewmmer.StopSessionRequest) (*brewmmer.StopSessionResponse, error) {
	return nil, nil
}
