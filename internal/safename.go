package internal

import "path"

// isSafeName returns true iff the file name given fulfills all of the
// following:
//
//     - is just a file basename with no path
//     - does not refer to a hidden file (begin with dot)
//     - does not begin with a dash (that could cause human mistakes)
func IsSafeName(name string) bool {
	dir, base := path.Split(name)
	if dir != "" {
		return false
	}
	if base == "" {
		return false
	}
	switch base[0] {
	case '.', '-':
		return false
	}
	return true
}
