package main

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"path"
	"runtime"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/google/go-github/v43/github"
	"golang.org/x/oauth2"
	"xarantolus/generator/util"
)

func main() {
	// Set one up @ https://github.com/settings/tokens/new
	ghToken := os.Getenv("GITHUB_TOKEN")
	if strings.TrimSpace(ghToken) == "" {
		log.Fatalf("no GITHUB_TOKEN env variable available\n")
	}

	repoOwner := os.Getenv("INITIAL_REPO_OWNER")
	if strings.TrimSpace(repoOwner) == "" {
		repoOwner = "xarantolus"
	}
	repoName := os.Getenv("INITIAL_REPO_NAME")
	if strings.TrimSpace(repoName) == "" {
		repoName = "filtrite"
	}

	log.Printf("[Info] Working on %s/%s\n", repoOwner, repoName)

	ctx := context.Background()
	token := oauth2.NewClient(ctx, oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: ghToken},
	))

	client := github.NewClient(token)

	mainRepo, _, err := client.Repositories.Get(ctx, repoOwner, repoName)
	if err != nil {
		log.Fatalf("loading main repo: %s\n", err.Error())
	}
	_ = mainRepo

	forks, err := util.LoadAllForks(client, repoOwner, repoName)
	if err != nil {
		log.Fatalf("loading forks: %s\n", err.Error())
	}

	forks = append([]*github.Repository{mainRepo}, forks...)

	log.Printf("[Info] Fetched %d forks\n", len(forks))

	var filterLists = make(map[string][]ListFileInfo)

	for _, fork := range forks {
		err := getForkInfo(client, fork, filterLists)
		if err != nil {
			log.Printf("[Warning] Error for %s/%s: %s\n", fork.GetOwner().GetLogin(), fork.GetName(), err.Error())
		}
	}

	// Afterwards bring it into a presentable format

	util.SaveJSON("lists.json", filterLists)
	util.LoadJSON("lists.json", &filterLists)

	var (
		filterListNameMapping     = make(map[string]string)
		filterListNameMappingLock sync.Mutex
	)

	var deduplicatedFilterlists = make(map[string][]ListFileInfo)
	for k, fls := range filterLists {
		deduplicatedFilterlists[k] = deduplicateFilterlists(fls)
	}

	util.SaveJSON("lists-dedup.json", deduplicatedFilterlists)
	util.LoadJSON("lists-dedup.json", &deduplicatedFilterlists)

	var urls = getUniqueURLs(deduplicatedFilterlists)

	parallelize(urls, func(u string) {
		name, err := util.GetFilterListNameFromURL(u)
		if err != nil {
			parsed, err := url.ParseRequestURI(u)
			if err != nil {
				return
			}

			name = stripExtension(parsed.Path)
		}

		filterListNameMappingLock.Lock()
		filterListNameMapping[u] = name
		filterListNameMappingLock.Unlock()
	})

	util.SaveJSON("url-mapping.json", filterListNameMapping)
	util.LoadJSON("url-mapping.json", filterListNameMapping)

	fmt.Println(filterListNameMapping)

	// var output []PresentableListFile

	// for _, info := range filterLists {

	// }

	_ = filterListNameMapping
}

func stripExtension(p string) string {
	name := path.Base(p)
	ext := path.Ext(name)
	if ext != "" {
		name = name[:len(name)-len(ext)]
	}
	return name
}

var correctCasingReplacer = strings.NewReplacer(
	"Ublock", "uBlock",
	"Adblock", "AdBlock",
)

func makeListTitle(name string) (out string) {
	fields := strings.FieldsFunc(name, func(r rune) bool {
		return !unicode.IsLetter(r)
	})

	return correctCasingReplacer.Replace(strings.Title(strings.ToLower(strings.Join(fields, " "))))
}

type ListFileInfo struct {
	Name string `json:"name"`

	FilterFileURL string `json:"filter_file_url"`

	RepoOwner string `json:"repo_owner"`
	RepoName  string `json:"repo_name"`

	ListURL string `json:"list_url"`

	URLs []string `json:"urls"`
}

type PresentableListFile struct {
	DisplayName string       `json:"display_name"`
	Info        ListFileInfo `json:"info"`
}

type URLMapping struct {
	URL   string `json:"url"`
	Title string `json:"title"`
}

var ignoredFileNames = map[string]bool{
	"bromite-default.txt": true,
}

func getUniqueURLs(info map[string][]ListFileInfo) (urls []string) {
	var deduplicate = map[string]bool{}
	for _, filters := range info {
		for _, filter := range filters {
			for _, u := range filter.URLs {
				if deduplicate[u] {
					continue
				}
				deduplicate[u] = true
				urls = append(urls, u)
			}
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

func getAssetByName(assets []*github.ReleaseAsset, fn string) *github.ReleaseAsset {
	for _, a := range assets {
		if a.GetName() == fn {
			return a
		}
	}

	return nil
}

func getLatestURL(a *github.ReleaseAsset, owner, repo string) string {
	// https://github.com/USERNAME/filtrite/releases/latest/download/FILENAME.dat
	return (&url.URL{
		Scheme: "https",
		Host:   "github.com",
		Path:   fmt.Sprintf("%s/%s/releases/latest/download/%s", url.PathEscape(owner), url.PathEscape(repo), url.PathEscape(*a.Name)),
	}).String()
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

func getForkInfo(client *github.Client, fork *github.Repository, filterLists map[string][]ListFileInfo) (err error) {
	// Only look at forks with a compatible license
	if fork.License == nil || !strings.EqualFold(fork.License.GetSPDXID(), "MIT") {
		err = fmt.Errorf("license identifier incompatible, must be MIT")
		return
	}

	ctx := context.Background()

	forkUser, forkRepoName := fork.GetOwner().GetLogin(), fork.GetName()

	fc, listFiles, _, err := client.Repositories.GetContents(ctx, forkUser, forkRepoName, "lists", nil)
	if err != nil {
		return
	}
	if fc != nil {
		err = fmt.Errorf("invalid file \"lists\" instead of directory\n")
		return
	}

	latestRelease, _, err := client.Repositories.GetLatestRelease(ctx, forkUser, forkRepoName)
	if err != nil {
		err = fmt.Errorf("no latest release available: %s\n", err.Error())
		return
	}
	if len(latestRelease.Assets) == 0 {
		err = fmt.Errorf("latest release has no assets\n")
		return
	}

	for _, listFile := range listFiles {
		fn := listFile.GetName()
		// We only want text files, just like filtrite itself
		if listFile.GetType() != "file" || !strings.HasSuffix(fn, ".txt") ||
			listFile.GetSize() == 0 || ignoredFileNames[fn] {
			continue
		}

		datFileName := fn[:len(fn)-4] + ".dat"
		asset := getAssetByName(latestRelease.Assets, datFileName)
		if asset == nil {
			log.Printf("[Warning] Looks like the list %q (%q) in %s/%s is not being released\n", fn, datFileName, forkUser, forkRepoName)
			continue
		}
		// Ignore outdated repos; however still support repos that decided to always use the same release
		// (instead of creating a new release, they just update the current release)
		// Assets should be generated at least every 10 days or so
		if time.Since(asset.GetUpdatedAt().Time) > 10*24*time.Hour {
			continue
		}

		lists, err := util.RequestListURLs(listFile.GetDownloadURL())
		if err != nil {
			log.Printf("requesting list %q in %s/%s: %s\n", fn, forkUser, forkRepoName, err.Error())
			continue
		}
		if len(lists) == 0 {
			continue
		}

		filterLists[fn] = append(filterLists[fn], ListFileInfo{
			Name: makeListTitle(stripExtension(fn)),

			RepoOwner: forkUser,
			RepoName:  forkRepoName,
			URLs:      lists,

			ListURL:       listFile.GetDownloadURL(),
			FilterFileURL: getLatestURL(asset, forkUser, forkRepoName),
		})
	}
	return
}
