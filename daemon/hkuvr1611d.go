package main

import (
	"flag"
	"fmt"
	"reflect"
	"time"

	"github.com/brutella/gouvr/uvr"
	"github.com/brutella/gouvr/uvr/1611"
	"github.com/brutella/hc/hap"
	"github.com/brutella/hc/model"
	"github.com/brutella/hc/model/accessory"
	"github.com/brutella/hkuvr1611"
	"github.com/brutella/log"

	"github.com/brutella/hkuvr1611/gpio"
	"github.com/brutella/hkuvr1611/mock"
)

func HandlePacket(p uvr1611.Packet) []*accessory.Accessory {
	var sensors []*accessory.Accessory
	inputs := hkuvr1611.InputValuesFromPacket(p)
	for i, v := range inputs {
		s := HandleInputValueWithName(v, fmt.Sprintf("Sensor %d", i+1))
		sensors = append(sensors, s.Accessory)
	}

	outlets := uvr1611.OutletsFromValue(p.Outgoing)
	for i, v := range outlets {
		s := HandleOutletWithName(v, fmt.Sprintf("Outlet %d", i+1))
		sensors = append(sensors, s.Accessory)
	}

	h1, h2 := uvr1611.AreHeatMetersEnabled(p.HeatRegister)
	if h1 == true {
		s := HandleHeatMeterWithName(p.HeatMeter1, "Heat Meter 1")
		sensors = append(sensors, s.Accessory)
	}

	if h2 == true {
		s := HandleHeatMeterWithName(p.HeatMeter2, "Heat Meter 2")
		sensors = append(sensors, s.Accessory)
	}

	return sensors
}

func HandleInputValueWithName(v uvr.Value, name string) (s *hkuvr1611.Sensor) {
	var found bool
	if s, found = sensors[name]; found == false {
		s = hkuvr1611.NewSensorForInputValue(v, InfoForAccessoryName(name))
		if s != nil {
			log.Println("[INFO]", reflect.TypeOf(s.Model), "with name", name)
			sensors[name] = s
		}
	} else {
		err := hkuvr1611.UpdateAccessoryWithInputValue(s.Model, v)
		if err != nil {
			log.Println("[ERRO]", err)
		}
	}
	return s
}

func HandleOutletWithName(o uvr1611.Outlet, name string) (s *hkuvr1611.Sensor) {
	var found bool
	if s, found = sensors[name]; found == false {
		s = hkuvr1611.NewSensorForOutlet(o, InfoForAccessoryName(name))
		if s != nil {
			log.Println("[INFO]", reflect.TypeOf(s.Model), "with name", name)
			sensors[name] = s
		}
	} else {
		err := hkuvr1611.UpdateAccessoryWithOutlet(s.Model, o)
		if err != nil {
			log.Println("[ERRO]", err)
		}
	}

	return s
}

func HandleHeatMeterWithName(hm uvr.HeatMeterValue, name string) (s *hkuvr1611.Sensor) {
	return nil
}

func InfoForAccessoryName(name string) model.Info {
	info := model.Info{
		Name:         name,
		Manufacturer: "TA",
		Model:        "UVR1611",
	}

	return info
}

// Access to HAP app
var transport hap.Transport
var uvrAccessory *accessory.Accessory

// List of sensors
var sensors map[string]*hkuvr1611.Sensor

// Timer to remove not reachable accessories
var timer *time.Timer

type Connection interface {
	Close()
}

// This app can connect to the UVR1611 data bus and provide the sensor values to HomeKit clients
//
// Optimizations: To improve the performance on a Raspberry Pi B+, the interrupt handler of the
// gpio pin is removed every time after successfully decoding a packet. This allows other goroutines
// (e.g. HAP server) to do their job more quickly.
func main() {
	var (
		mode    = flag.String("conn", "mock", "Connection type; mock, gpio, replay")
		file    = flag.String("file", "", "Log file from which to replay packets")
		port    = flag.String("port", "P8_07", "GPIO port; default P8_07")
		timeout = flag.Int("timeout", 120, "Timeout in seconds until accessories are not reachable")
	)

	flag.Parse()
	sensors = map[string]*hkuvr1611.Sensor{}

	info := InfoForAccessoryName("UVR1611")
	uvrAccessory = accessory.New(info)

	timer_duration := time.Duration(*timeout) * time.Second
	timer = time.AfterFunc(timer_duration, func() {
		log.Println("[INFO] Not Reachable")
		if transport != nil {
			sensors = map[string]*hkuvr1611.Sensor{}
			transport.Stop()
			transport = nil
		}
	})

	var conn Connection
	callback := func(packet uvr1611.Packet) {
		sensors := HandlePacket(packet)
		if transport == nil {
			var err error
			transport, err = hap.NewIPTransport("00102003", uvrAccessory, sensors...)

			transport.OnStop(func() {
				timer.Stop()
				conn.Close()
			})

			go func() {
				transport.Start()
			}()

			if err != nil {
				log.Fatal(err)
			}
		}
		timer.Reset(timer_duration)
	}

	switch *mode {
	case "mock":
		conn = mock.NewConnection(callback)
	case "replay":
		conn = mock.NewReplayConnection(*file, callback)
	case "gpio":
		conn = gpio.NewConnection(*port, callback)
	default:
		log.Fatal("Incorrect -conn flag")
	}

	select {}
}
