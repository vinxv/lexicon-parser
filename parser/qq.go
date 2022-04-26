package parser

import (
	"compress/zlib"
	"encoding/binary"
	"io/ioutil"
	"log"
	"os"
)

type QQParser struct {
	BaseParser
}

func (p *QQParser) Parse(reader *os.File) ([]Result, error) {
	p.reader = reader
	size, err := reader.Stat()
	if err != nil {
		return nil, err
	}
	p.size = size.Size()

	p.seek(0x38)
	bytes, err := p.read(4)
	if err != nil {
		return nil, err
	}

	pos := int64(binary.LittleEndian.Uint32(bytes))

	p.seek(0x44)

	bytes, err = p.read(4)
	if err != nil {
		return nil, err
	}

	// word count
	count := binary.LittleEndian.Uint32(bytes)

	p.seek(0x60)
	metaInfo := p.readString(pos - 0x60)
	log.Printf("meta info:\n %s\ncount:%d", metaInfo, count)
	p.seek(pos)

	unzipReader, err := zlib.NewReader(p.reader)
	if err != nil {
		log.Fatal(err)
	}

	data, err := ioutil.ReadAll(unzipReader)
	if err != nil {
		return nil, err
	}

	var addr int

	words := make([]Result, count)

	for i := 0; i < int(count); i++ {

		b1 := data[addr] & 0xff
		pinyinLen := int(b1)

		b2 := data[addr+1] & 0xff
		wordLen := int(b2)

		p1 := data[addr+6 : addr+10]
		pinyinAddr := binary.LittleEndian.Uint32(p1)

		wordAddr := int(pinyinAddr) + pinyinLen

		pinyinData := data[pinyinAddr:wordAddr]
		pinyin := string(pinyinData)

		wordData := data[wordAddr : wordAddr+wordLen]
		word := bytes2string(wordData)

		words[i] = Result{Pinyin: pinyin, Word: word}
		addr += 10
	}

	return words, nil
}
