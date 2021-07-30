package main

import (
	"fmt"

	"github.com/jlaffaye/ftp"
)

func main() {
	f, err := ftp.Dial("five.sh:443")
	fmt.Println(f, err)
	fmt.Println(f.List("."))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
