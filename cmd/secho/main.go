package main

import (
	"os"

  "github.com/bloudermilk/secho/pkg/secho"
)

func main() {
  config := secho.LoadConfig(os.Args[1])

  // go secho.StartSensor(config)
  go secho.Advertise(config)

  select {}
}
