package parser

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Measurement struct {
	SubjID  string             `json:"subj_id"`
	DataTxt map[string]string  `json:"data_txt"`
	DataNum map[string]float64 `json:"data_num"`
}

func Parse(baseURL string, totalSub int) (b []byte) {
	var (
		measurements map[string]Measurement
		doc          *goquery.Document
		err          error
	)

	measurements = make(map[string]Measurement)

	for i := 2; i < totalSub; i++ {
		subID := fmt.Sprintf("IRC_1%03d", i)
		url := fmt.Sprintf("%s%s", baseURL, subID)
		fmt.Println(subID, url)
		_, err := http.Get(url)
		if err != nil {
			fmt.Printf("URL %s does not exist, error:%s\n", url, err)
			continue
		}
		doc, err = goquery.NewDocument(url)
		if err != nil {
			fmt.Printf("URL %s can't not be parsed, error:%s\n", url, err)
			continue
		}

		dTxt := make(map[string]string)
		dNum := make(map[string]float64)

		doc.Find("ul").Each(func(i int, s *goquery.Selection) {
			if i != 1 {
				return
			}

			s.Find("li").Each(func(j int, s *goquery.Selection) {
				text := strings.Fields(s.Text())
				if len(text) > 2 {
					txt := strings.Join(text[2:], " ")
					f, err := strconv.ParseFloat(txt, 64)
					if err == nil {
						dNum[text[0]] = f
						return
					}
					dTxt[text[0]] = txt
				}
			})
		})

		if len(dNum)+len(dTxt) != 0 {
			measurements[subID] = Measurement{SubjID: subID, DataTxt: dTxt, DataNum: dNum}
		}

	}

	b, err = json.MarshalIndent(measurements, "", "	")
	if err != nil {
		fmt.Printf("error encodin JSON: %v", err)
	}
	return
}
