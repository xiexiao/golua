package main

import (
	"io"
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
	fs, err := os.OpenFile(nname, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	defer fs.Close()
	if _, err = fs.Seek(0, os.SEEK_SET); err != nil {
		log.Fatal(err)
	}
	old, err := io.ReadAll(fs)
	if len(old) > 0 {
		if err = gen.CheckHeader(string(old)); err != nil {
			log.Fatal(err)
		}
		if _, err = fs.Seek(0, os.SEEK_SET); err != nil {
			log.Fatal(err)
		}
		if err = fs.Truncate(0); err != nil {
			log.Fatal(err)
		}
	}
	if _, err = fs.Write(out); err != nil {
		log.Fatal(err)
	}
}
