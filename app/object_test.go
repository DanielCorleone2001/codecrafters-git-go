package main

import (
	"fmt"
	"testing"
)

func Test_Parse(t *testing.T) {
	p := NewFileParser("3b18e512dba79e4c8300dd08aeb37f8e728b8dad")
	p.withAbs = true
	obj := p.ParseObject()
	fmt.Println(obj.String())
}
