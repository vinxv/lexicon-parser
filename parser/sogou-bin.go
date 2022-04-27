package parser

import (
	"fmt"
	"log"
	"os"
)

type SogouBinParser struct {
	BaseParser
}

func (p *SogouBinParser) Parse(r *os.File) ([]Result, error) {
	p.reader = r
	size, err := r.Stat()
	if err != nil {
		return nil, err
	}
	p.size = size.Size()

	// file header
	fileChecksum, err := p.readUint32()
	if err != nil {
		return nil, err
	}

	v0, _ := p.readUint32()
	v1, _ := p.readUint32()
	v2, _ := p.readUint32()
	v3, _ := p.readUint32()
	log.Printf("v0:%d, v1:%d, v2:%d, v3:%d", v0, v1, v2, v3)
	configSize := v0
	checksum := v0 + v1 + v2 + v3

	if v0 > uint32(p.size) {
		return nil, fmt.Errorf("invalid config size: %d", v0)
	}

	log.Printf("file checksum: %d, configSize:%d, checksum:%d", fileChecksum, configSize, checksum)

	p.seek(20)

	for i := 0; i < int(v1); i++ {
		dictTypeDef, _ := p.readUint16()
		numDataType, _ := p.readUint16()
		dataTypes := []uint16{}
		if dictTypeDef > 100 {
			return nil, fmt.Errorf("invalid dict type def: %d", dictTypeDef)
		}
		for j := 0; j < int(numDataType); j++ {
			dataType, _ := p.readUint16()
			dataTypes = append(dataTypes, dataType)
			log.Printf("[%d-%d] data type: %d", i, j, dataType)
		}
		attrIdx, _ := p.readUint32()
		keyDataIdx, _ := p.readUint32()
		dataIdx, _ := p.readUint32()
		v6, _ := p.readUint32()
		log.Printf("attrIdx: %d, keyDataIdx: %d, dataIdx: %d, v2:%d", attrIdx, keyDataIdx, dataIdx, v6)
	}

	for i := 0; i < int(v2); i++ {
		count, _ := p.readUint32()
		a2, _ := p.readUint32()
		id, _ := p.readUint32()
		b2, _ := p.readUint32()
		log.Printf("count: %d, a2: %d, id: %d, b2: %d", count, a2, id, b2)
	}

	for i := 0; i < int(v3); i++ {
		aint, _ := p.readUint32()
		log.Printf("[%d] aint: %d", i, aint)
	}

	if v0+8 != uint32(p.pos) {
		return nil, fmt.Errorf("invalid pos: %d", p.pos)
	}

	headerSize := 12*(v2+v3+v1) + 24
	log.Printf("header size: %d", headerSize)

	b2Version, _ := p.readUint32()
	b2Fmt, _ := p.readUint32()

	log.Printf("b2Version: %d, b2Fmt: %d", b2Version, b2Fmt)

	totalSize, _ := p.readUint32()
	var userDictSize uint32 = 4 + 76
	if totalSize+headerSize+configSize+8 != uint32(p.size)-userDictSize {
		return nil, fmt.Errorf("invalid total size: %d", totalSize)
	}

	b2Size3, _ := p.readUint32()
	b2Size4, _ := p.readUint32()
	b2Size5, _ := p.readUint32()
	log.Printf("b2Size3: %d, b2Size4: %d, b2Size5: %d\n", b2Size3, b2Size4, b2Size5)

	for i := 0; i < int(b2Size3); i++ {
		offset, _ := p.readUint32()
		dataSize, _ := p.readUint32()
		usedDataSize, _ := p.readUint32()
		checksum += offset + dataSize + usedDataSize
	}

	for i := 0; i < int(b2Size4); i++ {
		offset, _ := p.readUint32()
		dataSize, _ := p.readUint32()
		usedDataSize, _ := p.readUint32()
		checksum += offset + dataSize + usedDataSize
	}

	for i := 0; i < int(b2Size5); i++ {
		offset, _ := p.readUint32()
		dataSize, _ := p.readUint32()
		usedDataSize, _ := p.readUint32()
		checksum += offset + dataSize + usedDataSize
	}

	log.Printf("current pos: %d", p.pos)
	if v0+headerSize+8 != uint32(p.pos) {
		return nil, fmt.Errorf("invalid pos: %d", p.pos)
	}

	// currentPos := p.pos

	p.seek(p.size - 0x4c)
	var p2 uint32
	var p3 uint32
	for i := 0; i < 19; i++ {
		v, _ := p.readUint32()
		if i == 14 {
			p2 = v
		}
		if i == 15 {
			p3 = v
		}
	}
	log.Printf("p2: %d, p3: %d", p2, p3)

	return nil, nil
}
