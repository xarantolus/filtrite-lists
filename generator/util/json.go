package util

import (
	"encoding/json"
	"os"
)

func LoadJSON(filename string, target interface{}) (err error) {
	f, err := os.Open(filename)
	if err != nil {
		return
	}
	defer f.Close()

	err = json.NewDecoder(f).Decode(target)

	return
}

func SaveJSON(filename string, source interface{}) (err error) {
	tmpFn := filename + ".tmp"

	f, err := os.Create(tmpFn)
	if err != nil {
		return
	}
	var closed bool
	defer func() {
		var cerr error
		if !closed {
			cerr = f.Close()
		}

		if err == nil {
			err = cerr
		} else {
			_ = os.Remove(tmpFn)
		}
	}()

	err = json.NewEncoder(f).Encode(source)

	err = f.Close()
	closed = true
	if err != nil {
		return
	}

	return os.Rename(tmpFn, filename)
}
