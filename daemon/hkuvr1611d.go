package main

import (
    "fmt"
    "time"
    "flag"
    "reflect"
    
    "github.com/brutella/hkuvr1611"
    "github.com/brutella/gouvr/uvr"
    "github.com/brutella/gouvr/uvr/1611"
    "github.com/brutella/log"
    "github.com/brutella/hap/app"
    "github.com/brutella/hap/server"
    "github.com/brutella/hap/model"
    "github.com/brutella/hap/common"
    
    "github.com/brutella/hkuvr1611/gpio"
    "github.com/brutella/hkuvr1611/mock"
)

func HandlePacket(p uvr1611.Packet) {
    inputs := hkuvr1611.InputValuesFromPacket(p)
    for i, v := range inputs {
        HandleInputValueWithName(v, fmt.Sprintf("Sensor %d", i+1))
    }
    
    outlets := uvr1611.OutletsFromValue(p.Outgoing)
    for i, v := range outlets {
        HandleOutletWithName(v, fmt.Sprintf("Outlet %d", i+1))
    }
    
    h1, h2 := uvr1611.AreHeatMetersEnabled(p.HeatRegister)
    if h1 == true {
        HandleHeatMeterWithName(p.HeatMeter1, "Heat Meter 1")
    }
    
    if h2 == true {
        HandleHeatMeterWithName(p.HeatMeter2, "Heat Meter 2")
    }
}

func HandleInputValueWithName(v uvr.Value, name string) {
    var s *hkuvr1611.Sensor
    var found bool
    if s, found = sensors[name]; found == false {
        s := hkuvr1611.NewSensorForInputValue(v, InfoForAccessoryName(name))
        if s != nil {
            log.Println("[INFO]", reflect.TypeOf(s.Model), "with name", name)
            application.AddAccessory(s.Accessory)
            sensors[name] = s
        }
    } else {
        err := hkuvr1611.UpdateAccessoryWithInputValue(s.Model, v)
        if err != nil {
            log.Println("[ERRO]", err)
        }
    }
}

func HandleOutletWithName(o uvr1611.Outlet, name string) {
    var s *hkuvr1611.Sensor
    var found bool
    if s, found = sensors[name]; found == false {
        s := hkuvr1611.NewSensorForOutlet(o, InfoForAccessoryName(name))
        if s != nil {
            log.Println("[INFO]", reflect.TypeOf(s.Model), "with name", name)
            application.AddAccessory(s.Accessory)
            sensors[name] = s
        }
    } else {
        err := hkuvr1611.UpdateAccessoryWithOutlet(s.Model, o)
        if err != nil {
            log.Println("[ERRO]", err)
        }
    }
}

func HandleHeatMeterWithName(hm uvr.HeatMeterValue, name string) {
    // TODO
}

func InfoForAccessoryName(name string) model.Info {
    serial := common.GetSerialNumberForAccessoryName(name, application.Storage)
    info := model.Info{
        Name: name,
        SerialNumber: serial, 
        Manufacturer: "TA",
        Model: "UVR1611",
    }
    
    return info
}

// Access to HAP app
var application *app.App

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
        mode = flag.String("conn", "mock", "Connection type; mock, gpio, replay")
        file = flag.String("file", "", "Log file from which to replay packets")
        port = flag.String("port", "P8_07", "GPIO port; default P8_07")
        timeout = flag.Int("timeout", 120, "Timeout in seconds until accessories are not reachable")
    )
    
    flag.Parse()
    
    conf := app.NewConfig()
    conf.DatabaseDir = "./data"
    conf.BridgeName = "UVR1611Bridge"
    
    pwd, _ := server.NewPassword("11122333")
    conf.BridgePassword = pwd
    conf.BridgeManufacturer = "Matthias H."
    
    var err error
    application, err = app.NewApp(conf)
    if err != nil {
        log.Fatal(err)
    }
    
    sensors = map[string]*hkuvr1611.Sensor{}
    
    timer_duration := time.Duration(*timeout)* time.Second
    timer = time.AfterFunc(timer_duration, func() {
        log.Println("[INFO] Not Reachable")
        application.SetReachable(false)
    })
    
    callback := func(packet uvr1611.Packet) {
        application.PerformBatchUpdates(func(){
            HandlePacket(packet)
            application.SetReachable(true)
        })
        // fmt.Println(time.Now().Format(time.Stamp))
        // packet.Log()
        timer.Reset(timer_duration)
    }
    
    var conn Connection
    
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
    
    application.OnExit(func(){
        timer.Stop()
        conn.Close()
    })
    
    application.RunAndPublish(false)
}
