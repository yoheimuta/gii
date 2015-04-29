package main

import (
	"fmt"

	"github.com/google/go-github/github"
	//"github.com/k0kubun/pp"
)

type GitHub struct {
	client *github.Client
}

type GistInfo struct {
	gist     *github.Gist
	comments []github.GistComment
}

func CreateGitHub() *GitHub {
	g := &GitHub{
		client: github.NewClient(nil),
	}
	return g
}

func (g *GitHub) GetIssues(gistIds []string) (gistInfo []GistInfo, err error) {
	for _, id := range gistIds {
		fmt.Printf("Downloading a gist and comments: %v\n", id)

		gist, res, err := g.client.Gists.Get(id)
		if err != nil {
			fmt.Printf("Skipped against a failed Gists.Get API: %v\n", err)
			continue
		}
		if res.StatusCode != 200 {
			fmt.Printf("Skipped against a invalid response of Gists.Get API: %v\n", res.StatusCode)
			continue
		}
		//pp.Print(gist)

		comments, res1, err1 := g.client.Gists.ListComments(id, &github.ListOptions{})
		if err1 != nil {
			fmt.Printf("Skipped against a failed Gists.ListComments API: %v\n", err1)
			continue
		}
		if res1.StatusCode != 200 {
			fmt.Printf("Skipped against a invalid response of Gists.ListComments API: %v\n", res1.StatusCode)
			continue
		}
		//pp.Print(comments)

		gistInfo = append(gistInfo, GistInfo{
			gist:     gist,
			comments: comments,
		})

		fmt.Printf("Downloaded a gist and comments: %v\n", id)
	}

	return gistInfo, nil
}
