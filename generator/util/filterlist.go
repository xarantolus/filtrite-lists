package util

import (
	"bufio"
	"fmt"
	"io"
	"net/url"
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

		name, ok := getName(txt)
		if ok {
			return name, nil
		}
	}

	if err = scan.Err(); err != nil {
		return
	}

	err = fmt.Errorf("no name/title found in filter list")

	return
}

func getName(line string) (name string, ok bool) {
	trimmed := strings.TrimLeftFunc(line, func(r rune) bool {
		return r == '!' || unicode.IsSpace(r)
	})

	if strings.HasPrefix(strings.ToLower(trimmed), "title") {
		split := strings.SplitN(trimmed, ":", 2)
		if len(split) != 2 {
			return
		}

		name = strings.TrimSpace(split[1])

		return name, true
	}

	const sub = "abp:subscribe"
	if strings.HasPrefix(strings.ToLower(trimmed), sub) {
		data, err := url.ParseQuery(trimmed[len(sub):])
		if err != nil {
			return
		}
		name = strings.TrimSpace(data.Get("title"))
		if name != "" {
			return name, true
		}
	}

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
