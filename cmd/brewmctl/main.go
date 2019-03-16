package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/brewm/gobrewmmer/pkg/api/brewmmer"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
)

const version = "0.2"

var endpoint string

var rClient brewmmer.RecipeServiceClient
var bClient brewmmer.BrewServiceClient
var sClient brewmmer.SessionServiceClient
var tClient brewmmer.TemperatureServiceClient
var ctx context.Context

func main() {

	// Set up a connection to the server.
	endpoint = getEnv("BREWM_ENDPOINT", "localhost:6999")

	conn, err := grpc.Dial(endpoint, grpc.WithInsecure())
	if err != nil {
		fmt.Printf("did not connect: %v", err)
	}
	defer conn.Close()

	rClient = brewmmer.NewRecipeServiceClient(conn)
	bClient = brewmmer.NewBrewServiceClient(conn)
	sClient = brewmmer.NewSessionServiceClient(conn)
	tClient = brewmmer.NewTemperatureServiceClient(conn)

	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	app := cli.NewApp()
	app.Name = "brewmctl"
	app.Usage = "command line interface to control brewmmer"
	app.Version = version

	app.Commands = []cli.Command{
		{
			Name:  "get",
			Usage: "get <temperature|sessions|recipes|brews>",
			Subcommands: cli.Commands{
				cli.Command{
					Name:   "temperature",
					Action: getTemperature,
				},
				cli.Command{
					Name:      "sessions",
					Action:    getSessions,
					ArgsUsage: "<id>",
					Flags: []cli.Flag{
						cli.BoolFlag{Name: "active"},
					},
				},
				cli.Command{
					Name:      "recipes",
					Action:    getRecipes,
					ArgsUsage: "<id>",
					Flags: []cli.Flag{
						cli.BoolFlag{Name: "raw"},
					},
				},
				cli.Command{
					Name:      "brews",
					Action:    getBrews,
					ArgsUsage: "<id>",
					Flags: []cli.Flag{
						cli.BoolFlag{Name: "active"},
					},
				},
			},
		},
		{
			Name:  "start",
			Usage: "start <session|brew|tbd...>",
			Subcommands: cli.Commands{
				cli.Command{
					Name:   "session",
					Action: startSession,
					Flags: []cli.Flag{
						cli.StringFlag{Name: "note"},
					},
				},
				cli.Command{
					Name:      "brew",
					ArgsUsage: "<recipe_id>",
					Action:    startBrew,
					Flags: []cli.Flag{
						cli.StringFlag{Name: "note"},
					},
				},
			},
		},
		{
			Name:  "stop",
			Usage: "stop <session|brew|tbd...>",
			Subcommands: cli.Commands{
				cli.Command{
					Name:      "session",
					ArgsUsage: "<id>",
					Action:    stopSession,
				},
				cli.Command{
					Name:      "brew",
					ArgsUsage: "<id>",
					Action:    stopBrew,
				},
			},
		},
		{
			Name:  "create",
			Usage: "create <recipe|tbd...>",
			Subcommands: cli.Commands{
				cli.Command{
					Name:      "recipe",
					Action:    createRecipe,
					ArgsUsage: "<json>",
				},
			},
		},
		{
			Name:  "complete",
			Usage: "complete <step|tbd...>",
			Subcommands: cli.Commands{
				cli.Command{
					Name:      "step",
					Action:    completeStep,
					ArgsUsage: "<brew_id>",
				},
			},
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func getTemperature(c *cli.Context) error {
	req := brewmmer.GetTemperatureRequest{}
	res, err := tClient.Get(ctx, &req)
	if err != nil {
		fmt.Printf("ERROR: Grpc call failed: %v", err)
		return err
	}

	fmt.Println("Current temperature: ", res.Temperature)

	return nil
}

func getSessions(c *cli.Context) error {
	activeFlag := c.Bool("active")

	switch {
	case c.NArg() == 0 && !activeFlag:
		req := brewmmer.ListSessionRequest{}
		res, err := sClient.List(ctx, &req)
		if err != nil {
			fmt.Printf("ERROR: Grpc call failed: %v", err)
			return err
		}
		prettyPrintSessions(res.Sessions)

	case c.NArg() == 0 && activeFlag:
		req := brewmmer.GetActiveSessionRequest{}
		res, err := sClient.GetActive(ctx, &req)
		if err != nil {
			fmt.Printf("ERROR: Grpc call failed: %v", err)
			return err
		}
		prettyPrintSession(res.Session)

	case c.NArg() == 1:
		id, err := strconv.ParseInt(c.Args().Get(0), 10, 64)
		if err != nil {
			fmt.Printf("Failed to parese id as integer!")
			return err
		}
		req := brewmmer.GetSessionRequest{
			Id: id,
		}
		res, err := sClient.Get(ctx, &req)
		if err != nil {
			fmt.Printf("ERROR: Grpc call failed: %v", err)
			return err
		}
		prettyPrintSession(res.Session)

	case c.NArg() == 1 && activeFlag:
		fmt.Println("ERROR: Argument and active flag can't be used together!")
		cli.ShowSubcommandHelp(c)

	case c.NArg() > 1:
		fmt.Println("ERROR: Only one argument is accepted!")
		cli.ShowSubcommandHelp(c)
	}

	return nil
}

func startSession(c *cli.Context) error {
	req := brewmmer.StartSessionRequest{}
	res, err := sClient.Start(ctx, &req)
	if err != nil {
		fmt.Printf("ERROR: Grpc call failed: %v", err)
		return err
	}

	fmt.Println("Session created with id: ", res.Id)

	return nil
}

func stopSession(c *cli.Context) error {
	if c.NArg() != 1 {
		fmt.Println("ERROR: Missing command argument!")
		cli.ShowSubcommandHelp(c)
		return nil
	}

	id, err := strconv.ParseInt(c.Args().Get(0), 10, 64)
	if err != nil {
		fmt.Printf("Failed to parese id as integer!")
		return err
	}

	req := brewmmer.StopSessionRequest{
		Id: id,
	}
	_, err = sClient.Stop(ctx, &req)
	if err != nil {
		fmt.Printf("ERROR: Grpc call failed: %v", err)
		return err
	}

	fmt.Println("Session stopped")

	return nil
}

func getBrews(c *cli.Context) error {
	activeFlag := c.Bool("active")

	switch {
	case c.NArg() == 0 && !activeFlag:
		req := brewmmer.ListBrewRequest{}
		_, err := bClient.ListBrews(ctx, &req)
		if err != nil {
			fmt.Printf("ERROR: Grpc call failed: %v", err)
			return err
		}
		// prettyPrintBrews(res.Brews)

	case c.NArg() == 0 && activeFlag:
		req := brewmmer.GetActiveBrewRequest{}
		res, err := bClient.GetActiveBrew(ctx, &req)
		if err != nil {
			fmt.Printf("ERROR: Grpc call failed: %v", err)
			return err
		}
		prettyPrintBrew(res.Brew)

	case c.NArg() == 1:
		id, err := strconv.ParseInt(c.Args().Get(0), 10, 64)
		if err != nil {
			fmt.Printf("Failed to parese id as integer!")
			return err
		}
		req := brewmmer.GetBrewRequest{
			Id: id,
		}
		res, err := bClient.GetBrew(ctx, &req)
		if err != nil {
			fmt.Printf("ERROR: Grpc call failed: %v", err)
			return err
		}
		prettyPrintBrew(res.Brew)

	case c.NArg() == 1 && activeFlag:
		fmt.Println("ERROR: Argument and active flag can't be used together!")
		cli.ShowSubcommandHelp(c)

	case c.NArg() > 1:
		fmt.Println("ERROR: Only one argument is accepted!")
		cli.ShowSubcommandHelp(c)
	}

	return nil
}

func startBrew(c *cli.Context) error {
	if c.NArg() != 1 {
		fmt.Println("ERROR: Missing command argument!")
		cli.ShowSubcommandHelp(c)
		return nil
	}

	recipeId, err := strconv.ParseInt(c.Args().Get(0), 10, 64)
	if err != nil {
		fmt.Printf("Failed to parese recipe id as integer!")
		return err
	}

	req := brewmmer.StartBrewRequest{
		RecipeId: recipeId,
	}

	res, err := bClient.StartBrew(ctx, &req)
	if err != nil {
		fmt.Printf("ERROR: Grpc call failed: %v", err)
		return err
	}

	fmt.Println("Brew created with id: ", res.Id)

	return nil
}

func stopBrew(c *cli.Context) error {
	if c.NArg() != 1 {
		fmt.Println("ERROR: Missing command argument!")
		cli.ShowSubcommandHelp(c)
		return nil
	}

	id, err := strconv.ParseInt(c.Args().Get(0), 10, 64)
	if err != nil {
		fmt.Printf("Failed to parese id as integer!")
		return err
	}

	req := brewmmer.StopBrewRequest{
		Id: id,
	}
	_, err = bClient.StopBrew(ctx, &req)
	if err != nil {
		fmt.Printf("ERROR: Grpc call failed: %v", err)
		return err
	}

	fmt.Println("Session stopped")

	return nil
}

func completeStep(c *cli.Context) error {
	if c.NArg() != 1 {
		fmt.Println("ERROR: Missing command argument!")
		cli.ShowSubcommandHelp(c)
		return nil
	}

	id, err := strconv.ParseInt(c.Args().Get(0), 10, 64)
	if err != nil {
		fmt.Printf("Failed to parese id as integer!")
		return err
	}

	req := brewmmer.CompleteBrewStepRequest{
		Id: id,
	}
	res, err := bClient.CompleteBrewStep(ctx, &req)
	if err != nil {
		fmt.Printf("ERROR: Grpc call failed: %v", err)
		return err
	}

	fmt.Println("Step completed, next Step:")
	fmt.Println("\t", "Phase", "\t", "Temperature", "\t", "Duration")
	prettyPrintStep(res.NextStep)

	return nil
}

func createRecipe(c *cli.Context) error {
	if c.NArg() != 1 {
		fmt.Println("ERROR: Missing command argument!")
		cli.ShowSubcommandHelp(c)
		return nil
	}

	json := c.Args().Get(0)

	// Unmarshal json string to Recipe struct
	um := jsonpb.Unmarshaler{}
	unserialized := &brewmmer.Recipe{}
	err := um.Unmarshal(strings.NewReader(json), unserialized)
	if err != nil {
		fmt.Printf("json unmarshaling error: %v", err)
	}

	// Call Create
	req := brewmmer.CreateRecipeRequest{
		Recipe: unserialized,
	}
	res, err := rClient.Create(ctx, &req)
	if err != nil {
		fmt.Printf("ERROR: Grpc call failed: %v", err)
	}

	fmt.Printf("Created recipe with id: <%+v>", res.Id)

	return nil
}

func getRecipes(c *cli.Context) error {
	rawFlag := c.Bool("raw")

	switch {
	case c.NArg() == 0:
		// Get all recipes
		req := brewmmer.ListRecipeRequest{}
		res, err := rClient.List(ctx, &req)
		if err != nil {
			fmt.Printf("Grpc call failed.")
			return err
		}
		prettyPrintRecipes(res.Recipes)
	case c.NArg() == 1:
		// Get recipe with id
		id := c.Args().Get(0)
		return getRecipeWithID(id, rawFlag)
	case c.NArg() > 1:
		fmt.Println("ERROR: Max one argument is accepted!")
		cli.ShowSubcommandHelp(c)
		return nil
	}

	return nil
}

func getRecipeWithID(id string, raw bool) error {
	i, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		fmt.Printf("Failed to parese id as integer!")
		return err
	}

	req := brewmmer.GetRecipeRequest{
		Id: i,
	}
	res, err := rClient.Get(ctx, &req)
	if err != nil {
		fmt.Printf("Grpc call failed.")
		return err
	}

	if raw {
		m := jsonpb.Marshaler{}
		recipeJSON, err := m.MarshalToString(res.Recipe)
		if err != nil {
			fmt.Printf("JSON marshaling failed.")
			return err
		}

		fmt.Printf(recipeJSON)
	} else {
		prettyPrintRecipe(res.Recipe)
	}

	return nil
}

//
// Helpers
//

func prettyPrintSessions(ss []*brewmmer.Session) {
	var maybeTimestamp func(*timestamp.Timestamp) string
	maybeTimestamp = func(t *timestamp.Timestamp) string {
		if t != nil {
			return getTimeString(t)
		}
		return "-"
	}

	fmt.Println("ID", "\t", "StartTime", "\t", "StopTime", "\t", "Note")
	for _, s := range ss {
		fmt.Println(s.Id, "\t", getTimeString(s.StartTime), "\t", maybeTimestamp(s.StopTime), "\t", s.Note)
	}
}

func prettyPrintSession(s *brewmmer.Session) {
	fmt.Println("Note: ", s.Note)
	fmt.Println("Start Time: ", getTimeString(s.StartTime))
	if s.StopTime != nil {
		fmt.Println("Stop Time: ", getTimeString(s.StopTime))
	} else {
		fmt.Println("Stop Time: ", "-")
	}
	fmt.Println("Measurements: ")

	fmt.Println("\t", "Timestamp", "\t", "Temperature")
	for _, m := range s.Measurements {
		fmt.Println("\t", getTimeString(m.Timestamp), "\t", m.Temperature)
	}
}

func prettyPrintRecipes(rs []*brewmmer.Recipe) {
	fmt.Println("ID", "\t", "Name", "\t", "Description")
	for _, r := range rs {
		fmt.Println(r.Id, "\t", r.Name, "\t", r.Description)
	}
}

func prettyPrintRecipe(r *brewmmer.Recipe) {
	fmt.Println("Name: ", r.Name)
	fmt.Println("Description: ", r.Description)
	fmt.Println("Ingredients: ")
	fmt.Println("\t", "Type", "\t", "Name", "\t", "Quantity")
	for _, i := range r.Ingredients {
		fmt.Println("\t", i.Type, "\t", i.Name, "\t", quantityToStr(i.Quantity))
	}
	fmt.Println("Steps: ")
	fmt.Println("\t", "Phase", "\t", "Temperature", "\t", "Duration")
	for _, s := range r.Steps {
		prettyPrintStep(s)
	}
}

func prettyPrintBrew(r *brewmmer.Brew) {
	fmt.Println("Id: ", r.Id)
	fmt.Println("Note: ", r.Note)
	fmt.Println("StartTime: ", getTimeString(r.StartTime))
	fmt.Println("CompletedTime: ", getTimeString(r.CompletedTime))
	fmt.Println("Active Step: ")
	fmt.Println("\t", "StartTime", "\t", "CompletedTime", "\t", "Phase", "\t", "Temperature", "\t", "Duration")
	prettyPrintBrewStep(r.CompletedSteps[len(r.CompletedSteps)-1])

	fmt.Println("Finished Steps: ")
	fmt.Println("\t", "StartTime", "\t", "CompletedTime", "\t", "Phase", "\t", "Temperature", "\t", "Duration")
	for _, step := range r.CompletedSteps[:len(r.CompletedSteps)-1] {
		if step.CompletedTime != nil {
			prettyPrintBrewStep(step)
		}
	}
}

func prettyPrintBrewStep(step *brewmmer.BrewStep) {
	fmt.Println("\t", getTimeString(step.StartTime), "\t", getTimeString(step.CompletedTime), "\t", step.Step.Phase, "\t", step.Step.Temperature, "\t", quantityToStr(step.Step.Duration))
	if len(step.Step.Ingredients) != 0 {
		fmt.Println("\t\t", "Type", "\t", "Name", "\t", "Quantity")
	}
	for _, si := range step.Step.Ingredients {
		fmt.Println("\t\t", si.Type, "\t", si.Name, "\t", quantityToStr(si.Quantity))
	}
}

func prettyPrintStep(s *brewmmer.Step) {
	fmt.Println("\t", s.Phase, "\t", s.Temperature, "\t", quantityToStr(s.Duration))
	if len(s.Ingredients) != 0 {
		fmt.Println("\t\t", "Type", "\t", "Name", "\t", "Quantity")
	}
	for _, si := range s.Ingredients {
		fmt.Println("\t\t", si.Type, "\t", si.Name, "\t", quantityToStr(si.Quantity))
	}
}

func quantityToStr(q *brewmmer.Quantity) string {
	if q != nil {
		return fmt.Sprint(q.Volume, " ", q.Unit)
	}
	return ""
}

func getTimeString(ts *timestamp.Timestamp) string {
	t, err := ptypes.Timestamp(ts)
	if err != nil {
		return fmt.Sprintf("(%v)", err)
	}
	if t.Year() == 1970 {
		return "<n/a>"
	}
	return t.Format(time.RFC3339)
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
