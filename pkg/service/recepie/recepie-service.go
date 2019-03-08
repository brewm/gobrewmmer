package recepie

import (
	"context"
	"strings"

	"github.com/brewm/gobrewmmer/cmd/brewmserver/global"
	"github.com/brewm/gobrewmmer/pkg/api/recepie"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type recepieServiceServer struct{}

func NewRecepieServiceServer() recepie.RecepieServiceServer {
	return &recepieServiceServer{}
}

func (s *recepieServiceServer) Create(ctx context.Context, req *recepie.CreateRequest) (*recepie.CreateResponse, error) {
	m := jsonpb.Marshaler{}
	recepieJson, err := m.MarshalToString(req.Recepie)
	if err != nil {
		return nil, status.Error(codes.Unknown, "json marshaling error-> "+err.Error())
	}

	sqlStatement := `INSERT INTO recepies (recepie) VALUES ($1)`
	res, err := global.BrewmDB.Exec(sqlStatement, recepieJson)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to insert into Recepie-> "+err.Error())
	}

	// get ID of creates ToDo
	id, err := res.LastInsertId()
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve id for created Recepie-> "+err.Error())
	}

	return &recepie.CreateResponse{
		Id: id,
	}, nil
}

func (s *recepieServiceServer) Get(ctx context.Context, req *recepie.GetRequest) (*recepie.GetResponse, error) {
	sqlStatement := `
    SELECT
      recepie
    FROM recepies
    WHERE id=$1`
	row := global.BrewmDB.QueryRow(sqlStatement, req.Id)

	var recepieJson string
	row.Scan(&recepieJson)

	um := jsonpb.Unmarshaler{}
	unserialized := &recepie.Recepie{}
	err := um.Unmarshal(strings.NewReader(recepieJson), unserialized)
	if err != nil {
		return nil, status.Error(codes.Unknown, "json unmarshaling error-> "+err.Error())
	}

	return &recepie.GetResponse{
		Recepie: unserialized,
	}, nil
}

func (s *recepieServiceServer) Delete(ctx context.Context, req *recepie.DeleteRequest) (*empty.Empty, error) {
	return nil, nil
}

func (s *recepieServiceServer) Update(ctx context.Context, req *recepie.UpdateRequest) (*empty.Empty, error) {
	return nil, nil
}

func (s *recepieServiceServer) List(ctx context.Context, empty *empty.Empty) (*recepie.ListResponse, error) {
	return &recepie.ListResponse{}, nil
}
