package main

import (
  "os"
  "fmt"
  "log"
  "bytes"
  "net/http"
  "net/url"
  "io/ioutil"

  "github.com/urfave/cli"
)

const version = "0.1"

const endpoint = "http://localhost:8080/"

func main() {
  app := cli.NewApp()
  app.Name = "brewmctl"
  app.Usage = "command line interface to control brewmmer"
  app.Version = version

  app.Commands = []cli.Command{
    {
      Name:  "get",
      Usage: "Get <resource>",
      Subcommands: cli.Commands{
        cli.Command{
          Name:   "temperature",
          Action: getTemperature,
        },
        cli.Command{
          Name:   "sessions",
          Action: getSessions,
        },
      },

    },
    {
      Name:  "start",
      Usage: "start <process>",
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
      Usage: "stop <process>",
      Subcommands: cli.Commands{
        cli.Command{
          Name:   "session",
          ArgsUsage: "<id>",
          Action: stopSession,
        },
      },
    },
  }

  err := app.Run(os.Args)
  if err != nil {
    log.Fatal(err)
  }
}


func getTemperature(c *cli.Context) error {
  data, err := requestWrapper("GET", endpoint + "v1/sense/", nil)
  if err != nil {
    return err
  }

  fmt.Println(data)

  return nil
}


func getSessions(c *cli.Context) error {
  data, err := requestWrapper("GET", endpoint + "v1/sessions/", nil)
  if err != nil {
    return err
  }

  fmt.Println(data)

  return nil
}

func startSession(c *cli.Context) error {
  payload := url.Values{}
  payload.Set("note", c.String("note"))

  data, err := requestWrapper("POST", endpoint + "v1/sessions/", &payload)
  if err != nil {
    return err
  }

  fmt.Println(data)

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

  data, err := requestWrapper("PUT", endpoint + "v1/sessions/", &payload)
  if err != nil {
    return err
  }

  fmt.Println(data)

  return nil
}

func requestWrapper(method string, endpoint string, payload *url.Values) (string, error) {
  var req *http.Request
  var err error

  if payload != nil {
    req, err = http.NewRequest(method, endpoint, bytes.NewBufferString(payload.Encode()))
  } else {
    req, err = http.NewRequest(method, endpoint, nil)
  }

  if err != nil {
    return "", err
  }
  req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")

  res, err := http.DefaultClient.Do(req)
  if err != nil {
    return "", err
  }


  defer res.Body.Close()

  body, err := ioutil.ReadAll(res.Body)
  if err != nil {
    return "", err
  }

  return string(body), nil
}
