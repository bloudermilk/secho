package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"
  "math"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/spi"
	"gobot.io/x/gobot/platforms/raspi"

	"gopkg.in/yaml.v2"
)

const ADCBits = 10

type ChannelConfig struct {
	Channel int
	Name    string
	Max     float64
}

type LoggerConfig struct {
	VRef      float64
	Bits      int
	Frequency float64 // In Hertz
	Channels  []ChannelConfig
}

type ChannelReading struct {
	Channel ChannelConfig
	RawReading int
	Time    time.Time
}

func (loggerConfig LoggerConfig) PollingInterval() time.Duration {
	return time.Duration(float64(time.Second) / loggerConfig.Frequency)
}

func (channelReading ChannelReading) ScaledReading() float64 {
  return float64(channelReading.RawReading) * (channelReading.Channel.Max / float64(math.Pow(2, ADCBits)))
}

func MapChannelToReading(vs []ChannelConfig, f func(ChannelConfig) ChannelReading) []ChannelReading {
	vsm := make([]ChannelReading, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
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
			readings := MapChannelToReading(config.Channels, func(channel ChannelConfig) ChannelReading {
				reading, err := adc.Read(channel.Channel)

				check(err)

        channelReading := ChannelReading{Channel: channel, RawReading: reading, Time: time.Now()}

        fmt.Printf("%s at %f\n", channelReading.Channel.Name, channelReading.ScaledReading())

				return channelReading
			})

			fmt.Println("\nReadings: ", readings)
		})
	}

	robot := gobot.NewRobot("logger",
		[]gobot.Connection{raspi},
		[]gobot.Device{adc},
		work,
	)

	robot.Start()
}
