package hkuvr1611

import (
	"errors"

	"github.com/brutella/gouvr/uvr"
	"github.com/brutella/gouvr/uvr/1611"
	"github.com/brutella/hc/model"
	"github.com/brutella/hc/model/accessory"
)

type Sensor struct {
	Model     model.Accessory
	Accessory *accessory.Accessory
}

func InputValuesFromPacket(p uvr1611.Packet) []uvr.Value {
	return []uvr.Value{
		p.Input1, p.Input2, p.Input3, p.Input4, p.Input5,
		p.Input6, p.Input7, p.Input8, p.Input9, p.Input10,
		p.Input11, p.Input12, p.Input13, p.Input14, p.Input15,
		p.Input16,
	}
}

func NewSensorForOutlet(o uvr1611.Outlet, info model.Info) *Sensor {
	a := accessory.NewOutlet(info)
	UpdateAccessoryWithOutlet(a, o)
	return &Sensor{a, a.Accessory}
}

func NewSensorForInputValue(v uvr.Value, info model.Info) *Sensor {
	input_type, _ := uvr1611.DecodeInputValue(v)

	var s *Sensor
	switch input_type {
	case uvr1611.InputTypeUnused:
		break
	case uvr1611.InputTypeDigital:
		o := accessory.NewOutlet(info)
		s = &Sensor{o, o.Accessory}
	case uvr1611.InputTypeRoomTemperature:
		// -50 ... +199°C
		t := accessory.NewThermostat(info, 0, -50, 199, 1)
		s = &Sensor{t, t.Accessory}
	case uvr1611.InputTypeTemperature:
		t := accessory.NewThermometer(info, 0, -50, 199, 1)
		s = &Sensor{t, t.Accessory}
	case uvr1611.InputTypeVolumeFlow:
		// TODO(brutella) ?
		break
	case uvr1611.InputTypeRadiation:
		// TODO(brutella) ?
		break
	}

	if s != nil {
		UpdateAccessoryWithInputValue(s.Model, v)
	}

	return s
}

func NewSensorForHeatMeter(hm uvr.HeatMeterValue, info model.Info) *Sensor {
	// TODO
	return nil
}

func UpdateAccessoryWithOutlet(a model.Accessory, o uvr1611.Outlet) error {
	if outlet, ok := a.(model.Outlet); ok == true {
		outlet.SetOn(o.Enabled)
		outlet.SetInUse(true)
	} else {
		return errors.New("Outlet expects outlet accessory")
	}

	return nil
}

func UpdateAccessoryWithInputValue(a model.Accessory, v uvr.Value) error {
	input_type, value := uvr1611.DecodeInputValue(v)

	switch input_type {
	case uvr1611.InputTypeUnused:
		break
	case uvr1611.InputTypeDigital:
		if outlet, ok := a.(model.Outlet); ok == true {
			on := false
			if value == 1.0 {
				on = true
			}
			outlet.SetOn(on)
			outlet.SetInUse(true)
		} else {
			return errors.New("Digital input expects outlet accessory")
		}
	case uvr1611.InputTypeRoomTemperature:
		if thermostat, ok := a.(model.Thermostat); ok == true {
			mode := model.HeatCoolModeOff
			switch uvr1611.RoomTemperatureModeFromValue(v) {
			case uvr1611.RoomTemperatureModeAutomatic:
				mode = model.HeatCoolModeAuto
			case uvr1611.RoomTemperatureModeNormal:
				mode = model.HeatCoolModeHeat
			case uvr1611.RoomTemperatureModeLowering:
				mode = model.HeatCoolModeCool
			case uvr1611.RoomTemperatureModeStandby:
				mode = model.HeatCoolModeOff
			}
			if mode == model.HeatCoolModeAuto {
				thermostat.SetMode(model.HeatCoolModeOff)
			} else {
				thermostat.SetMode(mode)
			}

			thermostat.SetTemperature(float64(value))
			// Target == Current
			thermostat.SetTargetTemperature(float64(value))
			thermostat.SetTargetMode(mode)
		} else {
			return errors.New("Room temperature input expects thermostat accessory")
		}
	case uvr1611.InputTypeTemperature:
		if thermometer, ok := a.(model.Thermometer); ok == true {
			thermometer.SetTemperature(float64(value))
		} else {
			return errors.New("Temperature input expects thermometer accessory")
		}
	}

	return nil
}

func UpdateAccessoryWithHeatMeter(a model.Accessory, hm uvr.HeatMeterValue) error {
	// TODO
	return nil
}
