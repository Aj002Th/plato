package tcp

import (
	"bytes"
	"encoding/binary"
)

// DataPkg 数据包结构
type DataPkg struct {
	Len  uint32
	Data []byte
}

func (d *DataPkg) Marshal() []byte {
	bytesBuffer := bytes.NewBuffer([]byte{})
	_ = binary.Write(bytesBuffer, binary.BigEndian, d.Len)
	return append(bytesBuffer.Bytes(), d.Data...)
}
