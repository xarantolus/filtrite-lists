package util

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"unicode"
)

func GetFilterListName(f io.Reader) (name string, err error) {
	scan := bufio.NewScanner(f)

	// Basically search for lines like "! Title: XY's Annoyance List"
	for scan.Scan() {
		txt := strings.TrimSpace(scan.Text())

		if len(txt) == 0 || !strings.HasPrefix(txt, "!") {
			continue
		}

		trimmed := strings.TrimLeftFunc(txt, func(r rune) bool {
			return r == '!' || unicode.IsSpace(r)
		})

		if !strings.HasPrefix(strings.ToLower(trimmed), "title") {
			continue
		}

		split := strings.SplitN(trimmed, ":", 2)
		if len(split) != 2 {
			continue
		}

		name = strings.TrimSpace(split[1])

		return
	}

	if err = scan.Err(); err != nil {
		return
	}

	err = fmt.Errorf("no name/title found in filter list")

	return
}

func GetFilterListNameFromURL(url string) (name string, err error) {
	resp, err := req(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	return GetFilterListName(resp.Body)
}
