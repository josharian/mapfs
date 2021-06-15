// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mapfs

import (
	"fmt"
	"io/fs"
	"strings"
	"testing"
	"testing/fstest"
)

func TestMapFS(t *testing.T) {
	m := MapFS{
		"hello":             {Data: []byte("hello, world\n")},
		"fortune/k/ken.txt": {Data: []byte("If a program is too slow, it must have a loop.\n")},
	}
	if err := fstest.TestFS(m, "hello", "fortune/k/ken.txt"); err != nil {
		t.Fatal(err)
	}
}

func TestMapFSChmodDot(t *testing.T) {
	m := MapFS{
		"a/b.txt": &MapFile{Mode: 0666},
		".":       &MapFile{Mode: 0777 | fs.ModeDir},
	}
	buf := new(strings.Builder)
	fs.WalkDir(m, ".", func(path string, d fs.DirEntry, err error) error {
		fi, err := d.Info()
		if err != nil {
			return err
		}
		fmt.Fprintf(buf, "%s: %v\n", path, fi.Mode())
		return nil
	})
	want := `
.: drwxrwxrwx
a: d---------
a/b.txt: -rw-rw-rw-
`[1:]
	got := buf.String()
	if want != got {
		t.Errorf("MapFS modes want:\n%s\ngot:\n%s\n", want, got)
	}
}
