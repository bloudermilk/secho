package secho

import (
	"time"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type SechoConfig struct {
  Name      string
	Frequency float64 // In Hz
	Sensors   []Sensor
}

type Sensor struct {
	Channel 		 int
	Label    		 string
	Unit    		 string
	LowerLimit	 float64 // In Unit
	UpperLimit   float64 // In Unit
  UUID         string
}

func LoadConfig(path string) SechoConfig {
	data, err := ioutil.ReadFile(path)
	CheckError(err)

	config := SechoConfig{}
	err = yaml.Unmarshal([]byte(data), &config)
	CheckError(err)

	if config.Frequency == 0 {
		config.Frequency = 1
	}

	return config
}

func (config SechoConfig) PollingInterval() time.Duration {
	return time.Duration(float64(time.Second) / config.Frequency)
}
