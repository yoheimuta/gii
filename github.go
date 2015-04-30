package main

import (
	"fmt"

	"golang.org/x/oauth2"

	"github.com/google/go-github/github"
	"github.com/k0kubun/pp"
)

type GitHub struct {
	client    *github.Client
	isVerbose bool
}

type GistInfo struct {
	gist     *github.Gist
	comments []github.GistComment
}

// tokenSource is an oauth2.TokenSource which returns a static access token
type tokenSource struct {
	token *oauth2.Token
}

// Token implements the oauth2.TokenSource interface
func (t *tokenSource) Token() (*oauth2.Token, error) {
	return t.token, nil
}

func CreateGitHub(token string, isVerbose bool) *GitHub {
	ts := &tokenSource{
		&oauth2.Token{AccessToken: token},
	}

	tc := oauth2.NewClient(oauth2.NoContext, ts)

	g := &GitHub{
		client:    github.NewClient(tc),
		isVerbose: isVerbose,
	}
	return g
}

func (g *GitHub) GetGists(gistIds []string) (gistsInfo []GistInfo, err error) {
	for _, id := range gistIds {
		fmt.Printf("Downloading a gist and comments: %v\n", id)

		gist, res, err := g.client.Gists.Get(id)
		if err != nil {
			fmt.Printf("Skipped against a failed Gists.Get API: %v\n", err)
			continue
		}
		if res.StatusCode != 200 {
			fmt.Printf("Skipped against an invalid response of Gists.Get API: %v\n", res.StatusCode)
			continue
		}
		if g.isVerbose {
			pp.Println(gist)
		}

		comments, res1, err1 := g.client.Gists.ListComments(id, &github.ListOptions{})
		if err1 != nil {
			fmt.Printf("Skipped against a failed Gists.ListComments API: %v\n", err1)
			continue
		}
		if res1.StatusCode != 200 {
			fmt.Printf("Skipped against a invalid response of Gists.ListComments API: %v\n", res1.StatusCode)
			continue
		}
		if g.isVerbose {
			pp.Println(comments)
		}

		gistsInfo = append(gistsInfo, GistInfo{
			gist:     gist,
			comments: comments,
		})

		fmt.Printf("Downloaded  a gist and comments: %v\n", id)
	}

	return gistsInfo, nil
}

func (g *GitHub) ImportGistsToIssues(gistsInfo []GistInfo, repo string, dry_run bool) (processed_count int, err error) {
	for _, gistInfo := range gistsInfo {
		var gist map[string]*string

		gist, err = g.extractGist(gistInfo.gist)
		if err != nil {
			fmt.Printf("Failed to extract gist: %v. Skipped\n", err)
			continue
		}

		var issue *github.Issue

		if !dry_run {
			owner := *gist["Owner"]
			body := fmt.Sprintf("Automatically imported from %s.\n\n%s", *gist["URL"], *gist["Content"])

			var res *github.Response
			issue, res, err = g.client.Issues.Create(owner, repo, &github.IssueRequest{
				Title: gist["Title"],
				Body:  &body,
			})

			if err != nil {
				fmt.Printf("Skipped against a failed Issues.Create API: %v\n", err)
				continue
			}
			if res.StatusCode != 201 {
				fmt.Printf("Skipped against an invalid response of Issues.Create API: %v\n", res.StatusCode)
				continue
			}

			fmt.Printf("Created an issue: %v\n", *issue.Number)
			if g.isVerbose {
				//TODO: pp panicked
				//pp.Println(*issue)
				//pp.Println(*res)
			}
		} else {
			fmt.Println("Dry-run to create an issue")
			if g.isVerbose {
				pp.Println(gist)
			}
		}

		var comment map[string]*string
		for _, gistComment := range gistInfo.comments {
			comment, err = g.extractGistComment(gistComment)
			if err != nil {
				fmt.Printf("Skipped. Failed to extract gistComment: %v\n", err)
				continue
			}
			if !dry_run {
				number := *issue.Number
				commentOwner := *comment["Owner"]

				var issueComment *github.IssueComment
				var res *github.Response
				issueComment, res, err = g.client.Issues.CreateComment(commentOwner, repo, number, &github.IssueComment{
					Body: comment["Body"],
				})

				if err != nil {
					fmt.Printf("Skipped against a failed Issues.CreateComment API: %v\n", err)
					continue
				}
				if res.StatusCode != 201 {
					fmt.Printf("Skipped against an invalid response of Issues.CreateComment API: %v\n", res.StatusCode)
					continue
				}

				fmt.Printf("Created an comment: %v\n", *issueComment.ID)
				if g.isVerbose {
					//TODO: pp panicked
					//pp.Println(issueComment)
					//pp.Println(res)
				}
			} else {
				fmt.Println("Dry-run to create a comment")
				if g.isVerbose {
					pp.Println(comment)
				}
			}
		}
		processed_count++
	}
	return processed_count, nil
}

func (g *GitHub) extractGist(gist *github.Gist) (extracted map[string]*string, err error) {
	var gistTitle *string
	var gistContent *string
	for title, file := range gist.Files {
		gistTitle0 := string(title)
		gistTitle = &gistTitle0
		gistContent = file.Content
	}
	gistOwner := gist.Owner.Login
	gistHTMLURL := gist.HTMLURL

	if *gistTitle == "" {
		return nil, fmt.Errorf("Not Found gistTitle")
	}
	if *gistContent == "" {
		return nil, fmt.Errorf("Not Found gistContent")
	}
	if *gistOwner == "" {
		return nil, fmt.Errorf("Not Found gistOwner")
	}
	if *gistHTMLURL == "" {
		return nil, fmt.Errorf("Not Found gistHTMLURL")
	}

	extracted = map[string]*string{
		"Title":   gistTitle,
		"Content": gistContent,
		"Owner":   gistOwner,
		"URL":     gistHTMLURL,
	}

	return extracted, nil
}

func (g *GitHub) extractGistComment(comment github.GistComment) (extracted map[string]*string, err error) {
	commentBody := comment.Body
	commentOwner := comment.User.Login

	if *commentBody == "" {
		return nil, fmt.Errorf("Not Found commentBody")
	}
	if *commentOwner == "" {
		return nil, fmt.Errorf("Not Found commentOwner")
	}

	extracted = map[string]*string{
		"Body":  commentBody,
		"Owner": commentOwner,
	}

	return extracted, nil
}
