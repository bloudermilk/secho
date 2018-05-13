package secho

import (
  "fmt"
  "log"
  "context"
  "time"
  "math"
  "encoding/binary"

  "github.com/pkg/errors"
  "github.com/go-ble/ble"
  "github.com/go-ble/ble/linux"
)

var (
  SensorServiceUUID = ble.UUID("00000000-8dc0-41e7-b525-c226a9b1f5ad")
  CharacteristicUserDescriptionUUID = ble.UUID16(0x2901)
  CharacteristicPresentationFormatUUID = ble.UUID16(0x2904)
)

func Advertise(config SechoConfig, fanout Fanout) {
  device, err := linux.NewDeviceWithName(config.Name)

  CheckError(err)

  ble.SetDefaultDevice(device)

  sensorsService := ble.NewService(SensorServiceUUID)

  for _, sensor := range(config.Sensors) {
    characteristic := NewSensorCharacteristic(sensor, func() chan Reading {
      return fanout.Subscribe()
    })

    sensorsService.AddCharacteristic(characteristic)
  }

  err = ble.AddService(sensorsService)
  CheckError(err)

  fmt.Println("Advertising...")

  ctx := ble.WithSigHandler(context.WithTimeout(context.Background(), time.Hour))

  chkErr(ble.AdvertiseNameAndServices(ctx, config.Name))
}

func NewSensorCharacteristic(sensor Sensor, subscribe func() chan Reading) *ble.Characteristic {
  characteristic := ble.NewCharacteristic(ble.UUID(sensor.UUID))

  characteristic.NewDescriptor(CharacteristicUserDescriptionUUID).SetValue([]byte(sensor.Label))

  characteristic.HandleRead(ble.ReadHandlerFunc(func(req ble.Request, rsp ble.ResponseWriter) {
		fmt.Fprintf(rsp, "count: Read %d", 0)
  }))

  characteristic.HandleNotify(ble.NotifyHandlerFunc(func(_ ble.Request, notifier ble.Notifier) {
    channel := subscribe()

    for {
      reading := <-time.After(time.Second)
      fmt.Printf("got that reading %+v\n", reading)

      if (reading.Sensor == sensor) {
        update := float64bytes(reading.Value)
        notifier.Write(update)
      }
    }
  }))

  // TODO: Do format string
  // characteristic.NewDescriptor(ble.UUID16(CharacteristicPresentationFormatUUID)).SetValue([]byte(sensor.Label))

  return characteristic
}

func float64bytes(float float64) []byte {
    bits := math.Float64bits(float)
    bytes := make([]byte, 8)
    binary.LittleEndian.PutUint64(bytes, bits)
    return bytes
}

func chkErr(err error) {
	switch errors.Cause(err) {
	case nil:
	case context.DeadlineExceeded:
		fmt.Println("done")
	case context.Canceled:
		fmt.Println("canceled")
	default:
		log.Fatalf(err.Error())
	}
}
