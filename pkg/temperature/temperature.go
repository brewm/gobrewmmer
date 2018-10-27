package temperature

import (
  "time"
  "fmt"
  "strings"
  "strconv"
  "net/http"
  "io/ioutil"

  "github.com/gin-gonic/gin"

  conn "github.com/brewm/gobrewmmer/pkg/connections"

  // "periph.io/x/periph/host"
  // "periph.io/x/periph/conn/onewire/onewirereg"
  // "periph.io/x/periph/devices/ds18b20"
  )

const sensorId = "28-0000051e015b"

type Session struct {
  Id             int64            `json:"id"`
  StartTime      time.Time        `json:"startTime"`
  StopTime       time.Time        `json:"stopTime"`
  Measurements   []Measurement    `json:"measurements,omitempty"`
  Note           string           `json:"note"`
}

type Measurement struct {
  SessionId   int64     `json:"sessionId,omitempty"`
  Timestamp   time.Time `json:"timestamp"`
  Temperature float64   `json:"temperature"`
}


func Sense(c *gin.Context) {
  m := Measurement{Timestamp: time.Now(), Temperature: readTemperature()}
  c.JSON(200, m)
}

func readTemperature() float64 {
  sensorPath := fmt.Sprintf("/sys/bus/w1/devices/%s/w1_slave", sensorId)
  data, err := ioutil.ReadFile(sensorPath)
  if err != nil {
    fmt.Println(err)
    return 0.0
  }

  raw := string(data)
  tString := raw[strings.LastIndex(raw, "t=")+2:len(raw)-1]

  t, err := strconv.ParseFloat(tString, 64)
  if err != nil {
    fmt.Println(err)
    return 0.0
  }

  return t / 1000.0
}

func readTestTemperature() float64 {
  return 21.3
}


func AllSession(c *gin.Context) {
  sessions := []Session{}

  if err := fetchAllSession(&sessions); err == nil {
    c.JSON(http.StatusOK, sessions)
  } else {
    c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
  }
}

func fetchAllSession(sessions *[]Session) error {
  rows, err := conn.BrewmmerDB.Query(`
    SELECT
      id,
      start_time,
      stop_time,
      note
    FROM session`)

  if err != nil {
    return err
  }
  defer rows.Close()

  for rows.Next() {
    // The sqlite driver can't handle nullable Time type so here is the workaround
    var nullableStopTime *time.Time

    s := Session{}
    err = rows.Scan(
      &s.Id,
      &s.StartTime,
      &nullableStopTime,
      &s.Note,
    )

    if err != nil {
      return err
    }

    if nullableStopTime != nil {
      s.StopTime = *nullableStopTime
    } else {
      s.StopTime = *new(time.Time)
    }

    *sessions = append(*sessions, s)
  }
  err = rows.Err()
  if err != nil {
    return err
  }

  return nil
}


func SingleSession(c *gin.Context) {
  session := Session{}

  id, err := strconv.ParseInt(c.Param("id"), 10, 64)
  if err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
  }

  session.Id = id

  if err := fetchSingleSession(&session); err == nil {
    c.JSON(http.StatusOK, session)
  } else {
    c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
  }
}

func fetchSingleSession(session *Session) error {
  sqlStatement := `
    SELECT
      id,
      start_time,
      stop_time,
      note
    FROM session
    WHERE id=$1`
  row := conn.BrewmmerDB.QueryRow(sqlStatement, session.Id)

  var nullableStopTime *time.Time

  err := row.Scan(
    &session.Id,
    &session.StartTime,
    &nullableStopTime,
    &session.Note,
  )

  if err != nil {
    return err
  }

  if nullableStopTime != nil {
    session.StopTime = *nullableStopTime
  } else {
    session.StopTime = *new(time.Time)
  }

  err = fetchMeasurements(session)
  if err != nil {
    return err
  }

  return nil
}

func fetchMeasurements(session *Session) error {
  sqlStatement := `
    SELECT
      timestamp,
      temperature
    FROM measurement
    WHERE session_id=$1`

  rows, err := conn.BrewmmerDB.Query(sqlStatement, session.Id)

  if err != nil {
    return err
  }
  defer rows.Close()

  for rows.Next() {
    m := Measurement{}
    err = rows.Scan(
      &m.Timestamp,
      &m.Temperature,
    )

    if err != nil {
      return err
    }

    session.Measurements = append(session.Measurements, m)
  }
  err = rows.Err()
  if err != nil {
    return err
  }

  return nil
}

// func FetchSession(db *sql.DB, sessionId int) Session {

// }

// func recordMeasurement(db *sql.DB) {
//   _, err = db.Exec("insert into foo(id, name) values(1, 'foo'), (2, 'bar'), (3, 'baz')")
//   if err != nil {
//     fmt.Println(err)
//   }
// }



// This is not working. :(
// func Temperature(){
//   // Make sure periph is initialized.
//   if _, err := host.Init(); err != nil {
//     fmt.Println(err)
//   }

//   // Use onewirereg 1-wire bus registry to find the first available 1-wire bus.
//   bus, err := onewirereg.Open("")
//   if err != nil {
//     fmt.Println(err)
//   }
//   defer bus.Close()

//   fmt.Println(ds18b20.ConvertAll(bus, 10))
//   // 0x7a00000131825228
//   dev, err := ds18b20.New(bus, 0x7a0000051e015b, 10)
//   if err != nil {
//       fmt.Println(err)
//   }
//   defer dev.Halt()

//   fmt.Println(dev.LastTemp())
// }


