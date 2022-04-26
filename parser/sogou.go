package parser

import (
	"encoding/binary"
	"log"
	"os"
	"strings"
)

const (
	sg_titlestart   int64 = 0x130
	sg_titleend     int64 = 0x338
	sg_typeend      int64 = 0x540
	sg_descend      int64 = 0xd40
	sg_pinyinstart  int64 = 0x1540
	sg_chinesestart int64 = 0x2628
)

type metaInfo struct {
	title    string
	category string
	desc     string
	sample   string
}

type SogouParser struct {
	BaseParser
}

func (p *SogouParser) byte2int(bytes []byte) int {
	return int(binary.LittleEndian.Uint16(bytes))
}

func (p *SogouParser) readeInt() int {
	bytes, _ := p.read(2)
	return int(binary.LittleEndian.Uint16(bytes))
}

func (p *SogouParser) parseMetaInfo() metaInfo {
	meta := metaInfo{}
	meta.title = p.readString(sg_titleend - p.pos)
	meta.category = p.readString(sg_typeend - p.pos)
	meta.desc = p.readString(sg_descend - p.pos)
	meta.sample = p.readString(sg_pinyinstart - p.pos)
	return meta
}

func (p *SogouParser) parsePinyinTable() []string {

	total := p.readeInt()
	pinyinTable := make([]string, total)

	p.read(2)

	for i := 0; i < total; i++ {

		index := p.readeInt()
		length := p.readeInt()

		pinyin := p.readString(int64(length))
		pinyinTable[index] = pinyin
	}

	return pinyinTable
}

func (p *SogouParser) parseWord(pinyinTable []string) ([]Result, error) {

	sameNum := p.readeInt()
	tableLen := p.readeInt()

	pinyinChars := make([]string, tableLen)
	pinyinIndexByte, err := p.read(int64(tableLen))
	if err != nil {
		return nil, err
	}

	for i := 0; i < tableLen/2; i++ {
		indexByte := pinyinIndexByte[2*i : 2*i+2]
		index := p.byte2int(indexByte)
		pinyinChars[i] = pinyinTable[index]
	}

	pinyin := strings.Join(pinyinChars, "'")
	pinyin = strings.Trim(pinyin, "'")

	wordItems := make([]Result, sameNum)
	for i := 0; i < sameNum; i++ {

		wordLen := p.readeInt()
		word := p.readString(int64(wordLen))

		extLen := p.readeInt()
		_ = extLen

		extBytes, err := p.read(int64(extLen))
		if err != nil {
			return nil, err
		}
		countBytes := extBytes[:2]
		count := p.byte2int(countBytes)
		wordItems[i] = Result{word, pinyin, count}
	}
	return wordItems, nil
}

func (p *SogouParser) Parse(reader *os.File) ([]Result, error) {
	p.reader = reader
	size, err := reader.Stat()
	if err != nil {
		return nil, err
	}
	p.size = size.Size()

	meta := p.parseMetaInfo()
	log.Printf("metainfo\ntitle: %s\ncategoty:%s\ndescription:%s\nsample:%v\n\n", meta.title, meta.category, meta.desc, meta.sample)

	pinyinTable := p.parsePinyinTable()
	words := []Result{}
	for p.pos < p.size {
		wordItems, err := p.parseWord(pinyinTable)
		if err != nil {
			continue
		}
		words = append(words, wordItems...)
	}
	return words, nil
}
