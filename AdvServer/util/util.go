package util

import (
	"bytes"
	"encoding/binary"
)

func IntToBytes(n int) []byte {
	tmp := int32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, tmp)
	return bytesBuffer.Bytes()
}

func Int64ToBytes(n int64) []byte {
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.LittleEndian, n)
	return bytesBuffer.Bytes()
}

func BytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)
	var tmp int32
	binary.Read(bytesBuffer, binary.BigEndian, &tmp)
	return int(tmp)
}

func Int2Byte(data int) []byte {
	b := []byte{0x00, 0x00, 0x03, 0xe8}
	b_buf := bytes.NewBuffer(b)
	var x int32
	binary.Read(b_buf, binary.BigEndian, &x)
	return b_buf.Bytes()
}
