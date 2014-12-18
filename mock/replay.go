package mock

import(
    "github.com/brutella/gouvr/uvr"
    "github.com/brutella/gouvr/uvr/1611"
    "fmt"
    "log"
    "time"
)

type replay struct {
    filePath string
    callback uvr1611.PacketCallback
}

func NewReplayConnection(filePath string, callback uvr1611.PacketCallback) *replay {    
    r := &replay{filePath, callback}
    go r.SimulatePackets()
    
    return r
}

func (r *replay) Close() {
    fmt.Println("Close")
}

func (r *replay) SimulatePackets() {
    ticker := time.NewTicker(10 * time.Second)
    
    r.sendPacket()
    for _ = range ticker.C {
        r.sendPacket()
    }
}

func (r *replay) sendPacket() {
    packetReceiver  := uvr1611.NewPacketReceiver()
    packetDecoder   := uvr1611.NewPacketDecoder(packetReceiver)
    byteDecoder     := uvr.NewByteDecoder(packetDecoder, uvr.NewTimeout(488.0, 0.2))
    syncDecoder     := uvr1611.NewSyncDecoder(byteDecoder, byteDecoder, uvr.NewTimeout(488.0*2, 0.3))
    interrupt       := uvr.NewEdgeInterrupt(syncDecoder)
    replayer        := uvr.NewReplayer(interrupt)
    
    packetReceiver.RegisterCallback(func(packet uvr1611.Packet) {
        r.callback(packet)
        syncDecoder.Reset()
        byteDecoder.Reset()
        packetDecoder.Reset()
    })
    err := replayer.Replay(r.filePath)
    if err != nil {
        log.Fatal("Could not replay file.", err)
    }
}
