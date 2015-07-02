# hkuvr1611

The hkuvr1611 is a daemon to simulate a HomeKit bridge for an [UVR1611][uvr1611] device. It uses the [gouvr][gouvr] library to read data from the device's data bus and the [HomeControl][hc] library to create a HomeKit bridge.

[uvr1611]: http://www.ta.co.at/en/products/uvr1611/
[hc]: https://github.com/brutella/hc
[gouvr]: https://github.com/brutella/gouvr


## Build

Build `hkuvr1611d.go` using `go build daemon/hkuvr1611d.go` or use the Makefile to build for Beaglebone Black
    
    make bbb
    
or Raspberry Pi

    make rpi

## Run

You need to provide the following arguments when running the `hkuvr1611d` daemon.

- pin: Accessory pin required for pairing
- conn: Specifies the connection type
    - mock: Simulates the data bus with random values; *default*
    - gpio: Uses a gpio to connect to the data bus
    - replay: Replays a log file (see *logs/testlog.log*)
- file: Log file from which to replay packets
- port: GPIO port; *default: P8_07*
- timeout: Timeout in seconds until accessories are not reachable; *default: 120*

#### Examples
    
    // Simulate
    hkuvr1611d -pin=32112321
    
    // Connect via GPIO
    hkuvr1611d -pin=32112321 -conn=gpio -port=P8_07

## HomeKit Client

You need an iOS app to control HomeKit accessories. 
You can use [Home][home] which runs on iPhone, iPad and Apple Watch.

Read the [Getting Started][home-getting-started] guide.

[home]: http://selfcoded.com/home/
[home-getting-started]: http://selfcoded.com/home/getting-started/

# Contact

Matthias Hochgatterer

Github: [https://github.com/brutella](https://github.com/brutella/)

Twitter: [https://twitter.com/brutella](https://twitter.com/brutella)


# License

hkuvr1611 is available under a non-commercial license. See the LICENSE file for more info.