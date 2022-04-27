package parser

import (
	"io"
	"os"
)

type Result struct {
	Word   string
	Pinyin string
	Count  int
}

type IParser interface {
	Parse(r *os.File) ([]Result, error)
}

type BaseParser struct {
	reader io.ReadSeeker
	pos    int64
	size   int64
}

func (p *BaseParser) read(num int64) ([]byte, error) {
	var data = make([]byte, num)
	n, err := p.reader.Read(data)
	if err != nil {
		return nil, err
	}

	p.pos += int64(n)
	return data[:n], err
}

func (p *BaseParser) seek(num int64) {
	p.reader.Seek(num, 0)
	p.pos = num
}

func (p *BaiduParser) isFinished() bool {
	return p.pos >= p.size
}

func (p *BaseParser) readString(num int64) string {
	bytes, err := p.read(num)
	if err != nil {
		return ""
	}
	return bytes2string(bytes)
}

func (p *BaseParser) readUint32() (uint32, error) {
	bytes, err := p.read(4)
	if err != nil {
		return 0, err
	}
	val := byte2Uint32(bytes)
	return val, err
}

func (p *BaseParser) readUint16() (uint16, error) {
	bytes, err := p.read(2)
	if err != nil {
		return 0, err
	}
	val := bytesUint16(bytes)
	return val, err
}
