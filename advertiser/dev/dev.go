package dev

import "github.com/go-ble/ble"

// NewDevice ...
func NewDevice(name string) (dev ble.Device, err error) {
	return DefaultDevice(name)
}
