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

	b *bytes.Buffer
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
	p.objFilePath = filepath.Join(".git/objects", p.hash[0:2], p.hash[2:])
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
	n, err := io.Copy(buf, r)
	if err != nil {
		panic(err)
	}
	fmt.Printf("read %d bytes after decode, from path:%s\n", n, p.objFilePath)

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

func (p *ObjectParser) toObject() Object {
	objectType, err := p.b.ReadString(0x20)
	if err != nil {
		panic(err)
	}
	fmt.Printf("object type is:%s\b", objectType)
	switch objectType {
	case "blob":
		_, _ = p.b.ReadByte()
		s, err := p.b.ReadString(0x00)
		if err != nil {
			panic(err)
		}
		size, err := strconv.Atoi(s)
		if err != nil {
			panic(err)
		}
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
