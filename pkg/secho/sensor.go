package secho

import (
	"fmt"
	"math"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/spi"
	"gobot.io/x/gobot/platforms/raspi"
)

const ADCBits = 10

func ScaleReading(sensor Sensor, reading int) float64 {
	return float64(reading) * (sensor.UpperLimit / float64(math.Pow(2, ADCBits)))
}

func StartSensor(config SechoConfig) {
	adaptor := raspi.NewAdaptor()
	source := spi.NewMCP3008Driver(adaptor)

	work := func() {
		fmt.Printf("Polling at %fHz (every %s)\n", config.Frequency, config.PollingInterval())

		gobot.Every(config.PollingInterval(), func() {
			for _, sensor := range(config.Sensors) {
				digitalReading, err := source.Read(sensor.Channel)

				CheckError(err)

				scaledReading := ScaleReading(sensor, digitalReading)

				fmt.Printf("Made a reading for '%s' - %f", sensor.Label, scaledReading)

				sensor.Readings <- scaledReading
			}
		})
	}

	robot := gobot.NewRobot(config.Name,
		[]gobot.Connection{adaptor},
		[]gobot.Device{source},
		work,
	)

	robot.Start()
}
