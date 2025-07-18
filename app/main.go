package main

import (
	"fmt"
	"github.com/codecrafters-io/git-starter-go/app/tree"
	"os"
	"strings"
)

// Usage: your_program.sh <command> <arg1> <arg2> ...
func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Fprintf(os.Stderr, "Logs from your program will appear here!\n")

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: mygit <command> [<args>...]\n")
		os.Exit(1)
	}

	switch command := os.Args[1]; command {
	case "init":
		// Uncomment this block to pass the first stage!
		for _, dir := range []string{".git", ".git/objects", ".git/refs"} {
			if err := os.MkdirAll(dir, 0755); err != nil {
				fmt.Fprintf(os.Stderr, "Error creating directory: %s\n", err)
			}
		}

		headFileContents := []byte("ref: refs/heads/main\n")
		if err := os.WriteFile(".git/HEAD", headFileContents, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing file: %s\n", err)
		}

		fmt.Println("Initialized git directory")
	case "cat-file":
		fileHash := os.Args[3]
		p := NewFileParser(fileHash)
		obj := p.ParseObject()
		fmt.Print(obj.String())
	case "hash-object":
		fileName := os.Args[3]
		content, err := os.ReadFile(fileName)
		if err != nil {
			panic(err)
		}
		w := NewObjectWriter(content)
		w.HashObject()
	case "ls-tree":
		l := tree.NewListTreeOp()
		if strings.Contains(os.Args[2], "name-only") {
			l.WithNameOnly(true)
			l.WithTreeSHA(os.Args[3])
		} else {
			l.WithTreeSHA(os.Args[2])
		}
		l.ListTreeContent()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command %s\n", command)
		os.Exit(1)
	}
}
