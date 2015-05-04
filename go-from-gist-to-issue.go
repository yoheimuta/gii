package main

import (
	"fmt"
	"os"
	"sync"
	"sync/atomic"

	"github.com/codegangsta/cli"
	"github.com/k0kubun/pp"
)

func main() {
	app := cli.NewApp()
	app.Name = "go-from-gist-to-issue"
	app.Usage = "importing gists to github issues"
	app.EnableBashCompletion = true
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "gist",
			Value: "",
			Usage: "a text file to list up gist url",
		},
		cli.StringFlag{
			Name:  "repo",
			Value: "",
			Usage: "a repository name to be imported from gists",
		},
		cli.StringFlag{
			Name:  "token",
			Value: "",
			Usage: "a github personal access token",
		},
		cli.BoolFlag{
			Name:  "dry-run",
			Usage: "a flag to run without any changes",
		},
		cli.BoolFlag{
			Name:  "verbose",
			Usage: "a flag to log verbosely",
		},
		cli.BoolFlag{
			Name:  "sequence",
			Usage: "a flag to import sequentially",
		},
	}
	app.Action = func(c *cli.Context) {
		if c.String("gist") == "" {
			fmt.Println("Not Found an argument of a gist")
			return
		}
		if c.String("repo") == "" {
			fmt.Println("Not Found an argument of a repo")
			return
		}
		if c.String("token") == "" {
			fmt.Println("Not Found an argument of a token")
			return
		}

		importedCount, err := action(
			c.String("gist"), c.String("repo"), c.String("token"),
			c.Bool("verbose"), c.Bool("sequence"), c.Bool("dry-run"),
		)
		if err != nil {
			fmt.Printf("Failed to import from gists to issues: %v\n", err)
		} else {
			fmt.Printf("Completed to import from gists to issues: count=%v\n", importedCount)
		}
	}
	app.Run(os.Args)
}

func action(gist string, repo string, token string, isVerbose bool, isSequence bool, dryRun bool) (count uint64, err error) {
	var gistIds []string
	gistIds, err = parse(gist)
	if err != nil {
		return 0, fmt.Errorf("Failed to parse a given file: err=%v\nAborted.\n", err)
	}
	if isVerbose {
		pp.Println(gistIds)
	}

	github := CreateGitHub(token, isVerbose)

	if !isSequence {
		var wg sync.WaitGroup

		for _, id := range gistIds {
			wg.Add(1)

			go func(_id string, _repo string, _dryRun bool) {
				defer wg.Done()

				err = github.Run(_id, _repo, _dryRun)
				if err != nil {
					fmt.Println(err)
					return
				}

				atomic.AddUint64(&count, 1)

			}(id, repo, dryRun)
		}
		wg.Wait()
	} else {
		for _, id := range gistIds {
			err = github.Run(id, repo, dryRun)
			if err != nil {
				fmt.Println(err)
				continue
			}

			count++
		}
	}
	return count, nil
}
