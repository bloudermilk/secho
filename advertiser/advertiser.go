package main

import (
  "fmt"
  "log"
  "context"
  "time"

  "github.com/go-ble/ble"
  "github.com/pkg/errors"

  "./dev"
)

var (
  ManufactererName = "Brendan Loudermilk"
  Name = "Vanberry"
	BatteryServiceUUID = ble.UUID16(0x180F)
  BatteryLevelCharacteristicUUID = ble.UUID16(0x2A19)
  DeviceInformationServiceUUID = ble.UUID16(0x180A)
  ManufacturerNameStringCharacteristicUUID = ble.UUID16(0x2A29)
)

func main() {
  device, err := dev.NewDevice(Name)

  if err != nil {
		log.Fatalf("Problem initializing device: %s", err)
	}

  ble.SetDefaultDevice(device)

  service := ble.NewService(BatteryServiceUUID)
  service.AddCharacteristic(BatteryCharacteristic())

  if err := ble.AddService(service); err != nil {
		log.Fatalf("can't add service: %s", err)
	}

  deviceService := NewDeviceInformationService()
  if err := ble.AddService(deviceService); err != nil {
		log.Fatalf("can't add service: %s", err)
	}

  fmt.Printf("Advertising...\n")

  ctx := ble.WithSigHandler(context.WithTimeout(context.Background(), time.Hour))
  chkErr(ble.AdvertiseNameAndServices(ctx, Name, BatteryServiceUUID, deviceService.UUID))
}

func chkErr(err error) {
	switch errors.Cause(err) {
	case nil:
	case context.DeadlineExceeded:
		fmt.Printf("done\n")
	case context.Canceled:
		fmt.Printf("canceled\n")
	default:
		log.Fatalf(err.Error())
	}
}

func BatteryCharacteristic() *ble.Characteristic {
  characteristic := ble.NewCharacteristic(BatteryLevelCharacteristicUUID)
  load := byte(100)

  characteristic.HandleRead(ble.ReadHandlerFunc(func(request ble.Request, response ble.ResponseWriter) {
    response.Write([]byte{load})
  }))

  // Characteristic User Description
  characteristic.NewDescriptor(ble.UUID16(0x2901)).SetValue([]byte("Batttery Level?"))

  // Characteristic Presentation Format
  characteristic.NewDescriptor(ble.UUID16(0x2904)).SetValue([]byte{4, 1, 39, 173, 1, 0, 0})

  return characteristic
}

func NewDeviceInformationService() *ble.Service {
  service := ble.NewService(DeviceInformationServiceUUID)

  service.AddCharacteristic(NewManufacturerNameStringCharacteristic())

  return service
}

func NewManufacturerNameStringCharacteristic() *ble.Characteristic {
  characteristic := ble.NewCharacteristic(ManufacturerNameStringCharacteristicUUID)

  characteristic.SetValue([]byte(ManufactererName))

  return characteristic
}
