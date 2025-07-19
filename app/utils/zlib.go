package utils

import (
	"bytes"
	"compress/zlib"
	"io"
)

func ZlibDecode(input []byte) *bytes.Buffer {
	r, err := zlib.NewReader(bytes.NewReader(input))
	if err != nil {
		panic(err)
	}
	defer r.Close()

	b := &bytes.Buffer{}
	_, err = io.Copy(b, r)
	if err != nil {
		panic(err)
	}
	return b
}
