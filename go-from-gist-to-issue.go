package main

import (
	"fmt"
	"os"
	"sync"

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

		gistIds, err := Parse(c.String("gist"))
		if err != nil {
			fmt.Printf("Failed to parse a given file: err=%v\nAborted.\n", err)
			return
		}

		if c.Bool("verbose") {
			pp.Println(gistIds)
		}

		importedCount := 0

		if !c.Bool("sequence") {
			var wg sync.WaitGroup

			for _, id := range gistIds {
				wg.Add(1)

				go func(_id string, _c *cli.Context) {
					defer wg.Done()

					err = run(_id, _c)
					if err != nil {
						fmt.Println(err)
						return
					}

					importedCount++

				}(id, c)
			}
			wg.Wait()
		} else {
			for _, id := range gistIds {
				err = run(id, c)
				if err != nil {
					fmt.Println(err)
					continue
				}

				importedCount++
			}
		}

		fmt.Printf("Completed to import gists info: count=%v\n", importedCount)
	}
	app.Run(os.Args)
}

func run(id string, c *cli.Context) (err error) {
	github := CreateGitHub(c.String("token"), c.Bool("verbose"))

	gistInfo, err := github.GetGist(id)
	if err != nil {
		return fmt.Errorf("Failed to get gist info: gistId=%v, err=%v\nSkipped.\n", id, err)
	}

	err = github.ImportGistToIssue(gistInfo, c.String("repo"), c.Bool("dry-run"))
	if err != nil {
		return fmt.Errorf("Failed to import gist info: gistId=%v, err=%v\nSkipped.\n", id, err)
	}

	return nil
}
