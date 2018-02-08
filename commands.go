package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/marcelfarres/hrtftools/command"
)

var GlobalFlags = []cli.Flag{}

var Commands = []cli.Command{
	{
		Name:   "parse",
		Usage:  "Parse HRTF from url.",
		Action: command.CmdParse,
		Flags: []cli.Flag{
			cli.StringFlag{
				EnvVar: "ENV_URL",
				Name:   "url",
				Value:  "http://recherche.ircam.fr/equipes/salles/listen/infomorph_display.php?subject=",
				Usage:  "URL where to fetch DB.",
			},
			cli.IntFlag{
				EnvVar: "ENV_NUM-SUB",
				Name:   "num-subjects, s-num",
				Value:  60,
				Usage:  "Total number of subjects in the DB.",
			}},
	},
	{
		Name:   "best-match",
		Usage:  "Get best match",
		Action: command.CmdBestMatch,
		Flags: []cli.Flag{
			cli.StringFlag{
				EnvVar: "ENV_DB-FILE",
				Name:   "db-file, f",
				Value:  "output.json",
				Usage:  "File to extract the hrtf morphological data.",
			},
			cli.StringFlag{
				EnvVar: "ENV_METHODE",
				Name:   "methode, m",
				Value:  "min-dist",
				Usage:  "Selection methode.",
			},
			cli.StringFlag{
				EnvVar: "ENV_ACTIVE-MEASUREMENTS-N",
				Name:   "active-measurements-n, am-n",
				Value:  "",
				Usage:  "Which measurements you want to input, comma separeted.",
			},
			cli.StringFlag{
				EnvVar: "ENV_ACTIVE-MEASUREMENTS-",
				Name:   "active-measurements-v, am-v",
				Value:  "",
				Usage:  "Value of measurements you want to input, comma separeted.",
			},
			cli.IntFlag{
				EnvVar: "ENV_NUM-RESULTS",
				Name:   "num-results, r",
				Value:  5,
				Usage:  "Total Results to show.",
			}},
	},
}

func CommandNotFound(c *cli.Context, command string) {
	fmt.Fprintf(os.Stderr, "%s: '%s' is not a %s command. See '%s --help'.", c.App.Name, command, c.App.Name, c.App.Name)
	os.Exit(2)
}
