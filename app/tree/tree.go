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
	return fmt.Sprintf("%d %s %s    %s", e.filePermission(), e.Mode, e.Sha, e.Name)
}

func (e *Entry) filePermission() int {
	switch e.Mode {
	case "blob":
		return 100644
	case "tree":
		return 40000
	default:
		return 100644 // regular file
	}
}
