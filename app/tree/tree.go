package tree

import "fmt"

type Tree struct {
	size int

	Entries []*Entry
}

type Entry struct {
	Mode string
	Name string
	Sha  [20]byte
}

func (e *Entry) String(nameOnly bool) string {
	if nameOnly {
		return e.Name
	}
	return fmt.Sprintf("%s %s %x    %s", e.Mode, e.Mode, e.Sha, e.Name)
}

func (e *Entry) typeName() string {
	switch e.Mode {
	case "100644":
		return "blob"
	case "40000":
		return "tree"
	default:
		return "blob" // regular file
	}
}
