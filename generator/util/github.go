package util

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/google/go-github/v43/github"
)

func LoadAllForks(client *github.Client, owner, name string) (forks []*github.Repository, err error) {

	var repoPage = 1

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	for {
		repos, resp, ferr := client.Repositories.ListForks(ctx, owner, name, &github.RepositoryListForksOptions{
			ListOptions: github.ListOptions{
				PerPage: 100,
				Page:    repoPage,
			},
		})
		if ferr != nil {
			err = fmt.Errorf("loading fork page %d: %s", repoPage, ferr.Error())
			return
		}

		forks = append(forks, repos...)

		repoPage = resp.NextPage
		if repoPage == 0 {
			break
		}
	}

	return forks, nil
}

func GetAssetByName(assets []*github.ReleaseAsset, fn string) *github.ReleaseAsset {
	for _, a := range assets {
		if a.GetName() == fn {
			return a
		}
	}

	return nil
}

func GetLatestURL(a *github.ReleaseAsset, owner, repo string) string {
	// https://github.com/USERNAME/filtrite/releases/latest/download/FILENAME.dat
	return (&url.URL{
		Scheme: "https",
		Host:   "github.com",
		Path:   fmt.Sprintf("%s/%s/releases/latest/download/%s", url.PathEscape(owner), url.PathEscape(repo), url.PathEscape(*a.Name)),
	}).String()
}
