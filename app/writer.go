package main

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/codecrafters-io/git-starter-go/app/consts"
	"os"
	"path/filepath"
	"strconv"
)

type ObjectWriter struct {
	content []byte

	hash      string
	buf       *bytes.Buffer
	cryptoBuf *bytes.Buffer
}

func NewObjectWriter(content []byte) *ObjectWriter {
	return &ObjectWriter{
		content:   content,
		buf:       &bytes.Buffer{},
		cryptoBuf: &bytes.Buffer{},
	}
}

func (w *ObjectWriter) HashObject() {
	w.calculateHash()
	w.createAndFile()
	w.printHash()
}

func (w *ObjectWriter) calculateHash() {
	buf := &bytes.Buffer{}
	buf.Grow(len(w.content) + 128)

	buf.WriteString("blob")
	buf.WriteByte(consts.SpaceASCIIByte)
	buf.WriteString(strconv.Itoa(len(w.content)))
	buf.WriteByte(0x00)
	buf.Write(w.content)

	b := sha1.Sum(buf.Bytes())
	w.hash = hex.EncodeToString(b[:])
	w.buf = buf
}

func (w *ObjectWriter) createAndFile() {
	if err := os.MkdirAll(w.objectPathPrefix(), 0755); err != nil {
		panic(err)
	}
	w.encodeContent()
	if err := os.WriteFile(filepath.Join(w.objectPathPrefix(), w.objectFileName()), w.cryptoBuf.Bytes(), 0644); err != nil {
		panic(err)
	}
}

func (w *ObjectWriter) objectPathPrefix() string {
	return filepath.Join(".git", "objects", w.hash[0:2])
}

func (w *ObjectWriter) objectFileName() string {
	return w.hash[2:]
}
func (w *ObjectWriter) encodeContent() {
	ww := zlib.NewWriter(w.cryptoBuf)
	if _, err := ww.Write(w.buf.Bytes()); err != nil {
		panic(err)
	}
	_ = ww.Close()
}

func (w *ObjectWriter) printHash() {
	_, _ = fmt.Fprintf(os.Stdout, w.hash)
}
