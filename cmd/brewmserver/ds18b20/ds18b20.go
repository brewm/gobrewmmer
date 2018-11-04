package ds18b20

import (
  "fmt"
  "strings"
  "strconv"
  "io/ioutil"

  log "github.com/sirupsen/logrus"

  // "periph.io/x/periph/host"
  // "periph.io/x/periph/conn/onewire/onewirereg"
  // "periph.io/x/periph/devices/ds18b20"
)

const sensorId = "28-0000051e015b"

func ReadTemperature() float64 {
  sensorPath := fmt.Sprintf("/sys/bus/w1/devices/%s/w1_slave", sensorId)
  data, err := ioutil.ReadFile(sensorPath)
  if err != nil {
    log.WithFields(log.Fields{
      "err": err,
    }).Error("Reading temperature data failed!")
    return 0.0
  }

  raw := string(data)
  tString := raw[strings.LastIndex(raw, "t=")+2:len(raw)-1]

  t, err := strconv.ParseFloat(tString, 64)
  if err != nil {
    log.WithFields(log.Fields{
      "err": err,
    }).Error("Parsing temperature data failed!")

    return 0.0
  }

  return t / 1000.0
}


// This is not working. :(
// func ReadTemperature(){
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
