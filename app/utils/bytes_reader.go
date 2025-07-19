package utils

import (
	"bytes"
	"fmt"
	"io"
)

type IBytes interface {
	Bytes() []byte
	io.Reader
}

func ReadBytesUtil(bb IBytes, delim byte) []byte {
	l := bytes.IndexByte(bb.Bytes(), delim)
	b := make([]byte, l)
	n, err := io.ReadFull(bb, b)
	if err != nil {
		panic(err)
	}
	if n != l {
		panic(fmt.Sprintf("read:%d, want:%d", n, l))
	}
	return b
}

func ReadStringUtil(bb IBytes, delim byte) string {
	return string(ReadBytesUtil(bb, delim))
}
