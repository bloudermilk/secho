package secho

import (
  "fmt"
)

func Log(source *chan Reading) {
  for {
    reading := <-*source
    fmt.Printf("Made a reading for '%s' - %f\n", reading.Sensor.Label, reading.Value)
  }
}
