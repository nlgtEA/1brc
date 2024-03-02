package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

var file_path = flag.String("file", "test_cases/measurements-10.txt", "path to the file")

func main() {
	flag.Parse()

	output := evaluate(*file_path)

	fmt.Printf("{%s}\n", output)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type computedValue struct {
	city string
	min  float64
	avg  float64
	max  float64
}

func evaluate(inp string) string {
	// {"city": [min, sum, max, count]}
	resultMap := make(map[string][]int)

	f, err := os.Open(inp)
	check(err)
	defer f.Close()

	r := bufio.NewReader(f)

	for {
		line, err := r.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			} else {
				panic(err)
			}
		}

		parsed := strings.Split(line, ";")
		name := parsed[0]
		fTemp, _ := strconv.ParseFloat(parsed[1][:len(parsed[1])-1], 64)
		temp := int(fTemp * 10)

		if _, ok := resultMap[name]; ok {
			if temp < resultMap[name][0] {
				resultMap[name][0] = temp
			}
			if temp > resultMap[name][2] {
				resultMap[name][2] = temp
			}
			resultMap[name][1] += temp
			resultMap[name][3]++
		} else {
			resultMap[name] = []int{temp, temp, temp, 1}
		}
	}

	computedValues := make([]computedValue, len(resultMap))

	count := 0
	for k, v := range resultMap {
		computedValues[count] = computedValue{k, float64(v[0]) / 10, math.Round(float64(v[1])/float64(v[3])) / 10, float64(v[2]) / 10}
		count++
	}
	sort.Slice(computedValues, func(i, j int) bool {
		return computedValues[i].city < computedValues[j].city
	})

	strBuilder := strings.Builder{}
	for _, v := range computedValues {
		strBuilder.WriteString(v.city)
		strBuilder.WriteString("=")
		strBuilder.WriteString(fmt.Sprintf("%.1f", v.min))
		strBuilder.WriteString("/")
		strBuilder.WriteString(fmt.Sprintf("%.1f", v.avg))
		strBuilder.WriteString("/")
		strBuilder.WriteString(fmt.Sprintf("%.1f", v.max))
		strBuilder.WriteString(", ")
	}

	return strBuilder.String()[:strBuilder.Len()-2]
}
