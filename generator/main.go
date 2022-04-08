package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"strings"

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

	// forks, err := util.LoadAllForks(client, repoOwner, repoName)
	// if err != nil {
	// 	log.Fatalf("loading forks: %s\n", err.Error())
	// }

	// forks = append([]*github.Repository{mainRepo}, forks...)

	// log.Printf("[Info] Fetched %d forks\n", len(forks))

	var filterLists = make(map[string][]ListFileInfo)

	// for _, fork := range forks {
	// 	// Only look at forks with a compatible license
	// 	if fork.License == nil || !strings.EqualFold(fork.License.GetSPDXID(), "MIT") {
	// 		continue
	// 	}

	// 	forkUser, forkRepoName := fork.GetOwner().GetLogin(), fork.GetName()

	// 	fc, listFiles, _, err := client.Repositories.GetContents(ctx, forkUser, forkRepoName, "lists", nil)
	// 	if err != nil {
	// 		log.Printf("[Warning] Error listing %s/%s: %s\n", forkUser, forkRepoName, err.Error())
	// 		continue
	// 	}
	// 	if fc != nil {
	// 		log.Printf("[Warning] Invalid file \"lists\" instead of directory in %s/%s\n", forkUser, forkRepoName)
	// 		continue
	// 	}

	// 	latestRelease, _, err := client.Repositories.GetLatestRelease(ctx, forkUser, forkRepoName)
	// 	if err != nil {
	// 		log.Printf("[Warning] No latest release available in %s/%s: %s\n", forkUser, forkRepoName, err.Error())
	// 		continue
	// 	}
	// 	if len(latestRelease.Assets) == 0 {
	// 		continue
	// 	}

	// 	for _, listFile := range listFiles {
	// 		fn := listFile.GetName()
	// 		// We only want text files, just like filtrite itself
	// 		if listFile.GetType() != "file" || !strings.HasSuffix(fn, ".txt") ||
	// 			listFile.GetSize() == 0 || ignoredFileNames[fn] {
	// 			continue
	// 		}

	// 		datFileName := fn[:len(fn)-4] + ".dat"
	// 		asset := getAssetByName(latestRelease.Assets, datFileName)
	// 		if asset == nil {
	// 			log.Printf("[Warning] Looks like the list %q (%q) in %s/%s is not being released\n", fn, datFileName, forkUser, forkRepoName)
	// 			continue
	// 		}
	// 		// Ignore outdated repos; however still support repos that decided to always use the same release (instead of creating a new release, they just update the current release)
	// 		if time.Since(asset.GetUpdatedAt().Time) > 7*24*time.Hour {
	// 			continue
	// 		}

	// 		lists, err := util.RequestListURLs(listFile.GetDownloadURL())
	// 		if err != nil {
	// 			log.Printf("requesting list %q in %s/%s: %s\n", fn, forkUser, forkRepoName, err.Error())
	// 			continue
	// 		}
	// 		if len(lists) == 0 {
	// 			continue
	// 		}

	// 		filterLists[fn] = append(filterLists[fn], ListFileInfo{
	// 			Name: fn,

	// 			RepoOwner: forkUser,
	// 			RepoName:  forkRepoName,
	// 			URLs:      lists,

	// 			ListURL:       listFile.GetDownloadURL(),
	// 			FilterFileURL: getLatestURL(asset, forkUser, forkRepoName),
	// 		})
	// 	}
	// }

	// // Afterwards bring it into a presentable format

	// data, err := json.Marshal(filterLists)
	// if err != nil {
	// 	panic(err)
	// }

	// ioutil.WriteFile("lists.json", data, 0o664)

	c, err := ioutil.ReadFile("lists.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(c, &filterLists)
	if err != nil {
		panic(err)
	}

	var filterListNameMapping = make(map[string]string)

	for _, filters := range filterLists {
		for _, filter := range filters {
			for _, u := range filter.URLs {
				if _, ok := filterListNameMapping[u]; ok {
					continue
				}

				name, err := util.GetFilterListNameFromURL(u)
				if err != nil {
					log.Printf("[Warning] Could not get filter list title for %q", u)
					continue
				}

				filterListNameMapping[u] = name
			}
		}
	}

	fmt.Println(filterListNameMapping)

	_ = filterListNameMapping
}

type ListFileInfo struct {
	Name string

	FilterFileURL string

	RepoOwner, RepoName string

	ListURL string

	URLs []string
}

type PresentableListFile struct {
	DisplayName string

	Info ListFileInfo
}

var ignoredFileNames = map[string]bool{
	"bromite-default.txt": true,
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
