package main

import (
	"os"

	"github.com/codegangsta/cli"
	"github.com/k0kubun/pp"
)

func main() {
	app := cli.NewApp()
	app.Name = "go-from-gist-to-issue"
	app.Usage = "importing gists to github issues"
	app.EnableBashCompletion = true
	app.Action = func(c *cli.Context) {
		pp.Print(c.Args())
	}
	app.Run(os.Args)
}
