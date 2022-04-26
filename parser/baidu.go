package parser

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"strings"
)

var shengmu = []string{"c", "d", "b", "f", "g", "h", "ch", "j", "k", "l", "m", "n", "", "p", "q", "r", "s", "t", "sh", "zh", "w", "x", "y", "z"}
var yunmu = []string{"uang", "iang", "iong", "ang", "eng", "ian", "iao", "ing", "ong", "uai", "uan", "ai", "an", "ao", "ei", "en", "er", "ua", "ie", "in", "iu", "ou", "ia", "ue", "ui", "un", "uo", "a", "e", "i", "o", "u", "v"}
var BD_HEAD_FLAG = []byte("biptbdsw")
var BD_HEADPOS int64 = 0x60
var BD_DATAPOS int64 = 0x350

/// BaiduParser 百度词典解析工具
/// binary layout
///                  |          |
///              0x60|     0x350|
/// +--------+-------+----+-----+----+---+-----------+-----------+
/// |   8B   |       | 4B |     | 4B | 1B|   LEN*2   |   LEN*2   |
/// +--------+-------+----+-----+----+---+-----------+-----------+
/// |  FLAG  |       | END|     | LEN| F | pinyininex|   WORD    |
/// +--------+-------+----+-----+----+---+-----------+-----------+
///
///                             |      REPEAT                    |
///                             +------------------------------- +
type BaiduParser struct {
	BaseParser
}

func NewBaiduParser() *BaiduParser {
	bp := BaseParser{
		reader: nil,
		pos:    0,
	}
	return &BaiduParser{bp}
}

func (p *BaiduParser) Parse(reader *os.File) ([]Result, error) {
	p.reader = reader
	size, err := reader.Stat()
	if err != nil {
		return nil, err
	}
	p.size = size.Size()

	flag, err := p.read(8)
	if err != nil {
		return nil, err
	}

	if !bytes.Equal(flag, BD_HEAD_FLAG) {
		return nil, fmt.Errorf("not baidufile")
	}

	p.seek(BD_HEADPOS)

	endPosBytes, err := p.read(4)
	if err != nil {
		return nil, err
	}

	endPos := binary.LittleEndian.Uint32(endPosBytes)
	p.seek(BD_DATAPOS)

	var result []Result = make([]Result, 0)

	for p.pos < int64(endPos) {

		lenBytes, err := p.read(4)
		if err != nil {
			return nil, err
		}

		len := binary.LittleEndian.Uint16(lenBytes)
		data, err := p.read(int64(len * 4))
		if err != nil {
			return nil, err
		}

		s := data[:len*2]

		var pinyin string

		for i := 0; i < int(len)*2; i++ {
			b := s[i : i+1]
			v, _ := binary.Uvarint(b)
			var c string
			if i%2 == 0 {
				c = shengmu[v]
			} else {
				c = yunmu[v] + "'"
			}
			pinyin += c
		}
		pinyin = strings.Trim(pinyin, "'")
		word := bytes2string(data[len*2:])
		result = append(result, Result{word, pinyin, 1})
	}
	return result, nil
}

var _ IParser = new(BaiduParser)
