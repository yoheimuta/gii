package main

import (
	"fmt"

	"golang.org/x/oauth2"

	"github.com/google/go-github/github"
	"github.com/k0kubun/pp"
)

type GitHub struct {
	client        *github.Client
	isVerbose     bool
	dryRun        bool
	noGistComment bool
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

func CreateGitHub(token string, isVerbose bool, dryRun bool, noGistComment bool) *GitHub {
	ts := &tokenSource{
		&oauth2.Token{AccessToken: token},
	}

	tc := oauth2.NewClient(oauth2.NoContext, ts)

	g := &GitHub{
		client:        github.NewClient(tc),
		isVerbose:     isVerbose,
		dryRun:        dryRun,
		noGistComment: noGistComment,
	}
	return g
}

func (g *GitHub) Run(gistId string, repo string) (err error) {
	gistInfo, err := g.getGist(gistId)
	if err != nil {
		return fmt.Errorf("Failed to get a gist info: gistId=%v, err=%v Skipped.\n", gistId, err)
	}

	issueURL, err := g.importGistToIssue(gistInfo, repo)
	if err != nil {
		return fmt.Errorf("Failed to import a gist info: gistId=%v, err=%v Skipped.\n", gistId, err)
	}

	if !g.noGistComment {
		if !g.dryRun {
			body := fmt.Sprintf("Automatically imported to %v.", issueURL)
			err = g.createGistComment(gistId, &body)
			if err != nil {
				return fmt.Errorf("Failed to create a gist comment: gistId=%v, err=%v Skipped.\n", gistId, err)
			}
		} else {
			fmt.Printf("Dry-run to create a gist comment: gistId=%v\n", gistId)
		}
	}

	return nil
}

func (g *GitHub) getGist(gistId string) (gistInfo GistInfo, err error) {
	fmt.Printf("Downloading a gist and comments: %v\n", gistId)

	gist, res, err := g.client.Gists.Get(gistId)
	if err != nil {
		err = fmt.Errorf("Skipped against a failed Gists.Get API: %v\n", err)
		return
	}
	if res.StatusCode != 200 {
		err = fmt.Errorf("Skipped against an invalid response of Gists.Get API: %v\n", res.StatusCode)
		return
	}
	if g.isVerbose {
		pp.Println(gist)
	}

	comments, res, err := g.client.Gists.ListComments(gistId, &github.ListOptions{})
	if err != nil {
		err = fmt.Errorf("Skipped against a failed Gists.ListComments API: %v\n", err)
		return
	}
	if res.StatusCode != 200 {
		err = fmt.Errorf("Skipped against a invalid response of Gists.ListComments API: %v\n", res.StatusCode)
		return
	}
	if g.isVerbose {
		pp.Println(comments)
	}

	gistInfo = GistInfo{
		gist:     gist,
		comments: comments,
	}

	fmt.Printf("Downloaded  a gist and comments: %v\n", gistId)

	return gistInfo, nil
}

func (g *GitHub) importGistToIssue(gistInfo GistInfo, repo string) (issueURL string, err error) {
	gist, err := g.extractGist(gistInfo.gist)
	if err != nil {
		return "", fmt.Errorf("Failed to extract gist: %v. Skipped\n", err)
	}

	var issue *github.Issue

	if !g.dryRun {
		owner := *gist["Owner"]
		body := fmt.Sprintf("Automatically imported from %v.\n\n%v", *gist["URL"], *gist["Content"])

		var res *github.Response
		issue, res, err = g.client.Issues.Create(owner, repo, &github.IssueRequest{
			Title: gist["Title"],
			Body:  &body,
		})

		if err != nil {
			return "", fmt.Errorf("Skipped against a failed Issues.Create API: %v\n", err)
		}
		if res.StatusCode != 201 {
			return "", fmt.Errorf("Skipped against an invalid response of Issues.Create API: %v\n", res.StatusCode)
		}

		fmt.Printf("Created an issue: from %v to %v\n", *gist["URL"], *issue.HTMLURL)
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

	for _, gistComment := range gistInfo.comments {
		comment, err := g.extractGistComment(gistComment)
		if err != nil {
			return *issue.HTMLURL, fmt.Errorf("Skipped. Failed to extract gistComment: %v\n", err)
		}
		if !g.dryRun {
			number := *issue.Number
			commentOwner := *comment["Owner"]

			issueComment, res, err := g.client.Issues.CreateComment(commentOwner, repo, number, &github.IssueComment{
				Body: comment["Body"],
			})

			if err != nil {
				return *issue.HTMLURL, fmt.Errorf("Skipped against a failed Issues.CreateComment API: %v\n", err)
			}
			if res.StatusCode != 201 {
				return *issue.HTMLURL, fmt.Errorf("Skipped against an invalid response of Issues.CreateComment API: %v\n", res.StatusCode)
			}

			fmt.Printf("Created a comment: %v\n", *issueComment.HTMLURL)
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
	if issue != nil {
		return *issue.HTMLURL, nil
	} else {
		return "", nil
	}
}

func (g *GitHub) extractGist(gist *github.Gist) (extracted map[string]*string, err error) {
	var gistContent *string
	//var gistTitle *string
	//for title, file := range gist.Files {
	for _, file := range gist.Files {
		//_title := string(title)
		//gistTitle = &_title
		gistContent = file.Content

		// TODO:
		break
	}
	gistTitle := gist.Description
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

func (g *GitHub) createGistComment(gistId string, body *string) (err error) {
	comment, res, err := g.client.Gists.CreateComment(gistId, &github.GistComment{
		Body: body,
	})

	if err != nil {
		return fmt.Errorf("Skipped against a failed Gists.CreateComment API: %v\n", err)
	}
	if res.StatusCode != 201 {
		return fmt.Errorf("Skipped against an invalid response of Gists.CreateComment API: %v\n", res.StatusCode)
	}

	fmt.Printf("Created a gist comment: %v\n", *comment.URL)
	if g.isVerbose {
		//TODO: pp panicked
		//pp.Println(comment)
		//pp.Println(res)
	}

	return nil
}
