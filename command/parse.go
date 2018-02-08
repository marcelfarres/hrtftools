package command

import (
	"fmt"
	"io/ioutil"

	"github.com/codegangsta/cli"
	"github.com/marcelfarres/hrtftools/parser"
)

func CmdParse(c *cli.Context) {

	totalSub := c.Int("num-subjects")
	baseURL := c.String("url")

	b := parser.Parse(baseURL, totalSub)

	err := ioutil.WriteFile("output.json", b, 0644)
	if err != nil {
		fmt.Printf("error writting JSON: %v", err)
	}
}
