package mapfs

import (
	"io/fs"
)

func toDirMode(mode fs.FileMode) fs.FileMode {
	// set each x bit if the corresponding r or w bit is set
	for shift := 0; shift < 9; shift += 3 {
		perm := (mode >> shift) & 0b111
		perm |= (perm >> 1 & 1) | (perm >> 2 & 1) // set x bit if r|w is set
		mode |= perm << shift
	}
	return mode | fs.ModeDir
}

// ChmodAll changes the files and directories in m to have mode mode.
// The mode will be adapted for directories: if mode is 0640, the mode for directories will be 0750|fs.ModeDir.
func (m MapFS) ChmodAll(fileMode fs.FileMode) {
	dirMode := toDirMode(fileMode)
	var missing []string
	err := fs.WalkDir(m, ".", func(path string, d fs.DirEntry, err error) error {
		f, ok := m[path]
		if !ok {
			missing = append(missing, path)
			return nil
		}
		if d.IsDir() {
			f.Mode = dirMode
		} else {
			f.Mode = fileMode
		}
		return nil
	})
	if err != nil {
		panic(err) // unreachable
	}
	for _, dir := range missing {
		m[dir] = &MapFile{
			Mode: dirMode,
		}
	}
}
