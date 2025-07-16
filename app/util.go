package main

import (
	"encoding/base64"
	"fmt"
)

func PrintBytes(b []byte) {
	fmt.Println(base64.StdEncoding.EncodeToString(b))
}
