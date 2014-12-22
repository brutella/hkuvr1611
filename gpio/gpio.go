package gpio

import(
    "github.com/brutella/gouvr/uvr"
    "github.com/brutella/gouvr/uvr/1611"
    
    "github.com/kidoman/embd"
    _"github.com/kidoman/embd/host/all"
    
    "log"
    "fmt"
    "math/big"
)

func InitGPIO(file string) (embd.DigitalPin, error) {
    embd.InitGPIO()
    pin, pin_err := embd.NewDigitalPin(file)
    if pin_err != nil {
        log.Fatal("Error opening pin! \n", pin_err)
        return nil, pin_err
    }
    
    pin.SetDirection(embd.In)
    
    return pin, nil
}

type gpio struct {
    pin embd.DigitalPin
}

func NewConnection(file string, callback uvr1611.PacketCallback) *gpio {
    pin, err := InitGPIO(file)
    if err != nil {
        log.Fatal(err)
    }
    
    packetReceiver  := uvr1611.NewPacketReceiver()
    packetDecoder   := uvr1611.NewPacketDecoder(packetReceiver)
    byteDecoder     := uvr.NewByteDecoder(packetDecoder, uvr.NewTimeout(488.0, 0.4))
    syncDecoder     := uvr1611.NewSyncDecoder(byteDecoder, byteDecoder, uvr.NewTimeout(488.0*2, 0.4))
    signal          := uvr.NewSignal(syncDecoder)
    
    packetReceiver.RegisterCallback(func(packet uvr1611.Packet) {
        if callback != nil {
            callback(packet)
        }
        
        syncDecoder.Reset()
        byteDecoder.Reset()
        packetDecoder.Reset()
    })
    
    err = pin.Watch(embd.EdgeBoth, func(pin embd.DigitalPin) {
        value, read_err := pin.Read()
        if read_err != nil {
            fmt.Println(read_err)
        } else {
            signal.Consume(big.Word(value))
        }
    })
     
    if err != nil {
        log.Fatal("Could not watch pin.", err)
    }
    
    return &gpio{
        pin: pin,
    }
}

func (g *gpio) Close() {
    g.pin.Close()
    embd.CloseGPIO()
}