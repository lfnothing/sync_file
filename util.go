package sync_file

import (
	"bytes"
	"encoding/binary"
)

func Int64ToBytes(n int64) []byte {
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, n)
	return bytesBuffer.Bytes()
}

func BytesToInt64(b []byte) (n int64) {
	bytesBuffer := bytes.NewBuffer(b)
	binary.Read(bytesBuffer, binary.BigEndian, &n)
	return
}
