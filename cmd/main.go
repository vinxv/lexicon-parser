package main

import (
	"flag"
	"fmt"
	"lexiconparser/parser"
	"log"
	"os"
	"path/filepath"
)

func guessKind(filename string) string {
	ext := filepath.Ext(filename)
	log.Printf("input %s. ext:%s", filename, ext)

	switch ext {
	case ".scel":
		return "sogou"
	case ".qpyd":
		return "qq"
	case ".bdict":
		return "baidu"
	}
	return ""
}

func main() {
	var input string
	var output string
	var kind string

	flag.StringVar(&input, "i", "", "input file path[required]")
	flag.StringVar(&kind, "t", "", "kind. eg: qq|sogou|baidu")
	flag.StringVar(&output, "o", "", "output file path")

	flag.Parse()
	if input == "" {
		flag.Usage()
		os.Exit(1)
	}

	if kind == "" {
		kind = guessKind(input)
	}

	var p parser.IParser
	switch kind {
	case "baidu":
		p = parser.NewBaiduParser()
	case "sogou":
		p = new(parser.SogouParser)
	case "qq":
		p = new(parser.QQParser)
	default:
		panic("unknown parser kind")
	}

	f, err := os.Open(input)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var out *os.File = os.Stdout
	if output != "" {
		out, err = os.Create(output)
		if err != nil {
			panic(err)
		}
		defer out.Close()
	}

	words, err := p.Parse(f)
	if err != nil {
		panic(err)
	}

	for _, word := range words {
		fmt.Fprintf(out, "%s\t%s\t%d\n", word.Word, word.Pinyin, word.Count)
	}
	text := fmt.Sprintf("Done. Parsed %d words from %s.", len(words), input)
	if output != "" {
		text += "output to " + output
	}
	log.Println(text)
}
