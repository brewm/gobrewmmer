package temperature

import (
  "time"
  "strconv"
  "net/http"

  "github.com/gin-gonic/gin"
  log "github.com/sirupsen/logrus"

  "github.com/brewm/gobrewmmer/pkg/ds18b20"
  global "github.com/brewm/gobrewmmer/pkg/global"
)

var sessionChannel chan struct{}

// in seconds
const measureInterval = 600

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
  rows, err := global.BrewmDB.Query(`
    SELECT
      id,
      start_time,
      stop_time,
      note
    FROM sessions`)

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
    FROM sessions
    WHERE id=$1`
  row := global.BrewmDB.QueryRow(sqlStatement, session.Id)

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
    FROM measurements
    WHERE session_id=$1`

  rows, err := global.BrewmDB.Query(sqlStatement, session.Id)

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
    c.JSON(http.StatusBadRequest, gin.H{"error": "Session is already in progress. One session can be active at a time."})
    return
  }

  note := c.PostForm("note")
  timestamp := time.Now()

  sqlStatement := `
    INSERT INTO sessions (start_time, note)
    VALUES ($1, $2)`

  result, err := global.BrewmDB.Exec(sqlStatement, timestamp, note)
  if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
    return
  }

  sessionId, err := result.LastInsertId()
  if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
    return
  }

  startSessionProcess(int(sessionId))

  c.JSON(http.StatusOK, gin.H{
    "message": "Session started.",
    "start_time": timestamp,
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

      time.Sleep(measureInterval * time.Second)
      select {
      case <-sessionChannel:
        log.WithFields(log.Fields{
          "session_id": id,
        }).Info("Stopping active session!")
        return
      default: // adding default will make it not block
        log.WithFields(log.Fields{
          "session_id": id,
        }).Debug("Rolling to next measurement!")
      }
    }
  }(id)
}

func insertTemperature(id int) {
  sqlStatement := `
    INSERT INTO measurements (session_id, timestamp, temperature)
    VALUES ($1, $2, $3);`

  _, err := global.BrewmDB.Exec(sqlStatement, id, time.Now(), ds18b20.ReadTemperature())
  if err != nil {
    log.WithFields(log.Fields{
      "session_id": id,
    }).Error("Failed to save measurement!")

  }
}

func StopSession(c *gin.Context) {
  log.Info(c)

  id := c.PostForm("id")

  if id == "" {
    c.JSON(http.StatusBadRequest, gin.H{"error": "No id found. Provide session id to stop."})
    return
  }

  sqlStatement := `
    SELECT (CASE WHEN stop_time IS NULL THEN 1 ELSE 0 END) as is_active
    FROM sessions
    WHERE id = $1`
  row := global.BrewmDB.QueryRow(sqlStatement, id)

  var isActive bool
  err := row.Scan(&isActive)

  if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{
      "message": "Checking session with the given id was unsuccesfull!",
      "id": id,
      "error": err.Error(),
    })
    return
  }

  if isActive == false {
    c.JSON(http.StatusBadRequest, gin.H{"error": "Given session is not active. Can't stop."})
    return
  }

  sqlStatement = `
    UPDATE sessions
    SET stop_time = $1
    WHERE id = $2`

  timestamp := time.Now()

  _, err = global.BrewmDB.Exec(sqlStatement, timestamp, id)
  if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
  }

  if sessionChannel != nil {
    close(sessionChannel)
  }

  c.JSON(http.StatusOK, gin.H{
    "message": "Session stopped.",
    "stop_time": timestamp,
    "id": id,
  })
}

