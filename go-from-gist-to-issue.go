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
	app.Action = func(c *cli.Context) {
		if len(c.Args()) < 1 {
			fmt.Println("Not Found an argument of a filename")
			return
		}

		gistIds, err := Parse(c.Args()[0])
		if err != nil {
			fmt.Printf("Failed to parse a given file: err=%v\nAborted.\n", err)
			return
		}
		pp.Print(gistIds)
	}
	app.Run(os.Args)
}
