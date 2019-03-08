package brewmmer

import (
	"context"

	"github.com/brewm/gobrewmmer/cmd/brewmserver/ds18b20"
	"github.com/brewm/gobrewmmer/pkg/api/brewmmer"
	"github.com/golang/protobuf/ptypes"
)

type temperatureServiceServer struct{}

func NewTemperatureServiceServer() brewmmer.TemperatureServiceServer {
	return &temperatureServiceServer{}
}

func (s *temperatureServiceServer) Get(ctx context.Context, req *brewmmer.GetTemperatureRequest) (*brewmmer.GetTemperatureResponse, error) {
	return &brewmmer.GetTemperatureResponse{
		Timestamp:   ptypes.TimestampNow(),
		Temperature: ds18b20.ReadTemperature(),
	}, nil

}
