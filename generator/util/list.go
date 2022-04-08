package util

import (
	"bufio"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

// ReadList returns all URLs the given reader without duplicates, sorted
func ReadList(r io.Reader) (entries []string, err error) {
	var entriesMap = map[string]bool{}

	scan := bufio.NewScanner(r)

	for scan.Scan() {
		t := strings.TrimSpace(scan.Text())

		// Remove comments and empty lines
		if strings.HasPrefix(t, "#") || t == "" {
			continue
		}

		// Remove invalid URLs
		if _, uerr := url.ParseRequestURI(t); uerr != nil {
			continue
		}

		if !entriesMap[t] {
			// If we haven't seen this URL before, we add it to the list
			entriesMap[t] = true
		}
	}

	for url := range entriesMap {
		entries = append(entries, url)
	}

	sort.Strings(entries)

	return
}

var client = &http.Client{
	Timeout: 30 * time.Second,
}

// RequestListURLs requests and parses list from the given URL
func RequestListURLs(url string) (s []string, err error) {
	resp, err := req(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	return ReadList(resp.Body)
}

func req(url string) (h *http.Response, err error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return
	}
	req.Header.Add("User-Agent", "xarantolus/filtrite-list")

	resp, err := client.Do(req)
	if err != nil {
		return
	}

	return resp, err
}
