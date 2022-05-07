package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/google/go-github/v43/github"
	"xarantolus/generator/util"
)

type OutputInfo struct {
	Date time.Time `json:"date"`

	Lists []PresentableListFile `json:"lists"`
}

type ListFileInfo struct {
	Name string `json:"name"`

	FilterFileURL string `json:"filter_file_url"`

	Stars int `json:"stars"`

	RepoOwner string `json:"repo_owner"`
	RepoName  string `json:"repo_name"`

	ListURL string `json:"list_url"`

	URLs []string `json:"urls"`
}

func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if strings.TrimSpace(val) == "" {
		return fallback
	}
	return val
}

func makePresentable(info ListFileInfo, urlTitles map[string]string) PresentableListFile {
	var urls []URLMapping
	for _, u := range info.URLs {
		ut := strings.TrimSpace(urlTitles[u])
		if ut == "" {
			ut = "Unknown"
		}
		urls = append(urls, URLMapping{
			URL:   u,
			Title: ut,
		})
	}

	sort.Slice(urls, func(i, j int) bool {
		if urls[i].Title == urls[j].Title {
			return urls[i].URL < urls[j].URL
		}
		return urls[i].Title < urls[j].Title
	})

	return PresentableListFile{
		DisplayName:   info.Name,
		URLs:          urls,
		FilterFileURL: info.FilterFileURL,
		Stars:         info.Stars,
		RepoOwner:     info.RepoOwner,
		RepoName:      info.RepoName,
		ListURL:       info.ListURL,
	}
}

type PresentableListFile struct {
	DisplayName string       `json:"display_name"`
	URLs        []URLMapping `json:"urls"`

	FilterFileURL string `json:"filter_file_url"`

	Stars int `json:"stars"`

	RepoOwner string `json:"repo_owner"`
	RepoName  string `json:"repo_name"`

	ListURL string `json:"list_url"`
}

type URLMapping struct {
	URL   string `json:"url"`
	Title string `json:"title"`
}

var ignoredFileNames = map[string]bool{
	"bromite-default.txt": true,
}

func getUniqueURLs(info []ListFileInfo) (urls []string) {
	var deduplicate = map[string]bool{}
	for _, fl := range info {
		for _, u := range fl.URLs {
			if deduplicate[u] {
				continue
			}
			deduplicate[u] = true
			urls = append(urls, u)
		}
	}

	return
}

func deduplicateFilterlists(lists []ListFileInfo) (output []ListFileInfo) {
	for _, list := range lists {
		// Check if the list contains all the same URLs as a list we already have in our output
		var isUnique = true
		for _, ol := range output {
			if len(list.URLs) == len(ol.URLs) && containsAll(list.URLs, ol.URLs) && containsAll(ol.URLs, list.URLs) {
				isUnique = false
				break
			}
		}

		// If unique, we keep it
		if isUnique {
			output = append(output, list)
		}
	}

	return output
}

func parallelize(urls []string, f func(url string)) {
	workChan := make(chan string)
	wg := &sync.WaitGroup{}
	worker := func(wg *sync.WaitGroup, c chan string) {
		defer wg.Done()
		for u := range c {
			f(u)
		}
	}
	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go worker(wg, workChan)
	}

	for _, u := range urls {
		workChan <- u
	}
	close(workChan)

	wg.Wait()
}

// containsAll returns if subset is a subset of set
func containsAll(subset, set []string) bool {
	var asmap = map[string]bool{}
	for _, s := range subset {
		asmap[s] = true
	}

	for _, s := range set {
		if asmap[s] == false {
			return false
		}
	}

	return true
}

func getForkInfo(client *github.Client, fork *github.Repository, filterLists []ListFileInfo) (out []ListFileInfo, err error) {
	out = filterLists
	// Only look at forks with a compatible license
	if fork.License == nil || !strings.EqualFold(fork.License.GetSPDXID(), "MIT") {
		err = fmt.Errorf("license identifier incompatible, must be MIT")
		return
	}

	forkUser, forkRepoName := fork.GetOwner().GetLogin(), fork.GetName()

	// Forks need to be at least 2 days old before we list them, sometimes people fork and delete their repo within a day
	if time.Since(fork.GetCreatedAt().Time) < 2*24*time.Hour {
		log.Printf("[Warning] Ignoring fork %s/%s because it's too young", forkUser, forkRepoName)
		return
	}

	ctx := context.Background()
	fc, listFiles, _, err := client.Repositories.GetContents(ctx, forkUser, forkRepoName, "lists", nil)
	if err != nil {
		return
	}
	if fc != nil {
		err = fmt.Errorf("invalid file \"lists\" instead of directory")
		return
	}

	latestRelease, _, err := client.Repositories.GetLatestRelease(ctx, forkUser, forkRepoName)
	if err != nil {
		err = fmt.Errorf("no latest release available: %s", err.Error())
		return
	}
	if len(latestRelease.Assets) == 0 {
		err = fmt.Errorf("latest release has no assets")
		return
	}

	// https://github.com/bromite/bromite/blob/4f10d11318703835bb201a54d606e2b8b2dd896b/build/patches/Bromite-AdBlockUpdaterService.patch#L1131
	const bromiteMaxFilterSize = 1024 * 1024 * 10

	for _, listFile := range listFiles {
		fn := listFile.GetName()
		// We only want text files, just like filtrite itself
		if listFile.GetType() != "file" || !strings.HasSuffix(fn, ".txt") ||
			listFile.GetSize() == 0 || ignoredFileNames[fn] {
			continue
		}

		datFileName := fn[:len(fn)-4] + ".dat"
		asset := util.GetAssetByName(latestRelease.Assets, datFileName)
		if asset == nil {
			log.Printf("[Warning] Looks like the list %q (%q) in %s/%s is not being released", fn, datFileName, forkUser, forkRepoName)
			continue
		}
		// Ignore outdated repos; however still support repos that decided to always use the same release
		// (instead of creating a new release, they just update the current release)
		// Assets should be generated at least every 10 days or so
		if time.Since(asset.GetUpdatedAt().Time) > 10*24*time.Hour {
			log.Printf("[Warning] Ignoring outdated asset %q in %s/%s", datFileName, forkUser, forkRepoName)
			continue
		}
		if asset.GetSize() > bromiteMaxFilterSize {
			log.Printf("[Warning] Ignoring asset %q in %s/%s because it's too large", datFileName, forkUser, forkRepoName)
			continue
		}

		lists, err := util.RequestListURLs(listFile.GetDownloadURL())
		if err != nil {
			log.Printf("[Error] Requesting list %q in %s/%s: %s", fn, forkUser, forkRepoName, err.Error())
			continue
		}
		if len(lists) == 0 {
			log.Printf("[Warning] List %q in %s/%s doesn't define any filterlists we could download", fn, forkUser, forkRepoName)
			continue
		}

		filterLists = append(filterLists, ListFileInfo{
			Name: util.MakeListTitle(util.StripExtension(fn)),

			RepoOwner: forkUser,
			RepoName:  forkRepoName,
			URLs:      lists,

			ListURL:       listFile.GetDownloadURL(),
			FilterFileURL: util.GetLatestURL(asset, forkUser, forkRepoName),

			Stars: fork.GetStargazersCount(),
		})
	}

	return filterLists, nil
}

const (
	bromiteDefaultOrg      = "bromite"
	bromiteDefaultRepo     = "filters"
	bromiteDefaultListFile = "https://www.bromite.org/filters/lists.txt"
)

func fetchBromiteDefaultList(client *github.Client) (l ListFileInfo, err error) {
	ctx := context.Background()

	repo, _, err := client.Repositories.Get(ctx, bromiteDefaultOrg, bromiteDefaultRepo)
	if err != nil {
		return
	}

	lists, err := util.RequestListURLs(bromiteDefaultListFile)
	if err != nil {
		return
	}

	return ListFileInfo{
		Name:          "Bromite Default",
		FilterFileURL: "https://www.bromite.org/filters/filters.dat",
		Stars:         repo.GetStargazersCount(),
		RepoOwner:     bromiteDefaultOrg,
		RepoName:      bromiteDefaultRepo,
		ListURL:       bromiteDefaultListFile,
		URLs:          lists,
	}, nil
}
