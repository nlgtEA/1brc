package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"runtime/pprof"
	"sort"
	"strings"
)

var file_path = flag.String("file", "test_cases/measurements-10.txt", "path to the file")
var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")

func main() {
	flag.Parse()

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		if err := pprof.StartCPUProfile(f); err != nil {
			panic(err)
		}
		defer pprof.StopCPUProfile()
	}

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

func parseTempToInt(rawTemp []byte) int {
	temp := 0
	isNegative := false

	if rawTemp[0] == 45 {
		isNegative = true
		rawTemp = rawTemp[1:]
	}

	l := len(rawTemp)
	if l == 4 {
		temp = int(rawTemp[0])*100 + int(rawTemp[1])*10 + int(rawTemp[3]) - 5328
	} else {
		temp = int(rawTemp[0])*10 + int(rawTemp[2]) - 528
	}

	if isNegative {
		temp = -temp
	}

	return temp
}

func evaluate(inp string) string {
	// {"city": [min, sum, max, count]}
	resultMap := make(map[string][]int)

	f, err := os.Open(inp)
	check(err)
	defer f.Close()

	r := bufio.NewReader(f)

	for {
		line, err := r.ReadSlice('\n')
		if err != nil {
			if err == io.EOF {
				break
			} else {
				panic(err)
			}
		}

		line = line[:len(line)-1]

		parsed := bytes.Split(line, []byte{';'})
		name := string(parsed[0])
		temp := parseTempToInt(parsed[1])

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
