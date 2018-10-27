package temperature

import (
  "log"
  "time"
  "fmt"
  "io/ioutil"
  "strings"
  "strconv"

  // "periph.io/x/periph/host"
  // "periph.io/x/periph/conn/onewire/onewirereg"
  // "periph.io/x/periph/devices/ds18b20"
  )

const sensorID = "28-0000051e015b"


type Measurement struct {
  Serial      int       `json:"serial,omitempty"`
  Timestamp   time.Time `json:"timestamp"`
  Temperature float64   `json:"temperature"`
}

func Sense() Measurement {
  return Measurement{Timestamp: time.Now(), Temperature: readTemperature()}
}

func readTemperature() float64 {
  sensorPath := fmt.Sprintf("/sys/bus/w1/devices/%s/w1_slave", sensorID)
  data, err := ioutil.ReadFile(sensorPath)
  if err != nil {
    log.Fatal(err)
    return 0.0
  }

  raw := string(data)
  tString := raw[strings.LastIndex(raw, "t=")+2:len(raw)-1]

  t, err := strconv.ParseFloat(tString, 64)
  if err != nil {
    log.Fatal(err)
    return 0.0
  }

  return t / 1000.0
}

func readTestTemperature() float64 {
  return 21.3
}

// This is not working. :(
// func Temperature3(){
//   // Make sure periph is initialized.
//   if _, err := host.Init(); err != nil {
//     log.Fatal(err)
//   }

//   // Use onewirereg 1-wire bus registry to find the first available 1-wire bus.
//   bus, err := onewirereg.Open("")
//   if err != nil {
//     log.Fatal(err)
//   }
//   defer bus.Close()

//   fmt.Println(ds18b20.ConvertAll(bus, 10))
//   // 0x7a00000131825228
//   dev, err := ds18b20.New(bus, 0x7a0000051e015b, 10)
//   if err != nil {
//       log.Fatal(err)
//   }
//   defer dev.Halt()

//   fmt.Println(dev.LastTemp())
// }


