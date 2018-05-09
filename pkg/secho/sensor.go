package secho

import (
	"fmt"
	"math"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/spi"
	"gobot.io/x/gobot/platforms/raspi"
)

const ADCBits = 10

type Reading struct {
	Sensor Sensor
	Value float64
	Timestamp time.Time
}

func ScaleReading(sensor Sensor, reading int) float64 {
	return float64(reading) * (sensor.UpperLimit / float64(math.Pow(2, ADCBits)))
}

func StartSensor(config SechoConfig, dest chan Reading) {
	adaptor := raspi.NewAdaptor()
	source := spi.NewMCP3008Driver(adaptor)

	work := func() {
		fmt.Printf("Polling at %fHz (every %s)\n", config.Frequency, config.PollingInterval())

		gobot.Every(config.PollingInterval(), func() {
			for _, sensor := range(config.Sensors) {
				rawReading, err := source.Read(sensor.Channel)
				scaledReading := ScaleReading(sensor, rawReading)

				CheckError(err)

				// fmt.Printf("Reading for %s: %f\n", sensor.Label, scaledReading)

				dest <- Reading{ Sensor: sensor, Value: scaledReading, Timestamp: time.Now() }
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
