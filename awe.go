package awe

import (
	"encoding/binary"
	"errors"
	"io"
)

const (
	Raw           byte = 0
	FastUnmarshal byte = 1 << iota
	Crc32
	Pub
	Sub
	Unsub
	Queue
	Compress

	MarkerSize uint32 = 2
	packetHeaderSize = 12
)

var (
	ErrInvalidPacket = errors.New("invalid Packet")
)

type Packet struct {
    Op    string
    Key   string
    Value []byte
}

func (p* Packet) Encode() (buf []byte) {
    opLen  := uint32(len(p.Op))
    keyLen := uint32(len(p.Key))
    valLen := uint32(len(p.Value))

    totalLen := packetHeaderSize + len(p.Op) + len(p.Key) + len(p.Value)
    buf = make([]byte, totalLen)

    binary.LittleEndian.PutUint32(buf[0:], opLen)
    binary.LittleEndian.PutUint32(buf[4:], keyLen)
    binary.LittleEndian.PutUint32(buf[8:], valLen)

    offset := packetHeaderSize
    copy(buf[offset:], p.Op)
    offset += int(opLen)
    copy(buf[offset:], p.Key)
    offset += int(keyLen)
    copy(buf[offset:], p.Value)
    return buf
}

func (p* Packet) Decode(buf []byte) (error) {
    if len(buf) < packetHeaderSize {
        return io.ErrUnexpectedEOF
    }

    opLen  := binary.LittleEndian.Uint32(buf[0:])
    keyLen := binary.LittleEndian.Uint32(buf[4:])
    valLen := binary.LittleEndian.Uint32(buf[8:])

    totalLen := int(packetHeaderSize + opLen + keyLen + valLen)
    if len(buf) < totalLen {
        return io.ErrUnexpectedEOF
    }

    opBytes  := buf[packetHeaderSize : packetHeaderSize+opLen]
    keyBytes := buf[packetHeaderSize+opLen : packetHeaderSize+opLen+keyLen]
    valBytes := buf[packetHeaderSize+opLen+keyLen : packetHeaderSize+opLen+keyLen+valLen]

	p.Key   = string(keyBytes)
	p.Op    = string(opBytes)
	p.Value = valBytes
	return nil
}
