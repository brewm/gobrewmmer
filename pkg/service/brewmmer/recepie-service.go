package brewmmer

import (
	"context"
	"database/sql"
	"strings"

	"github.com/brewm/gobrewmmer/pkg/api/brewmmer"
	"github.com/golang/protobuf/jsonpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type recipeServiceServer struct {
	db *sql.DB
}

func NewRecipeServiceServer(db *sql.DB) brewmmer.RecipeServiceServer {
	return &recipeServiceServer{db: db}
}

func (rss *recipeServiceServer) Create(ctx context.Context, req *brewmmer.CreateRecipeRequest) (*brewmmer.CreateRecipeResponse, error) {
	m := jsonpb.Marshaler{}

	// reset id
	req.Recipe.Id = 0
	recipeJson, err := m.MarshalToString(req.Recipe)
	if err != nil {
		return nil, status.Error(codes.Unknown, "json marshaling error-> "+err.Error())
	}

	sqlStatement := `INSERT INTO recipes (recipe) VALUES ($1)`
	res, err := rss.db.Exec(sqlStatement, recipeJson)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to insert into Recipe-> "+err.Error())
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve id for created Recipe-> "+err.Error())
	}

	return &brewmmer.CreateRecipeResponse{
		Id: id,
	}, nil
}

func (rss *recipeServiceServer) Get(ctx context.Context, req *brewmmer.GetRecipeRequest) (*brewmmer.GetRecipeResponse, error) {
	sqlStatement := `
    SELECT
      recipe
    FROM recipes
    WHERE id=$1`
	row := rss.db.QueryRow(sqlStatement, req.Id)

	var recipeJson string
	row.Scan(&recipeJson)

	um := jsonpb.Unmarshaler{}
	recipe := &brewmmer.Recipe{}
	err := um.Unmarshal(strings.NewReader(recipeJson), recipe)
	if err != nil {
		return nil, status.Error(codes.Unknown, "json unmarshaling error-> "+err.Error())
	}

	// Fix ID
	recipe.Id = req.Id
	return &brewmmer.GetRecipeResponse{
		Recipe: recipe,
	}, nil
}

func (rss *recipeServiceServer) Delete(ctx context.Context, req *brewmmer.DeleteRecipeRequest) (*brewmmer.DeleteRecipeResponse, error) {
	return nil, nil
}

func (rss *recipeServiceServer) Update(ctx context.Context, req *brewmmer.UpdateRecipeRequest) (*brewmmer.UpdateRecipeResponse, error) {
	return nil, nil
}

func (rss *recipeServiceServer) List(ctx context.Context, req *brewmmer.ListRecipeRequest) (*brewmmer.ListRecipeResponse, error) {
	recipes := []*brewmmer.Recipe{}
	um := jsonpb.Unmarshaler{}

	rows, err := rss.db.Query(`
		SELECT
		  id,
      recipe
    FROM recipes`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id *int64
		var recipeJson *string
		err = rows.Scan(
			&id,
			&recipeJson,
		)

		recipe := &brewmmer.Recipe{}
		err := um.Unmarshal(strings.NewReader(*recipeJson), recipe)
		if err != nil {
			return nil, status.Error(codes.Unknown, "json unmarshaling error-> "+err.Error())
		}

		// Fix ID
		recipe.Id = *id
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
