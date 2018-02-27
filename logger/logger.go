package main

import (
        "time"
        "fmt"

        "gobot.io/x/gobot"
        "gobot.io/x/gobot/drivers/spi"
        "gobot.io/x/gobot/platforms/raspi"
)

func main() {
        raspi := raspi.NewAdaptor()
        adc := spi.NewMCP3008Driver(raspi)

        work := func() {
                adc.Start()

                gobot.Every(1*time.Second, func() {
                        result, err := adc.Read(0)

                        if err == nil {
                          fmt.Printf("Hello world: %v\n", result)
                        } else {
                          fmt.Printf("Error\n")
                        }
                })
        }

        robot := gobot.NewRobot("logger",
                []gobot.Connection{raspi},
                []gobot.Device{adc},
                work,
        )

        robot.Start()
}
