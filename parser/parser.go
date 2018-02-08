package parser

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

type Measurement struct {
	SubjID  string             `json:"subj_id"`
	DataTxt map[string]string  `json:"data_txt"`
	DataNum map[string]float64 `json:"data_num"`
}

type job struct {
	url   string
	subID string
	r     chan Measurement
}

func Parse(baseURL string, totalSub int) (b []byte) {
	var (
		measurements map[string]Measurement
		err          error
	)

	jobs := make(chan job, totalSub)
	rs := make(chan Measurement, totalSub)
	maxWorkers := 20

	wg := &sync.WaitGroup{}
	wg.Add(maxWorkers)

	for i := 1; i <= maxWorkers; i++ {
		go func(i int) {
			defer wg.Done()

			for j := range jobs {
				parse(i, j)
			}
		}(i)
	}

	for i := 2; i < totalSub; i++ {
		subID := fmt.Sprintf("IRC_1%03d", i)
		url := fmt.Sprintf("%s%s", baseURL, subID)
		jobs <- job{url: url, subID: subID, r: rs}
	}
	close(jobs)
	wg.Wait()
	close(rs)

	measurements = make(map[string]Measurement)
	for m := range rs {
		measurements[m.SubjID] = m
	}

	b, err = json.MarshalIndent(measurements, "", "	")
	if err != nil {
		fmt.Printf("error encodin JSON: %v", err)
	}
	fmt.Printf("Parsing ended, total of %v subjects found", len(measurements))
	return
}

func parse(id int, j job) {
	var (
		doc *goquery.Document
	)

	fmt.Printf("Worker id:%v started job %v.\n", id, j)

	_, err := http.Get(j.url)
	if err != nil {
		fmt.Printf("URL %s does not exist, error:%s\n", j.url, err)
		return
	}
	doc, err = goquery.NewDocument(j.url)
	if err != nil {
		fmt.Printf("URL %s can't not be parsed, error:%s\n", j.url, err)
		return
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
				if err == nil && f != 0 {
					dNum[text[0]] = f
					return
				}
				dTxt[text[0]] = txt
			}
		})
	})

	if len(dNum)+len(dTxt) == 0 {
		return
	}

	j.r <- Measurement{SubjID: j.subID, DataTxt: dTxt, DataNum: dNum}
	fmt.Printf("Worker id:%v finished job %v.\n", id, j)
}
