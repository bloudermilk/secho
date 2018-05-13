package main

import (
	"os"
	"fmt"

  "github.com/bloudermilk/secho/pkg/secho"
)

func main() {
  config := secho.LoadConfig(os.Args[1])
	source := make(chan secho.Reading)
	fanout := secho.Fanout{Input: source}

	fmt.Printf("\nLoaded config: %+v\n\n", config)

	go fanout.DoFan()
  // go secho.StartSensor(config, source)
  go secho.Advertise(config, fanout)
	// go secho.Log(fanout.Subscribe())

  select {}
}
