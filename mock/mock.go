package mock

import(
    "github.com/brutella/gouvr/uvr"
    "github.com/brutella/gouvr/uvr/1611"
    "fmt"
    "log"
    "time"
    "math/rand"
)

type mock struct {
    callback uvr1611.PacketCallback
}

func NewConnection(callback uvr1611.PacketCallback) *mock {    
    m := &mock{callback: callback}
    go m.SimulatePackets()
    
    return m
}

func (m *mock) Close() {
    fmt.Println("Close")
}
func B(b byte) uvr.Byte {
    return uvr.Byte(b)
}

func (m *mock) SimulatePackets() {
    ticker := time.NewTicker(7 * time.Second)
    
    m.sendRandomPacket()
    for _ = range ticker.C {
        m.sendRandomPacket()
    }
}

func (m *mock) sendRandomPacket() {
    delta := byte(rand.Intn(120))
    bytes := []uvr.Byte{
        B(byte(uvr.DeviceTypeUVR1611)),
        B(0x00),
        B(0x00),
        B(0x00), B(0x00), B(0x00), B(0x00), B(0x00),
        B(0xFA + delta), B(0x20),
        B(0xAF + delta), B(0x72),
        B(0x11 + delta), B(0x20),
        B(0x22 + delta), B(0x20),
        B(0x33 + delta), B(0x20),
        B(0x44 + delta), B(0x20),
        B(0x55 + delta), B(0x20),
        B(0x00), B(0x90),
        B(0x00), B(0x00),
        B(0x00), B(0x00),
        B(0x00), B(0x00),
        B(0x00), B(0x00),
        B(0x00), B(0x00),
        B(0x00), B(0x00),
        B(0x00), B(0x00),
        B(0x00), B(0x00),
        B(0xA5), B(0x00),
        B(0x00),
        B(0x00),
        B(0x00),
        B(0x00),
        B(0x00),
        B(0x00), B(0x00), B(0x00), B(0x00), B(0x00), B(0x00), B(0x00), B(0x00),
        B(0x00), B(0x00), B(0x00), B(0x00), B(0x00), B(0x00), B(0x00), B(0x00),
    }
    bytes = append(bytes, uvr1611.ChecksumFromBytes(bytes))
    
    packet, err := uvr1611.PacketFromBytes(bytes)
    if err != nil {
        log.Fatal(err)
    }
    
    m.callback(packet)
}
