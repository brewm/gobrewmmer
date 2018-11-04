package main

import (
  "os"
  "fmt"
  "log"
  "bytes"
  "net/http"
  "net/url"
  "io/ioutil"
  "encoding/json"

  "github.com/urfave/cli"
)

const version = "0.1"

var endpoint string

func main() {
  endpoint = getEnv("BREWM_ENDPOINT", "http://localhost:8080")


  app := cli.NewApp()
  app.Name = "brewmctl"
  app.Usage = "command line interface to control brewmmer"
  app.Version = version

  app.Commands = []cli.Command{
    {
      Name:  "get",
      Usage: "get <temperature|sessions>",
      Subcommands: cli.Commands{
        cli.Command{
          Name:   "temperature",
          Action: getTemperature,
        },
        cli.Command{
          Name:   "sessions",
          Action: getSessions,
          ArgsUsage: "<id>",
          Flags: []cli.Flag{
            cli.BoolFlag{Name: "active"},
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
  data, err := requestWrapper("GET", endpoint + "/v1/sense/", nil)
  if err != nil {
    return err
  }

  fmt.Println(data)

  return nil
}

func getSessions(c *cli.Context) error {
  url := endpoint + "/v1/sessions"

  activeFlag := c.Bool("active")

  switch {
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

  var out bytes.Buffer
  json.Indent(&out, data, "", "\t")
  out.WriteTo(os.Stdout)

  return nil
}

func startSession(c *cli.Context) error {
  payload := url.Values{}
  payload.Set("note", c.String("note"))

  data, err := requestWrapper("POST", endpoint + "/v1/sessions/", &payload)
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

  data, err := requestWrapper("PUT", endpoint + "/v1/sessions/", &payload)
  if err != nil {
    return err
  }

  fmt.Println(string(data))

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

func getEnv(key, fallback string) string {
  if value, ok := os.LookupEnv(key); ok {
    return value
  }
  return fallback
}
