package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/brewm/gobrewmmer/pkg/api/recepie"
	"github.com/brewm/gobrewmmer/pkg/api/session"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/ptypes"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
)

const version = "0.2"

var endpoint string

var rClient recepie.RecepieServiceClient
var sClient session.SessionServiceClient
var ctx context.Context

func main() {

	// Set up a connection to the server.
	conn, err := grpc.Dial("localhost:6999", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("did not connect: %v", err)
	}
	defer conn.Close()

	rClient = recepie.NewRecepieServiceClient(conn)
	sClient = session.NewSessionServiceClient(conn)

	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	endpoint = getEnv("BREWM_ENDPOINT", "http://localhost:8080")

	app := cli.NewApp()
	app.Name = "brewmctl"
	app.Usage = "command line interface to control brewmmer"
	app.Version = version

	app.Commands = []cli.Command{
		{
			Name:  "get",
			Usage: "get <temperature|sessions|recepies>",
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
					Name:      "recepies",
					Action:    getRecepies,
					ArgsUsage: "<id>",
					Flags: []cli.Flag{
						cli.BoolFlag{Name: "raw"},
					},
				},
			},
		},
		{
			Name:  "start",
			Usage: "start <session|tbd...>",
			Subcommands: cli.Commands{
				cli.Command{
					Name:   "session",
					Action: startSession,
					Flags: []cli.Flag{
						cli.StringFlag{Name: "note"},
					},
				},
			},
		},
		{
			Name:  "stop",
			Usage: "stop <session|tbd...>",
			Subcommands: cli.Commands{
				cli.Command{
					Name:      "session",
					ArgsUsage: "<id>",
					Action:    stopSession,
				},
			},
		},
		{
			Name:  "create",
			Usage: "create <recepie|tbd...>",
			Subcommands: cli.Commands{
				cli.Command{
					Name:      "recepie",
					Action:    createRecepie,
					ArgsUsage: "<json>",
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
	data, err := requestWrapper("GET", endpoint+"/v1/sense", nil)
	if err != nil {
		return err
	}

	prettyJsonPrint(data)

	return nil
}

func getSessions(c *cli.Context) error {
	url := endpoint + "/v1/sessions"

	activeFlag := c.Bool("active")

	switch {
	case c.NArg() == 0 && !activeFlag:
		req := session.ListRequest{}
		res, err := sClient.List(ctx, &req)
		if err != nil {
			fmt.Printf("ERROR: Grpc call failed: %v", err)
		}

		prettyPrintSessions(res.Sessions)
		return nil

	case c.NArg() == 0 && activeFlag:
		url += "?active=true"
	case c.NArg() == 1:
		url += "/" + c.Args().Get(0)
	case c.NArg() == 1 && activeFlag:
		fmt.Println("ERROR: Argument and active flag can't be used togather!")
		cli.ShowSubcommandHelp(c)
		return nil
	case c.NArg() > 1:
		fmt.Println("ERROR: Only one argument is accepted!")
		cli.ShowSubcommandHelp(c)
		return nil
	}

	data, err := requestWrapper("GET", url, nil)
	if err != nil {
		return err
	}

	prettyJsonPrint(data)

	return nil
}

func startSession(c *cli.Context) error {
	payload := url.Values{}
	payload.Set("note", c.String("note"))

	data, err := requestWrapper("POST", endpoint+"/v1/sessions/", &payload)
	if err != nil {
		return err
	}

	fmt.Println(string(data))

	return nil
}

func stopSession(c *cli.Context) error {
	if c.NArg() != 1 {
		fmt.Println("ERROR: Missing command argument!")
		cli.ShowSubcommandHelp(c)
		return nil
	}

	id := c.Args().Get(0)

	payload := url.Values{}
	payload.Set("id", id)

	data, err := requestWrapper("PUT", endpoint+"/v1/sessions/", &payload)
	if err != nil {
		return err
	}

	fmt.Println(string(data))

	return nil
}

func createRecepie(c *cli.Context) error {
	if c.NArg() != 1 {
		fmt.Println("ERROR: Missing command argument!")
		cli.ShowSubcommandHelp(c)
		return nil
	}

	json := c.Args().Get(0)

	// Unmarshal json string to Recepie struct
	um := jsonpb.Unmarshaler{}
	unserialized := &recepie.Recepie{}
	err := um.Unmarshal(strings.NewReader(json), unserialized)
	if err != nil {
		fmt.Printf("json unmarshaling error: %v", err)
	}

	// Call Create
	req := recepie.CreateRequest{
		Recepie: unserialized,
	}
	res, err := rClient.Create(ctx, &req)
	if err != nil {
		fmt.Printf("ERROR: Grpc call failed: %v", err)
	}

	fmt.Printf("Created recepie with id: <%+v>", res.Id)

	return nil
}

func getRecepies(c *cli.Context) error {
	rawFlag := c.Bool("raw")

	switch {
	case c.NArg() == 0:
		// Get all recepies
		// TODO: implement call
		fmt.Println("ERROR: Not implemented yet!")
	case c.NArg() == 1:
		// Get recepie with id
		id := c.Args().Get(0)
		return getRecepieWithID(id, rawFlag)
	case c.NArg() > 1:
		fmt.Println("ERROR: Max one argument is accepted!")
		cli.ShowSubcommandHelp(c)
		return nil
	}

	return nil
}

func getRecepieWithID(id string, raw bool) error {
	i, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		fmt.Printf("Failed to parese id as integer!")
		return err
	}

	req := recepie.GetRequest{
		Id: i,
	}
	res, err := rClient.Get(ctx, &req)
	if err != nil {
		fmt.Printf("Grpc call failed.")
		return err
	}

	if raw {
		m := jsonpb.Marshaler{}
		recepieJSON, err := m.MarshalToString(res.Recepie)
		if err != nil {
			fmt.Printf("JSON marshaling failed.")
			return err
		}

		fmt.Printf(recepieJSON)
	} else {
		prettyPrintRecepie(res.Recepie)
	}

	return nil
}

func requestWrapper(method string, url string, payload *url.Values) ([]byte, error) {
	var req *http.Request
	var err error

	if payload != nil {
		req, err = http.NewRequest(method, url, bytes.NewBufferString(payload.Encode()))
	} else {
		req, err = http.NewRequest(method, url, nil)
	}

	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func prettyJsonPrint(data []byte) {
	var out bytes.Buffer
	json.Indent(&out, data, "", "\t")
	out.WriteTo(os.Stdout)
}

func prettyPrintSessions(ss []*session.Session) {
	fmt.Println("ID", "\t", "StartTime", "\t", "StopTime", "\t", "Note")
	for _, s := range ss {
		fmt.Println(s.Id, "\t", ptypes.TimestampString(s.StartTime), "\t", ptypes.TimestampString(s.StopTime), "\t", s.Note)
	}
}

func prettyPrintRecepie(r *recepie.Recepie) {
	var quantityToStr func(*recepie.Quantity) string
	quantityToStr = func(q *recepie.Quantity) string {
		if q != nil {
			return fmt.Sprint(q.Volume, " ", q.Unit)
		}
		return ""
	}

	fmt.Println("\tName: ", r.Name)
	fmt.Println("\tDescription: ", r.Description)
	fmt.Println("\tIngredients: ")
	fmt.Println("\t\t", "Type", "\t", "Name", "\t", "Quantity")
	for _, i := range r.Ingredients {
		fmt.Println("\t\t", i.Type, "\t", i.Name, "\t", quantityToStr(i.Quantity))
	}
	fmt.Println("\tSteps: ")
	fmt.Println("\t\t", "Phase", "\t", "Temperature", "\t", "Duration")
	for _, s := range r.Steps {
		fmt.Println("\t\t", s.Phase, "\t", s.Temperature, "\t", quantityToStr(s.Duration))
		if len(s.Ingredients) != 0 {
			fmt.Println("\t\t\t", "Type", "\t", "Name", "\t", "Quantity")
		}
		for _, si := range s.Ingredients {
			fmt.Println("\t\t\t", si.Type, "\t", si.Name, "\t", quantityToStr(si.Quantity))
		}
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
