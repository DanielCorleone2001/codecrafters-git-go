package tree

type Tree struct {
	size int

	Entries []*Entry
}

type Entry struct {
	Mode string
	Name string
	Sha  [20]byte
}
