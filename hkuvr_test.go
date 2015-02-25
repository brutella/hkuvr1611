package hkuvr1611

import (
	"fmt"
	"github.com/brutella/gouvr/uvr"
	"github.com/brutella/gouvr/uvr/1611"
	"github.com/brutella/hc/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

var info model.Info = model.Info{
	Name:         "Input",
	SerialNumber: "001",
	Manufacturer: "TA",
	Model:        "Accessory",
}

func Packet() uvr1611.Packet {
	bytes := []uvr.Byte{
		uvr.Byte(byte(uvr.DeviceTypeUVR1611)),
		uvr.Byte(0x00),
		uvr.Byte(0x00),
		uvr.Byte(0x00), uvr.Byte(0x00), uvr.Byte(0x00), uvr.Byte(0x00), uvr.Byte(0x00),
		uvr.Byte(0xFB), uvr.Byte(0x20),
		uvr.Byte(0xDD), uvr.Byte(0x72),
		uvr.Byte(0x11), uvr.Byte(0x20),
		uvr.Byte(0x22), uvr.Byte(0x20),
		uvr.Byte(0x33), uvr.Byte(0x20),
		uvr.Byte(0x44), uvr.Byte(0x20),
		uvr.Byte(0x55), uvr.Byte(0x20),
		uvr.Byte(0x00), uvr.Byte(0x00),
		uvr.Byte(0x00), uvr.Byte(0x00),
		uvr.Byte(0x00), uvr.Byte(0x00),
		uvr.Byte(0x00), uvr.Byte(0x00),
		uvr.Byte(0x00), uvr.Byte(0x00),
		uvr.Byte(0x00), uvr.Byte(0x00),
		uvr.Byte(0x00), uvr.Byte(0x00),
		uvr.Byte(0x00), uvr.Byte(0x00),
		uvr.Byte(0x00), uvr.Byte(0x00),
		uvr.Byte(0x01), uvr.Byte(0x00),
		uvr.Byte(0x00),
		uvr.Byte(0x00),
		uvr.Byte(0x00),
		uvr.Byte(0x00),
		uvr.Byte(0x00),
		uvr.Byte(0x00), uvr.Byte(0x00), uvr.Byte(0x00), uvr.Byte(0x00), uvr.Byte(0x00), uvr.Byte(0x00), uvr.Byte(0x00), uvr.Byte(0x00),
		uvr.Byte(0x00), uvr.Byte(0x00), uvr.Byte(0x00), uvr.Byte(0x00), uvr.Byte(0x00), uvr.Byte(0x00), uvr.Byte(0x00), uvr.Byte(0x00),
	}
	bytes = append(bytes, uvr1611.ChecksumFromBytes(bytes))

	packet, _ := uvr1611.PacketFromBytes(bytes)

	return packet
}

func TestThermometer(t *testing.T) {
	p := Packet()
	s := NewSensorForInputValue(p.Input1, info)
	assert.NotNil(t, s)

	thermometer := s.Model.(model.Thermometer)
	assert.NotNil(t, thermometer)
	assert.Equal(t, fmt.Sprintf("%.1f", thermometer.Temperature()), "25.1")
}

func TestThermostat(t *testing.T) {
	p := Packet()
	s := NewSensorForInputValue(p.Input2, info)
	assert.NotNil(t, s)

	thermostat := s.Model.(model.Thermostat)
	assert.NotNil(t, thermostat)
	assert.Equal(t, fmt.Sprintf("%.1f", thermostat.Temperature()), "22.1")
	assert.Equal(t, fmt.Sprintf("%.1f", thermostat.TargetTemperature()), "22.1")
	assert.Equal(t, thermostat.Mode(), model.ModeHeating)
	assert.Equal(t, thermostat.TargetMode(), model.ModeHeating)
}

func TestOutlet(t *testing.T) {
	p := Packet()
	outlets := uvr1611.OutletsFromValue(p.Outgoing)
	s := NewSensorForOutlet(outlets[0], info)
	assert.NotNil(t, s)

	outlet := s.Model.(model.Outlet)
	assert.NotNil(t, outlet)
	assert.True(t, outlet.IsOn())
	assert.True(t, outlet.IsInUse())
}
