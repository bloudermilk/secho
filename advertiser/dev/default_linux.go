package dev

import (
	"github.com/go-ble/ble"
	"github.com/go-ble/ble/linux"
)

// DefaultDevice ...
func DefaultDevice(name string) (d ble.Device, err error) {
	return linux.NewDeviceWithName(name)
}
