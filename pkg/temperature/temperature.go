package temperature

import (
  "time"
  "strconv"
  "net/http"

  "github.com/gin-gonic/gin"

  "github.com/brewm/gobrewmmer/pkg/ds18b20"
  conn "github.com/brewm/gobrewmmer/pkg/connections"
)

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
