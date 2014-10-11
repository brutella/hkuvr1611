package main

import (
    "log"
    "fmt"
    "time"
    "flag"
    
    "github.com/brutella/gouvr/uvr/1611"
    
    "github.com/brutella/hap/app"
    "github.com/brutella/hap/server"
    "github.com/brutella/hap/model"
    "github.com/brutella/hap/model/accessory"
    "github.com/brutella/hap/common"
    
    "github.com/brutella/hcuvr1611/gpio"
    "github.com/brutella/hcuvr1611/mock"
)

func updateAccessories(packet uvr1611.Packet) {
    _, in1 := uvr1611.DecodeInputValue(packet.Input1)
    _, in2 := uvr1611.DecodeInputValue(packet.Input2)
    _, in3 := uvr1611.DecodeInputValue(packet.Input3)
    _, in4 := uvr1611.DecodeInputValue(packet.Input4)
    _, in5 := uvr1611.DecodeInputValue(packet.Input5)
    _, in6 := uvr1611.DecodeInputValue(packet.Input6)
    _, in7 := uvr1611.DecodeInputValue(packet.Input7)
    
    thermostat1 := accessorySensorName("Aussen")
    thermostat2 := accessorySensorName("Fussbodenheizung Vorlauf")
    thermostat3 := accessorySensorName("Buffer Oben")
    thermostat4 := accessorySensorName("Buffer Mitte")
    thermostat5 := accessorySensorName("Buffer Unten")
    thermostat6 := accessorySensorName("Raum")
    thermostat7 := accessorySensorName("Wärmetauscher Sekundär")
    
    thermostat1.SetTemperature(float64(in1))
    thermostat2.SetTemperature(float64(in2))
    thermostat3.SetTemperature(float64(in3))
    thermostat4.SetTemperature(float64(in4))
    thermostat5.SetTemperature(float64(in5))
    thermostat6.SetTemperature(float64(in6))
    thermostat7.SetTemperature(float64(in7))
}

func accessorySensorName(name string) model.Thermostat {
    thermostat, found := thermostats[name]
    if found == true {
        return thermostat
    }
    
    fmt.Println("Create new thermostat for", name)
    
    serial := common.GetSerialNumberForAccessoryName(name, application.Storage)
    info := model.Info{
        Name: name,
        Serial: serial, 
        Manufacturer: "TA",
        Model: "UVR1611",
    }
    
    thermo := accessory.NewThermostat(info, 10, 0, 100, 1)
    application.AddAccessory(thermo.Accessory)
    thermostats[name] = thermo
    
    return thermo
}

var application *app.App
var thermostats map[string]model.Thermostat


type Connection interface {
    Close()
}

func main() {
    var (
        mode = flag.String("conn", "sim", "Connection type; sim or gpio")
        port = flag.String("port", "P8_07", "GPIO port; default P8_07")
    )
    
    flag.Parse()
    
    thermostats = map[string]model.Thermostat{}
    
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
    
    callback := func(packet uvr1611.Packet) {
        updateAccessories(packet)

        fmt.Println(time.Now().Format(time.Stamp))
        fmt.Println("Zeit:", packet.Timestamp.ToString())
        fmt.Println("Aussentemperatur:", uvr1611.InputValueToString(packet.Input1))
        fmt.Println("Fussbodenheizung Vorlauf:", uvr1611.InputValueToString(packet.Input2))
        fmt.Println("Buffer Warmwasser")
        fmt.Println("   Oben:", uvr1611.InputValueToString(packet.Input3))
        fmt.Println("   Mitte:", uvr1611.InputValueToString(packet.Input4))
        fmt.Println("   Unten:", uvr1611.InputValueToString(packet.Input5))
        fmt.Println("Raumtemperatur:", uvr1611.InputValueToString(packet.Input6))
        fmt.Println("Wärmetauscher Sekundär:", uvr1611.InputValueToString(packet.Input7))
    }
    
    var conn Connection
    
    switch *mode {
    case "sim":
        conn = mock.NewConnection(callback)
    case "gpio":
        conn = gpio.NewConnection(*port, callback)
    default:
        log.Fatal("Incorrect -conn flag")
    }
    
    application.OnExit(func(){
        conn.Close()
    })
    
    application.Run()
}
