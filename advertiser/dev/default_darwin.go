package dev

import (
	"github.com/go-ble/ble"
	"github.com/go-ble/ble/darwin"
)

// DefaultDevice ...
func DefaultDevice(name string) (d ble.Device, err error) {
	return darwin.NewDeviceWithName(name)
}
