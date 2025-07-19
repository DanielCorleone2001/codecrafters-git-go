package tree

import (
	"fmt"
	"os"
	"strings"
)

type ListTree struct {
	treeSha  string
	nameOnly bool

	p *FileParser
}

func NewListTreeOp() *ListTree {
	return &ListTree{}
}

func (l *ListTree) WithTreeSHA(sha string) *ListTree {
	l.treeSha = sha
	return l
}

func (l *ListTree) WithNameOnly(b bool) *ListTree {
	l.nameOnly = b
	return l
}

func (l *ListTree) ListTreeContent() {
	l.buildParser()
	l.printTreeEntries()
}

func (l *ListTree) buildParser() {
	l.p = NewTreeFileParser(l.treeSha)
}

func (l *ListTree) printTreeEntries() {
	sb := &strings.Builder{}
	t := l.p.ToTree()
	for _, e := range t.Entries {
		if l.nameOnly {
			_, _ = fmt.Fprintf(sb, "%s\n", e.Name)
		} else {
			_, _ = fmt.Fprintf(sb, "%s %s %s    %s", e.Mode, e.Mode, e.Sha, e.Name)
		}
	}
	_, _ = fmt.Fprintf(os.Stdout, sb.String())
}
