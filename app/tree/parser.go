package tree

import (
	"bytes"
	"fmt"
	"github.com/codecrafters-io/git-starter-go/app/consts"
	"github.com/codecrafters-io/git-starter-go/app/utils"
	"io"
	"os"
	"path/filepath"
	"strconv"
)

type FileParser struct {
	treeSha    string
	size       int
	fileBuf    *bytes.Buffer
	entriesBuf *bytes.Buffer

	t *Tree
}

func NewTreeFileParser(sha string) *FileParser {
	return &FileParser{
		treeSha:    sha,
		fileBuf:    &bytes.Buffer{},
		entriesBuf: &bytes.Buffer{},
	}
}

func (p *FileParser) ToTree() *Tree {
	p.readFileContent()
	p.parserHeader()
	p.parseEntries()
	p.t.size = p.size
	return p.t
}

func (p *FileParser) filePath() string {
	return filepath.Join(".git", "objects", p.treeSha[0:2])
}

func (p *FileParser) fileName() string {
	return p.treeSha[2:]
}

func (p *FileParser) readFileContent() {
	b, err := os.ReadFile(filepath.Join(p.filePath(), p.fileName()))
	if err != nil {
		panic(err)
	}
	p.fileBuf = bytes.NewBuffer(b)
}

func (p *FileParser) parserHeader() {
	mode := utils.ReadStringUtil(p.fileBuf, consts.SpaceASCIIByte)
	if mode != "tree" {
		panic("tree file mode is not tree,it's [" + mode + "]")
	}
	_, _ = p.fileBuf.ReadByte() // consts.SpaceASCIIByte
	sizes := utils.ReadStringUtil(p.fileBuf, consts.NullASCIIByte)
	size, err := strconv.Atoi(sizes)
	if err != nil {
		panic(fmt.Sprintf("conv sizes to int fail, sizes:%s, err:%s", sizes, err))
	}
	p.size = size
	_, _ = p.fileBuf.ReadByte()
}

func (p *FileParser) parseEntries() {
	p.readEntries2Buf()

	read := 0
	for read < p.size {
		read += p.readSingleEntry()
	}

}

func (p *FileParser) readEntries2Buf() {
	eb := make([]byte, p.size)
	n, err := io.ReadFull(p.fileBuf, eb)
	if err != nil {
		panic(err)
	}
	if n != p.size {
		panic(fmt.Errorf("read:%d, want:%d", n, p.size))
	}
	p.entriesBuf = bytes.NewBuffer(eb)
}

func (p *FileParser) readSingleEntry() int {
	mode := utils.ReadStringUtil(p.entriesBuf, consts.SpaceASCIIByte)
	_, _ = p.entriesBuf.ReadByte()
	name := utils.ReadStringUtil(p.entriesBuf, consts.NullASCIIByte)
	_, _ = p.entriesBuf.ReadByte()

	sha := make([]byte, 20)
	n, err := p.entriesBuf.Read(sha)
	if err != nil {
		panic(err)
	}
	if n != 20 {
		panic(fmt.Errorf("read:%d, want:%d", n, 20))
	}

	p.t.Entries = append(p.t.Entries, &Entry{
		Mode: mode,
		Name: name,
		Sha:  [20]byte(sha),
	})

	return len(mode) + 1 + len(name) + 20
}
