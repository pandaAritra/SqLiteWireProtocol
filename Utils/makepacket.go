package utils

import (
	"encoding/binary"
)

// MakePKT creates a wire protocol packet with the given message type and payload string.
// Message type 0x02 is typically used for SQL queries.
func MakePKT(msgType int, payload string) *[]byte {
	payloadByte := []byte(payload)
	lenBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(lenBuf, uint32(len(payloadByte)))

	packet := append([]byte{byte(msgType)}, lenBuf...)
	packet = append(packet, payloadByte...)
	return &packet
}
