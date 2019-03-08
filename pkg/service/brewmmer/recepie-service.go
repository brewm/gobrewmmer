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

type recipeServiceServer struct{}

func NewRecipeServiceServer() brewmmer.RecipeServiceServer {
	return &recipeServiceServer{}
}

func (s *recipeServiceServer) Create(ctx context.Context, req *brewmmer.CreateRecipeRequest) (*brewmmer.CreateRecipeResponse, error) {
	m := jsonpb.Marshaler{}
	recipeJson, err := m.MarshalToString(req.Recipe)
	if err != nil {
		return nil, status.Error(codes.Unknown, "json marshaling error-> "+err.Error())
	}

	sqlStatement := `INSERT INTO recipes (recipe) VALUES ($1)`
	res, err := global.BrewmDB.Exec(sqlStatement, recipeJson)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to insert into Recipe-> "+err.Error())
	}

	// get ID of creates ToDo
	id, err := res.LastInsertId()
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve id for created Recipe-> "+err.Error())
	}

	return &brewmmer.CreateRecipeResponse{
		Id: id,
	}, nil
}

func (s *recipeServiceServer) Get(ctx context.Context, req *brewmmer.GetRecipeRequest) (*brewmmer.GetRecipeResponse, error) {
	sqlStatement := `
    SELECT
      recipe
    FROM recipes
    WHERE id=$1`
	row := global.BrewmDB.QueryRow(sqlStatement, req.Id)

	var recipeJson string
	row.Scan(&recipeJson)

	um := jsonpb.Unmarshaler{}
	unserialized := &brewmmer.Recipe{}
	err := um.Unmarshal(strings.NewReader(recipeJson), unserialized)
	if err != nil {
		return nil, status.Error(codes.Unknown, "json unmarshaling error-> "+err.Error())
	}

	return &brewmmer.GetRecipeResponse{
		Recipe: unserialized,
	}, nil
}

func (s *recipeServiceServer) Delete(ctx context.Context, req *brewmmer.DeleteRecipeRequest) (*brewmmer.DeleteRecipeResponse, error) {
	return nil, nil
}

func (s *recipeServiceServer) Update(ctx context.Context, req *brewmmer.UpdateRecipeRequest) (*brewmmer.UpdateRecipeResponse, error) {
	return nil, nil
}

func (s *recipeServiceServer) List(ctx context.Context, req *brewmmer.ListRecipeRequest) (*brewmmer.ListRecipeResponse, error) {
	recipes := []*brewmmer.Recipe{}
	um := jsonpb.Unmarshaler{}

	rows, err := global.BrewmDB.Query(`
    SELECT
      recipe
    FROM recipes`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {

		var recipeJson *string
		err = rows.Scan(&recipeJson)

		recipe := &brewmmer.Recipe{}
		err := um.Unmarshal(strings.NewReader(*recipeJson), recipe)
		if err != nil {
			return nil, status.Error(codes.Unknown, "json unmarshaling error-> "+err.Error())
		}

		recipes = append(recipes, recipe)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return &brewmmer.ListRecipeResponse{
		Recipes: recipes,
	}, nil
}
