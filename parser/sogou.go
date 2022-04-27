package parser

import (
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

func (p *SogouParser) parseMetaInfo() metaInfo {
	meta := metaInfo{}
	meta.title = p.readString(sg_titleend - p.pos)
	meta.category = p.readString(sg_typeend - p.pos)
	meta.desc = p.readString(sg_descend - p.pos)
	meta.sample = p.readString(sg_pinyinstart - p.pos)
	return meta
}

func (p *SogouParser) parsePinyinTable() []string {

	total, _ := p.readUint16()
	pinyinTable := make([]string, total)

	p.read(2)

	for i := 0; i < int(total); i++ {

		index, _ := p.readUint16()
		length, _ := p.readUint16()

		pinyin := p.readString(int64(length))
		pinyinTable[index] = pinyin
	}

	return pinyinTable
}

func (p *SogouParser) parseWord(pinyinTable []string) ([]Result, error) {

	sameNum, _ := p.readUint16()
	tableLen, _ := p.readUint16()

	pinyinChars := make([]string, tableLen)
	pinyinIndexByte, err := p.read(int64(tableLen))
	if err != nil {
		return nil, err
	}

	for i := 0; i < int(tableLen/2); i++ {
		indexByte := pinyinIndexByte[2*i : 2*i+2]
		index := bytesUint16(indexByte)
		pinyinChars[i] = pinyinTable[index]
	}

	pinyin := strings.Join(pinyinChars, "'")
	pinyin = strings.Trim(pinyin, "'")

	wordItems := make([]Result, sameNum)
	for i := 0; i < int(sameNum); i++ {

		wordLen, _ := p.readUint16()
		word := p.readString(int64(wordLen))

		extLen, _ := p.readUint16()
		extBytes, err := p.read(int64(extLen))
		if err != nil {
			return nil, err
		}
		countBytes := extBytes[:2]
		count := bytesUint16(countBytes)
		wordItems[i] = Result{word, pinyin, int(count)}
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

	p.seek(sg_titlestart)

	meta := p.parseMetaInfo()
	_ = meta
	// log.Printf("metainfo\ntitle: %s\ncategoty:%s\ndescription:%s\nsample:%v\n\n", meta.title, meta.category, meta.desc, meta.sample)

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
