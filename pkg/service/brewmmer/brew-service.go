package brewmmer

import (
	"context"
	"database/sql"
	"errors"
	"reflect"
	"strings"
	"time"

	"github.com/brewm/gobrewmmer/pkg/api/brewmmer"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type brewServiceServer struct {
	db *sql.DB
}

func NewBrewServiceServer(db *sql.DB) brewmmer.BrewServiceServer {
	return &brewServiceServer{db: db}
}

func (bss *brewServiceServer) StartBrew(ctx context.Context, req *brewmmer.StartBrewRequest) (*brewmmer.StartBrewResponse, error) {
	timestamp := time.Now()

	sqlStatement := `
    INSERT INTO brews (recipe_id, start_time, note)
    VALUES ($1, $2, $3)`

	result, err := bss.db.Exec(sqlStatement, req.RecipeId, timestamp, req.Note)
	if err != nil {
		return nil, err
	}

	brewID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &brewmmer.StartBrewResponse{
		Id: brewID,
	}, nil
}

func (bss *brewServiceServer) CompleteBrewStep(ctx context.Context, req *brewmmer.CompleteBrewStepRequest) (*brewmmer.CompleteBrewStepResponse, error) {
	brew, err := fetchActiveBrew(bss.db)
	if err != nil {
		return nil, err
	}

	if brew != nil {
		return nil, errors.New("No brew is active atm")
	}

	completedSteps := brew.CompletedSteps
	completedStepCount := len(completedSteps)
	var nextStep *brewmmer.Step

	if completedStepCount == 0 {
		nextStep = brew.Recipe.Steps[0]
	} else {
		lastStep := completedSteps[completedStepCount-1]
		nextStep = getNextStep(lastStep.Step, brew.Recipe)
	}

	m := jsonpb.Marshaler{}
	stepJSON, err := m.MarshalToString(nextStep)
	if err != nil {
		return nil, status.Error(codes.Unknown, "json marshaling error-> "+err.Error())
	}

	sqlStatement := `INSERT INTO brew_steps (start_time, step, brew_id) VALUES ($1, $2, $3)`
	_, err = bss.db.Exec(sqlStatement, time.Now(), stepJSON, req.Id)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to insert into brew_steps-> "+err.Error())
	}

	return &brewmmer.CompleteBrewStepResponse{
		NextStep: nextStep,
	}, nil
}

func getNextStep(lastStep *brewmmer.Step, recipe *brewmmer.Recipe) *brewmmer.Step {
	for i, step := range recipe.Steps {
		if reflect.DeepEqual(lastStep, step) {
			return recipe.Steps[i] // not sure about this
		}
	}
	return nil
}

func (bss *brewServiceServer) StopBrew(ctx context.Context, req *brewmmer.StopBrewRequest) (*brewmmer.StopBrewResponse, error) {
	sqlStatement := `
    SELECT (CASE WHEN completed_time IS NULL THEN 1 ELSE 0 END) as is_active
    FROM brews
    WHERE id = $1`
	row := bss.db.QueryRow(sqlStatement, req.Id)

	var isActive bool
	err := row.Scan(&isActive)

	if err != nil {
		log.WithFields(log.Fields{
			"id":  req.Id,
			"err": err,
		}).Error("Checking brew with the given id failed!")
		return nil, err
	}

	if isActive == false {
		return nil, errors.New("given brew is not active, can't stop")
	}

	sqlStatement = `
    UPDATE brews
    SET completed_time = $1
    WHERE id = $2`

	timestamp := time.Now()

	_, err = bss.db.Exec(sqlStatement, timestamp, req.Id)
	if err != nil {
		log.WithFields(log.Fields{
			"id":  req.Id,
			"err": err,
		}).Error("Update completed time for brew database query failed!")

		return nil, err
	}

	return &brewmmer.StopBrewResponse{}, nil
}

func (bss *brewServiceServer) GetBrew(ctx context.Context, req *brewmmer.GetBrewRequest) (*brewmmer.GetBrewResponse, error) {
	log.WithFields(log.Fields{"id": req.Id}).Info("Getting brew!")

	brew, err := fetchBrew(bss.db, req.Id)
	if err != nil {
		return nil, err
	}

	return &brewmmer.GetBrewResponse{
		Brew: brew,
	}, nil
}

func fetchBrew(db *sql.DB, id int64) (*brewmmer.Brew, error) {
	sqlStatement := `
	SELECT
		id,
		recipe_id,
		start_time,
		completed_time,
		note
	FROM brews
	WHERE id=$1`
	row := db.QueryRow(sqlStatement, id)

	var startTime *time.Time
	var nullableCompletedTime *time.Time
	var recipeId *int64

	brew := new(brewmmer.Brew)
	err := row.Scan(
		&brew.Id,
		&recipeId,
		&startTime,
		&nullableCompletedTime,
		&brew.Note,
	)
	if err != nil {
		log.WithFields(log.Fields{
			"id":  id,
			"err": err,
		}).Error("Find id in 'brews' database query failed!")
		return nil, err
	}

	if nullableCompletedTime != nil {
		brew.CompletedTime, err = ptypes.TimestampProto(*nullableCompletedTime)
		if err != nil {
			return nil, err
		}
	} else {
		brew.CompletedTime = new(timestamp.Timestamp)
	}

	brew.StartTime, err = ptypes.TimestampProto(*startTime)
	if err != nil {
		return nil, err
	}

	log.WithFields(log.Fields{
		"id":             brew.Id,
		"recipe_id":      recipeId,
		"start_time":     brew.StartTime,
		"completed_time": brew.CompletedTime,
	}).Debug("Found brew!")

	recipe, err := fetchRecipe(db, *recipeId)
	if err != nil {
		return nil, err
	}
	brew.Recipe = recipe

	err = fetchCompletedSteps(db, brew)
	if err != nil {
		return nil, err
	}

	return brew, nil
}

func (bss *brewServiceServer) GetActiveBrew(ctx context.Context, req *brewmmer.GetActiveBrewRequest) (*brewmmer.GetActiveBrewResponse, error) {
	log.Info("Getting active brew!")
	brew, err := fetchActiveBrew(bss.db)

	return &brewmmer.GetActiveBrewResponse{
		Brew: brew,
	}, err
}

func fetchActiveBrew(db *sql.DB) (*brewmmer.Brew, error) {
	// Making sure to return only one session
	sqlStatement := `
	SELECT
		MAX(id)
	FROM brews
	WHERE completed_time IS NULL`
	row := db.QueryRow(sqlStatement)

	var id *int64

	err := row.Scan(&id)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("Find latest brew in 'brews' database query failed!")

		return nil, err
	}

	if id == nil {
		return nil, status.Error(codes.Unknown, "No actibe brew was found!")
	}

	return fetchBrew(db, *id)
}

func (bss *brewServiceServer) ListBrews(ctx context.Context, req *brewmmer.ListBrewRequest) (*brewmmer.ListBrewResponse, error) {
	panic("not implemented")
}

func fetchCompletedSteps(db *sql.DB, brew *brewmmer.Brew) error {
	log.WithFields(log.Fields{"brew_id": brew.Id}).Debug("Getting completed steps!")

	sqlStatement := `
    SELECT
      start_time,
			completed_time,
			step
    FROM brew_steps
    WHERE brew_id=$1`

	rows, err := db.Query(sqlStatement, brew.Id)

	if err != nil {
		log.WithFields(log.Fields{
			"session_id": brew.Id,
			"err":        err,
		}).Error("Get completed steps from 'brew_steps' database query failed!")

		return err
	}
	defer rows.Close()

	for rows.Next() {
		var startTime *time.Time
		var nullableCompletedTime *time.Time
		var stepJSON *string

		bs := new(brewmmer.BrewStep)
		err = rows.Scan(
			&startTime,
			&nullableCompletedTime,
			&stepJSON,
		)

		bs.StartTime, err = ptypes.TimestampProto(*startTime)
		if err != nil {
			return err
		}

		if nullableCompletedTime != nil {
			bs.CompletedTime, err = ptypes.TimestampProto(*nullableCompletedTime)
			if err != nil {
				return err
			}
		} else {
			bs.CompletedTime = new(timestamp.Timestamp)
		}

		um := jsonpb.Unmarshaler{}
		step := &brewmmer.Step{}
		err := um.Unmarshal(strings.NewReader(*stepJSON), step)
		if err != nil {
			return status.Error(codes.Unknown, "json unmarshaling error-> "+err.Error())
		}

		bs.Step = step

		brew.CompletedSteps = append(brew.CompletedSteps, bs)
	}
	err = rows.Err()
	if err != nil {
		return err
	}

	return nil
}
