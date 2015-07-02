# hkuvr1611

The hkuvr1611 is a daemon to simulate a HomeKit bridge for an [UVR1611][uvr1611] device. It uses the [gouvr][gouvr] library to read data from the device's data bus and the [HomeControl][hc] library to create a HomeKit bridge.

[uvr1611]: http://www.ta.co.at/en/products/uvr1611/
[hc]: https://github.com/brutella/hc
[gouvr]: https://github.com/brutella/gouvr

## Build

You can use the included Makefile to build the `daemon/hkuvr1611d.go` for a RaspberryPi and Beaglebone Black.

```make
// RaspberryPi
make rpi
    
// Beaglebone Black
make bbb
```

## Usage

You need to provide the following arguments when running the `hkuvr1611d` daemon.

- conn: Specifies the connection type
    - mock: Simulates the data bus with random values; *default*
    - gpio: Uses a gpio to connect to the data bus
    - replay: Replays a log file (see *logs/testlog.log*)
- file: Log file from which to replay packets
- port: GPIO port; *default: P8_07*
- timeout: Timeout in seconds until accessories are not reachable; *default: 120*

Example

    hkuvr1611d -conn=gpio -port=P8_07 -timeout=360

# Contact

Matthias Hochgatterer

Github: [https://github.com/brutella](https://github.com/brutella/)

Twitter: [https://twitter.com/brutella](https://twitter.com/brutella)


# License

hkuvr1611 is available under a non-commercial license. See the LICENSE file for more info.