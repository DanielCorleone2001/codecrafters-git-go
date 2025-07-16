package main

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func Hash2FilePath(hash string) string {
	return filepath.Join(".git/objects", hash[0:2], hash[2:])
}

type Object interface {
	fmt.Stringer

	ObjectName() string
}

type Blob struct {
	b []byte
}

func (b *Blob) String() string {
	r, err := zlib.NewReader(bytes.NewReader(b.b))
	if err != nil {
		panic(err)
	}
	sb := &strings.Builder{}
	io.Copy(sb, r)
	_ = r.Close()
	return sb.String()
}

func (*Blob) ObjectName() string {
	return "blob"
}

func ParseObjectFile(fc []byte) Object {
	b := bytes.NewBuffer(fc)
	// blob <size>\0<content>
	objectType, err := b.ReadString(0x20)
	if err != nil {
		panic(err)
	}
	switch objectType {
	case "blob":
		_, _ = b.ReadByte()
		s, err := b.ReadString(0x00)
		if err != nil {
			panic(err)
		}
		size, err := strconv.Atoi(s)
		if err != nil {
			panic(err)
		}
		_, _ = b.ReadByte()
		content := make([]byte, size)
		n, err := b.Read(content)
		if err != nil {
			panic(err)
		}
		if n != size {
			_, _ = fmt.Fprintf(os.Stderr, "read bytes len is %d,want:%d", n, size)
		}
		return &Blob{b: content}
	default:
		return nil
	}
}
