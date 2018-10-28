package temperature

import (
  "fmt"
  "time"
  "strconv"
  "net/http"

  "github.com/gin-gonic/gin"

  "github.com/brewm/gobrewmmer/pkg/ds18b20"
  conn "github.com/brewm/gobrewmmer/pkg/connections"
)

// https://play.golang.org/p/9TSzoxgzF13 for testing
var sessionChannel chan struct{}

type Session struct {
  Id             int              `json:"id"`
  StartTime      time.Time        `json:"startTime"`
  StopTime       time.Time        `json:"stopTime"`
  Measurements   []Measurement    `json:"measurements,omitempty"`
  Note           string           `json:"note"`
}

type Measurement struct {
  SessionId   int       `json:"sessionId,omitempty"`
  Timestamp   time.Time `json:"timestamp"`
  Temperature float64   `json:"temperature"`
}


func Sense(c *gin.Context) {
  m := Measurement{Timestamp: time.Now(), Temperature: ds18b20.ReadTemperature()}
  c.JSON(200, m)
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

  id, err := strconv.Atoi(c.Param("id"))
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

func StartSession(c *gin.Context) {
  if sessionChannel != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": "Session is already in progress. One session can be actibe at a time."})
    return
  }

  note := c.PostForm("note")

  sqlStatement := `
    INSERT INTO session (start_time, note)
    VALUES ($1, $2)`

  result, err := conn.BrewmmerDB.Exec(sqlStatement, time.Now(), note)
  if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
  }

  sessionId, err :=result.LastInsertId()
  if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
  }

  startSessionProcess(int(sessionId))

  c.JSON(http.StatusOK, gin.H{
    "id": sessionId,
  })
}

func startSessionProcess(id int) {
  sessionChannel = make(chan struct{})

  // Start goroutine to periodically run the insert
  go func(id int) {
    for {
      // Start goroutine to do an async insert
      go insertTemperature(id)

      time.Sleep(1 * time.Second)
      select {
      case <-sessionChannel:
        fmt.Printf("[%d] Stop session process\n", id)
        return
      default: // adding default will make it not block
        fmt.Printf("[%d] Rolling to next iteration\n", id)
      }
    }
  }(id)
}

func insertTemperature(id int) {
  sqlStatement := `
    INSERT INTO measurement (session_id, timestamp, temperature)
    VALUES ($1, $2, $3);`

  _, err := conn.BrewmmerDB.Exec(sqlStatement, id, time.Now(), ds18b20.ReadTemperature())
  if err != nil {
    fmt.Printf("ERROR: temperature recording failed for session with id '%d'\n", id)
  }
}

func StopSession(c *gin.Context) {
  if sessionChannel == nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": "There is no active session. Nothing to stop!"})
    return
  }
  close(sessionChannel)
}

