package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
)

var gistIdRule = regexp.MustCompile("(.*)/(.*)")

func parse(filename string) (gistIds []string, err error) {
	fp, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		group := gistIdRule.FindSubmatch([]byte(scanner.Text()))
		if len(group) < 3 {
			fmt.Printf("Skipped against an invalid gist url: %v\n", scanner.Text())
			continue
		}

		gistIds = append(gistIds, string(group[2]))
	}
	if err = scanner.Err(); err != nil {
		return nil, err
	}

	return gistIds, nil
}
