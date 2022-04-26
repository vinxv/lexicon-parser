package parser

import (
	"encoding/binary"
)

func bytes2string(bytes []byte) string {
	var rs []rune
	for i := 0; i < len(bytes)-1; i += 2 {

		val := binary.LittleEndian.Uint16(bytes[i : i+2])
		if val == 0 {
			continue
		}
		rs = append(rs, rune(val))
	}
	return string(rs)
}
