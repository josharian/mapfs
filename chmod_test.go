package mapfs

import (
	"fmt"
	"io/fs"
	"strings"
	"testing"
)

func TestToDirMode(t *testing.T) {
	tests := []struct {
		in   fs.FileMode
		want fs.FileMode
	}{
		{in: 0666, want: 0777 | fs.ModeDir},
		{in: 0664, want: 0775 | fs.ModeDir},
		{in: 0604, want: 0705 | fs.ModeDir},
		{in: 0004, want: 0005 | fs.ModeDir},
		{in: 0600, want: 0700 | fs.ModeDir},
		{in: 0060, want: 0070 | fs.ModeDir},
	}

	for _, test := range tests {
		got := toDirMode(test.in)
		if test.want != got {
			t.Errorf("toDirMode(%v) = %v want %v", test.in, got, test.want)
		}
	}
}

func TestChmodAll(t *testing.T) {
	m := MapFS{
		"a/b/c.txt": &MapFile{},
	}
	m.ChmodAll(0666)
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
a: drwxrwxrwx
a/b: drwxrwxrwx
a/b/c.txt: -rw-rw-rw-
`[1:]
	got := buf.String()
	if want != got {
		t.Errorf("MapFS modes want:\n%s\ngot:\n%s\n", want, got)
	}
}
