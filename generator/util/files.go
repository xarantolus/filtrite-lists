package util

import "path"

func StripExtension(p string) string {
	name := path.Base(p)
	ext := path.Ext(name)
	if ext != "" {
		name = name[:len(name)-len(ext)]
	}
	return name
}
