package browser

import (
	"fmt"
	"math"
	"sort"

	"github.com/marcelfarres/hrtftools/parser"
)

type maxmin struct {
	max float64
	min float64
}

type kv struct {
	Key   string
	Value float64
}

func MinDist(d map[string]parser.Measurement, resNum int, keys []string, values []float64) (r map[string]float64) {

	var (
		keyMaxMin map[string]*maxmin
	)

	keyMaxMin = make(map[string]*maxmin, 0)
	// Filter Map
	for _, key := range keys {
		keyFound := false
		keyMaxMin[key] = &maxmin{max: -math.MaxFloat64, min: math.MaxFloat64}

		for _, m := range d {
			for nKey, v := range m.DataNum {
				if nKey == key {
					keyFound = true
					if v > keyMaxMin[key].max {
						keyMaxMin[key].max = v
						continue
					}
					if v < keyMaxMin[key].min {
						keyMaxMin[key].min = v
					}
				}

			}
		}
		if !keyFound {
			fmt.Printf("Key %s not present in any participant (keys listed:%v)", key, keys)
			return nil
		}
	}
	subDist := make(map[string]float64)
	for i, key := range keys {
		ponderation := 1 / (keyMaxMin[key].max - keyMaxMin[key].min)
		for _, m := range d {
			subDist[m.SubjID] = math.Abs(m.DataNum[key]-values[i]) * ponderation
			// fmt.Println(math.Abs(m.DataNum[key]-values[i])*ponderation, m.DataNum[key], values[i], ponderation)
		}
	}

	var ss []kv
	for k, v := range subDist {
		ss = append(ss, kv{k, v})
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value < ss[j].Value
	})

	r = make(map[string]float64, resNum)
	for i := 0; i < resNum; i++ {
		r[ss[i].Key] = ss[i].Value
	}

	return r
}
