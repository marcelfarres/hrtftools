package command

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/marcelfarres/hrtftools/browser"
	"github.com/marcelfarres/hrtftools/parser"
)

func CmdBestMatch(c *cli.Context) {
	var (
		dbFile, methode, measureName, measureValueS string
		keys                                        []string
		resNum                                      int
		measureValue                                []float64
		measurements                                map[string]parser.Measurement
		r                                           map[string]float64
	)

	resNum = c.Int("num-results")
	dbFile = c.String("db-file")
	methode = c.String("methode")
	measureName = c.String("active-measurements-n")
	measureValueS = c.String("active-measurements-v")

	if len(measureName) == 0 || len(measureValueS) == 0 {
		fmt.Println("No measurements keys or values where provided, please enter subject data.")
		return
	}

	keys = strings.Split(measureName, ",")
	tmp := strings.Split(measureValueS, ",")
	for _, sF := range tmp {
		f, err := strconv.ParseFloat(sF, 64)
		if err != nil {
			fmt.Println("Unable to parse measurements values to float")
			return
		}
		measureValue = append(measureValue, f)
	}

	if len(keys) != len(measureValue) {
		fmt.Println("Measurements keys and values are different, did you miss some?")
		return
	}

	dbTxt, err := os.Open(dbFile)
	if err != nil {
		fmt.Println("opening db file", err.Error())
		return
	}

	jsonParser := json.NewDecoder(dbTxt)
	if err = jsonParser.Decode(&measurements); err != nil {
		fmt.Println("not able to parse json", err.Error())
		return
	}

	switch methode {
	case "min-dist":
		r = browser.MinDist(measurements, resNum, keys, measureValue)
	default:
		fmt.Printf("Method %s not defined", methode)
		return
	}

	b, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("Best subj id are:\n %s", string(b))
}
