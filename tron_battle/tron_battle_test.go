/*
Copyright (c) 2016 Hubert Jarosz

This software is provided 'as-is', without any express or implied
warranty. In no event will the authors be held liable for any damages
arising from the use of this software.

Permission is granted to anyone to use this software for any purpose,
including commercial applications, and to alter it and redistribute it
freely, subject to the following restrictions:

1. The origin of this software must not be misrepresented; you must not
   claim that you wrote the original software. If you use this software
   in a product, an acknowledgement in the product documentation would be
   appreciated but is not required.
2. Altered source versions must be plainly marked as such, and must not be
   misrepresented as being the original software.
3. This notice may not be removed or altered from any source distribution.
*/

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

	io.WriteString(in, "4 1\n"+"2 1 2 1\n"+"13 3 13 3\n"+"5 6 5 6\n"+"11 19 11 19\n")

	for i := 0; i < b.N; i++ {
		in.Seek(0, os.SEEK_SET)
		runGame(in)
	}
}
