package brewmmer

import (
	"context"
	"strings"

	"github.com/brewm/gobrewmmer/cmd/brewmserver/global"
	"github.com/brewm/gobrewmmer/pkg/api/brewmmer"
	"github.com/golang/protobuf/jsonpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type recepieServiceServer struct{}

func NewRecepieServiceServer() brewmmer.RecepieServiceServer {
	return &recepieServiceServer{}
}

func (s *recepieServiceServer) Create(ctx context.Context, req *brewmmer.CreateRecepieRequest) (*brewmmer.CreateRecepieResponse, error) {
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

	return &brewmmer.CreateRecepieResponse{
		Id: id,
	}, nil
}

func (s *recepieServiceServer) Get(ctx context.Context, req *brewmmer.GetRecepieRequest) (*brewmmer.GetRecepieResponse, error) {
	sqlStatement := `
    SELECT
      recepie
    FROM recepies
    WHERE id=$1`
	row := global.BrewmDB.QueryRow(sqlStatement, req.Id)

	var recepieJson string
	row.Scan(&recepieJson)

	um := jsonpb.Unmarshaler{}
	unserialized := &brewmmer.Recepie{}
	err := um.Unmarshal(strings.NewReader(recepieJson), unserialized)
	if err != nil {
		return nil, status.Error(codes.Unknown, "json unmarshaling error-> "+err.Error())
	}

	return &brewmmer.GetRecepieResponse{
		Recepie: unserialized,
	}, nil
}

func (s *recepieServiceServer) Delete(ctx context.Context, req *brewmmer.DeleteRecepieRequest) (*brewmmer.DeleteRecepieResponse, error) {
	return nil, nil
}

func (s *recepieServiceServer) Update(ctx context.Context, req *brewmmer.UpdateRecepieRequest) (*brewmmer.UpdateRecepieResponse, error) {
	return nil, nil
}

func (s *recepieServiceServer) List(ctx context.Context, req *brewmmer.ListRecepieRequest) (*brewmmer.ListRecepieResponse, error) {
	recepies := []*brewmmer.Recepie{}
	um := jsonpb.Unmarshaler{}

	rows, err := global.BrewmDB.Query(`
    SELECT
      recepie
    FROM recepies`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {

		var recepieJson *string
		err = rows.Scan(&recepieJson)

		recepie := &brewmmer.Recepie{}
		err := um.Unmarshal(strings.NewReader(*recepieJson), recepie)
		if err != nil {
			return nil, status.Error(codes.Unknown, "json unmarshaling error-> "+err.Error())
		}

		recepies = append(recepies, recepie)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return &brewmmer.ListRecepieResponse{
		Recepies: recepies,
	}, nil
}
