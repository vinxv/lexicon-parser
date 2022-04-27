package parser

import (
	"encoding/binary"
	"unicode/utf16"
)

func bytes2string(bytes []byte) string {
	var vals []uint16
	for i := 0; i < len(bytes)-1; i += 2 {

		val := binary.LittleEndian.Uint16(bytes[i : i+2])
		vals = append(vals, val)
	}
	runes := utf16.Decode(vals)
	return string(runes)
}

func bytesUint16(bytes []byte) uint16 {
	return binary.LittleEndian.Uint16(bytes)
}

func byte2Uint32(bytes []byte) uint32 {
	return binary.LittleEndian.Uint32(bytes)
}
