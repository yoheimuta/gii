package main

import (
	"fmt"
	"os"

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
		isVerbose := c.Bool("verbose")

		gistIds, err := Parse(c.String("gist"))
		if err != nil {
			fmt.Printf("Failed to parse a given file: err=%v\nAborted.\n", err)
			return
		}

		if isVerbose {
			pp.Println(gistIds)
		}

		github := CreateGitHub(c.String("token"), isVerbose)
		gistsInfo, err1 := github.GetGists(gistIds)
		if err1 != nil {
			fmt.Printf("Failed to get gists info: err=%v\nAborted.\n", err1)
			return
		}

		imported_count, err2 := github.ImportGistsToIssues(gistsInfo, c.String("repo"), c.Bool("dry-run"))
		if err2 != nil {
			fmt.Printf("Failed to import gists info: err=%v\nAborted.\n", err2)
			return
		}
		fmt.Printf("Completed to import gists info: count=%v\n", imported_count)
	}
	app.Run(os.Args)
}
