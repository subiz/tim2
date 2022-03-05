package tim2

import (
	"fmt"
	"encoding/csv"
	"io"
	"os"
	"strconv"
	"strings"
)

func Report() {
	csvfile, err := os.Open("./term")
	if err != nil {
		panic(err)
	}
	defer csvfile.Close()

	data := map[string]map[string]int{}
	i := 0
	r := csv.NewReader(csvfile)
	for {
		i++
		if i%100000 == 0 {
			fmt.Println("I", i)
		}
		rec, err := r.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}

		col, acc, term := rec[0], rec[1], rec[2]
		if data[col+"_"+acc] == nil {
			data[col+"_"+acc] = make(map[string]int)
		}
		data[col+"_"+acc][term]++
	}

	// trim
	trimmedData := map[string]map[string]int{}
	for k, v := range data {
		for term, count := range v {
			if trimmedData[k] == nil {
				trimmedData[k] = make(map[string]int)
			}
			if count < 2000 {
				continue
			}
			trimmedData[k][term] = count
		}
	}

	for k, v := range data {
		fmt.Println("TERMS OF", k, len(v))
	}

	data = trimmedData

	for k, v := range data {
		// if !strings.Contains(k, "acrdslqbraholyxzwjbu")  {
		//continue
		//}
		if len(v) == 0 {
			continue
		}
		fmt.Println("DISTRIBUTION OF", k)
		top10 := topK(v, 10)
		fmt.Println(drawGraph(top10))
	}
}

func formatLabel(label string, length int) string {
	if len(label) < length {
		return strings.Repeat(" ", length-len(label)) + label
	}

	return label[:length-3] + "..."
}

func topK(m map[string]int, k int) map[string]int {
	out := make(map[string]int)

	for i := 0; i < k; i++ {
		maxi := ""
		for word, freq := range m {
			if _, has := out[word]; has {
				continue
			}

			if maxi == "" {
				maxi = word
			}
			if m[maxi] < freq {
				maxi = word
			}
		}
		if maxi != "" {
			out[maxi] = m[maxi]
		}
	}
	return out
}

func drawGraph(full map[string]int) string {
	// convert data to sort
	labels := []string{}
	data := []int{}
	for k, v := range full {
		labels = append(labels, k)
		data = append(data, v)
	}

	// sort data
	for i := 0; i < len(data); i++ {
		for j := i + 1; j < len(data); j++ {
			if data[i] < data[j] { // swap
				data[i], data[j] = data[j], data[i]
				labels[i], labels[j] = labels[j], labels[i]
			}
		}
	}

	s := ""
	// draw graph
	for i, d := range data {
		l := formatLabel(labels[i], 50)
		numStroke := d * 120 / data[0] // max 60 strokes
		line := strings.Repeat("#", numStroke) + strings.Repeat(" ", 120-numStroke)
		s += "\n" + l + " " + line + "  " + strconv.Itoa(d)
	}
	return s
}
