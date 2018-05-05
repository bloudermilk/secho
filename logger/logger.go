package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/spi"
	"gobot.io/x/gobot/platforms/raspi"

	"gopkg.in/yaml.v2"
)

const ADCBits = 10

type SensorConfig struct {
	Channel int
	Name    string
	Unit    string
	Limit   float64 // In unit being measured
}

type LoggerConfig struct {
	Frequency float64 // In Hz
	Sensors   []SensorConfig
}

type SensorReading struct {
	Sensor         SensorConfig
	DigitalReading int
	Time           time.Time
}

func (loggerConfig LoggerConfig) PollingInterval() time.Duration {
	return time.Duration(float64(time.Second) / loggerConfig.Frequency)
}

func (sensorReading SensorReading) ScaledReading() float64 {
	return float64(sensorReading.DigitalReading) * (sensorReading.Sensor.Limit / float64(math.Pow(2, ADCBits)))
}

func ReadingsForSensors(sensors []SensorConfig, f func(SensorConfig) SensorReading) []SensorReading {
	readings := make([]SensorReading, len(sensors))
	for i, sensor := range sensors {
		readings[i] = f(sensor)
	}
	return readings
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func loadConfig() LoggerConfig {
	path := os.Args[1]

	fmt.Println("Reading config file from path: ", path)

	data, err := ioutil.ReadFile(path)
	check(err)

	fmt.Println("Parsing config file...")

	config := LoggerConfig{}
	err = yaml.Unmarshal([]byte(data), &config)
	check(err)

	if config.Frequency == 0 {
		config.Frequency = 1
	}

	fmt.Printf("Parsed config file successfully: %+v\n", config)

	return config
}

func main() {
	config := loadConfig()
	raspi := raspi.NewAdaptor()
	adc := spi.NewMCP3008Driver(raspi)

	work := func() {
		fmt.Printf("Polling at %fHz (every %s)\n", config.Frequency, config.PollingInterval())

		gobot.Every(config.PollingInterval(), func() {
			readings := ReadingsForSensors(config.Sensors, func(sensor SensorConfig) SensorReading {
				reading, err := adc.Read(sensor.Channel)

				check(err)

				sensorReading := SensorReading{Sensor: sensor, DigitalReading: reading, Time: time.Now()}

				fmt.Printf("%s at %f%s\n", sensorReading.Sensor.Name, sensorReading.ScaledReading(), sensorReading.Sensor.Unit)

				return sensorReading
			})

			fmt.Printf("\nReadings: %+v\n\n", readings)
		})
	}

	robot := gobot.NewRobot("logger",
		[]gobot.Connection{raspi},
		[]gobot.Device{adc},
		work,
	)

	robot.Start()
}
