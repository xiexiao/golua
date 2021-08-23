package main

import (
	"log"
	"os"
	"strings"

	"github.com/xiexiao/golua/gen"
)

func main() {
	var err error
	if len(os.Args) == 1 {
		log.Fatal("need a file!")
	}
	fname := os.Args[1]
	data, err := os.ReadFile(fname)
	if err != nil {
		log.Fatal(err)
	}
	out, err := gen.GenLua(fname, data)
	if err != nil {
		log.Fatal(err)
	}
	nname := fname + "_lua"
	i := strings.LastIndex(fname, ".")
	if i > -1 {
		nname = fname[:i] + "_lua" + fname[i:]
	}
	_, err = os.Stat(nname)
	if err != nil {
		if !os.IsNotExist(err) {
			os.Remove(nname)
		}
	}
	err = os.WriteFile(nname, []byte(out), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
}
