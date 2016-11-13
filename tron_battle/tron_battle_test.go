package main

import (
	"io"
	"io/ioutil"
	"os"
	"testing"
)

func BenchmarkGame(b *testing.B) {

	in, _ := ioutil.TempFile("", "")
	defer in.Close()

	io.WriteString(in, "3 0\n"+"2 1 2 1\n"+"13 3 13 3\n"+"5 6 5 6\n")

	for i := 0; i < b.N; i++ {
		in.Seek(0, os.SEEK_SET)
		runGame(in)
	}
}
