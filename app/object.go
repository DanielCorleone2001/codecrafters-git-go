package main

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
)

type ObjectParser struct {
	hash string

	objFilePath string
	fc          []byte

	b       *bytes.Buffer
	withAbs bool
}

const (
	hashLen = 40
)

func NewFileParser(hash string) *ObjectParser {
	if len(hash) != hashLen {
		panic(fmt.Sprintf("hash len must be:%d, instead of :%d", hashLen, len(hash)))
	}
	return &ObjectParser{
		hash: hash,
	}
}

func (p *ObjectParser) ParseObject() Object {
	p.parseFilePath()
	p.readFile()
	p.decodeFileContent()
	return p.toObject()
}

func (p *ObjectParser) parseFilePath() {
	path := filepath.Join(".git/objects", p.hash[0:2], p.hash[2:])
	if p.withAbs {
		path = filepath.Join("/Users/daniel/code/go/codecrafter/codecrafters-git-go", path)
	}
	p.objFilePath = path
}

func (p *ObjectParser) readFile() {
	fc, err := os.ReadFile(p.objFilePath)
	if err != nil {
		panic(err)
	}
	p.fc = fc
}

func (p *ObjectParser) decodeFileContent() {
	r, err := zlib.NewReader(bytes.NewReader(p.fc))
	if err != nil {
		panic(err)
	}
	defer r.Close()

	buf := &bytes.Buffer{}
	_, err = io.Copy(buf, r)
	if err != nil {
		panic(err)
	}

	p.b = buf
}

type Object interface {
	fmt.Stringer
	ObjectName() string
}

type Blob struct {
	b []byte
}

func (b *Blob) String() string {
	return string(b.b)
}

func (b *Blob) ObjectName() string {
	return "blob"
}

func (p *ObjectParser) readUntil(delim byte) []byte {
	l := bytes.IndexByte(p.b.Bytes(), delim)
	b := make([]byte, l)
	n, err := io.ReadFull(p.b, b)
	if err != nil {
		panic(err)
	}
	if n != l {
		panic(fmt.Sprintf("read:%d, want:%d", n, l))
	}
	return b
}

func (p *ObjectParser) parseObjectType() string {
	t := p.readUntil(0x20)
	return string(t)
}

func (p *ObjectParser) toObject() Object {
	switch objectType := p.parseObjectType(); objectType {
	case "blob":
		_, _ = p.b.ReadByte()
		s := p.readUntil(0x00)
		size, err := strconv.Atoi(string(s))
		if err != nil {
			panic(err)
		}
		_, _ = p.b.ReadByte()
		content := make([]byte, size)
		n, err := io.ReadFull(p.b, content)
		if err != nil {
			panic(err)
		}
		if n != size {
			panic(fmt.Errorf("read:%d, want:%d", n, size))
		}
		return &Blob{b: content}
	default:
		panic("not supported for:" + objectType)
	}
}
