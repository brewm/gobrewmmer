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

	"github.com/brewm/gobrewmmer/pkg/api/brewmmer"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
)

const version = "0.2"

var endpoint string

var rClient brewmmer.RecepieServiceClient
var sClient brewmmer.SessionServiceClient
var ctx context.Context

func main() {

	// Set up a connection to the server.
	conn, err := grpc.Dial("localhost:6999", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("did not connect: %v", err)
	}
	defer conn.Close()

	rClient = brewmmer.NewRecepieServiceClient(conn)
	sClient = brewmmer.NewSessionServiceClient(conn)

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
		fmt.Println("ERROR: Argument and active flag can't be used togather!")
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

func createRecepie(c *cli.Context) error {
	if c.NArg() != 1 {
		fmt.Println("ERROR: Missing command argument!")
		cli.ShowSubcommandHelp(c)
		return nil
	}

	json := c.Args().Get(0)

	// Unmarshal json string to Recepie struct
	um := jsonpb.Unmarshaler{}
	unserialized := &brewmmer.Recepie{}
	err := um.Unmarshal(strings.NewReader(json), unserialized)
	if err != nil {
		fmt.Printf("json unmarshaling error: %v", err)
	}

	// Call Create
	req := brewmmer.CreateRecepieRequest{
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

	req := brewmmer.GetRecepieRequest{
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

func prettyPrintSessions(ss []*brewmmer.Session) {
	var maybeTimestamp func(*timestamp.Timestamp) string
	maybeTimestamp = func(t *timestamp.Timestamp) string {
		if t != nil {
			return ptypes.TimestampString(t)
		}
		return "-"
	}

	fmt.Println("ID", "\t", "StartTime", "\t", "StopTime", "\t", "Note")
	for _, s := range ss {
		fmt.Println(s.Id, "\t", ptypes.TimestampString(s.StartTime), "\t", maybeTimestamp(s.StopTime), "\t", s.Note)
	}
}

func prettyPrintSession(s *brewmmer.Session) {
	fmt.Println("Note: ", s.Note)
	fmt.Println("Start Time: ", ptypes.TimestampString(s.StartTime))
	if s.StopTime != nil {
		fmt.Println("Stop Time: ", ptypes.TimestampString(s.StopTime))
	} else {
		fmt.Println("Stop Time: ", "-")
	}
	fmt.Println("Measurements: ")

	fmt.Println("\t", "Timestamp", "\t", "Temperature")
	for _, m := range s.Measurements {
		fmt.Println("\t", ptypes.TimestampString(m.Timestamp), "\t", m.Temperature)
	}
}

func prettyPrintRecepie(r *brewmmer.Recepie) {
	var quantityToStr func(*brewmmer.Quantity) string
	quantityToStr = func(q *brewmmer.Quantity) string {
		if q != nil {
			return fmt.Sprint(q.Volume, " ", q.Unit)
		}
		return ""
	}

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
		fmt.Println("\t", s.Phase, "\t", s.Temperature, "\t", quantityToStr(s.Duration))
		if len(s.Ingredients) != 0 {
			fmt.Println("\t\t", "Type", "\t", "Name", "\t", "Quantity")
		}
		for _, si := range s.Ingredients {
			fmt.Println("\t\t", si.Type, "\t", si.Name, "\t", quantityToStr(si.Quantity))
		}
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
