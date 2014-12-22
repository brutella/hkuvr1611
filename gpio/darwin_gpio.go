// +build darwin

package gpio

import(
    "github.com/brutella/gouvr/uvr/1611"
    "log"
)

type gpio struct {
}

func NewConnection(file string, callback uvr1611.PacketCallback) *gpio {
    log.Fatal("GPIO not supported")
    return &gpio{}
}

func (g *gpio) Close() {
}