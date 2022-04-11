package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/url"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/google/go-github/v43/github"
	"golang.org/x/oauth2"
	"xarantolus/generator/util"
)

func main() {
	// Set one up @ https://github.com/settings/tokens/new
	var (
		ghToken    = os.Getenv("GITHUB_TOKEN")
		repoOwner  = getEnv("INITIAL_REPO_OWNER", "xarantolus")
		repoName   = getEnv("INITIAL_REPO_NAME", "filtrite")
		outputFile = getEnv("OUTPUT_FILE", "filterlists_jsonp.js")
	)
	if strings.TrimSpace(ghToken) == "" {
		log.Fatalf("no GITHUB_TOKEN env variable available\n")
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

	forks, err := util.LoadAllForks(client, repoOwner, repoName)
	if err != nil {
		log.Fatalf("loading forks: %s\n", err.Error())
	}

	// We sort forks by an arbitrary measure of importance
	sort.Slice(forks, func(i, j int) bool {
		return forks[i].GetStargazersCount() > forks[j].GetStargazersCount()
	})
	// Add the main repo at the beginning, regardless of count
	forks = append([]*github.Repository{mainRepo}, forks...)

	log.Printf("[Info] Fetched %d forks\n", len(forks))

	var filterLists []ListFileInfo

	for _, fork := range forks {
		filterLists, err = getForkInfo(client, fork, filterLists)
		if err != nil {
			log.Printf("[Warning] Error for %s/%s: %s\n", fork.GetOwner().GetLogin(), fork.GetName(), err.Error())
		}
	}

	// Afterwards bring it into a presentable format

	var (
		filterListUrlNameMapping  = make(map[string]string)
		filterListNameMappingLock sync.Mutex
	)

	var deduplicatedFilterlists = deduplicateFilterlists(filterLists)

	var urls = getUniqueURLs(deduplicatedFilterlists)

	parallelize(urls, func(u string) {
		name, err := util.GetFilterListNameFromURL(u)
		if err != nil {
			parsed, err := url.ParseRequestURI(u)
			if err != nil {
				return
			}

			name = util.MakeListTitle(util.StripExtension(parsed.Path))
		}

		filterListNameMappingLock.Lock()
		filterListUrlNameMapping[u] = name
		filterListNameMappingLock.Unlock()
	})

	var outputFilterLists []PresentableListFile

	for _, info := range deduplicatedFilterlists {
		outputFilterLists = append(outputFilterLists, makePresentable(info, filterListUrlNameMapping))
	}

	if len(outputFilterLists) == 0 {
		panic("No output produced, something is wrong with this program")
	}

	var buf bytes.Buffer

	buf.WriteString("jsonp(")
	b, err := json.Marshal(OutputInfo{Date: time.Now().UTC(), Lists: outputFilterLists})
	if err != nil {
		panic(err)
	}
	buf.Write(b)
	buf.WriteString(")")

	err = os.WriteFile(outputFile, buf.Bytes(), 0o666)
	if err != nil {
		log.Fatalf("writing result file: %s\n", err.Error())
	}
}
