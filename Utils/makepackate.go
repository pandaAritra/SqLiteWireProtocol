package utils

import (
	"encoding/binary"
)

func MakePKT(rowNo int, row string) *[]byte {

	payloadByte := []byte(row)
	lenBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(lenBuf, uint32(len(payloadByte)))

	packet := append([]byte{byte(rowNo)}, lenBuf...)
	packet = append(packet, payloadByte...)
	return &packet

}
