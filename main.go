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

type weatherData struct {
	city string
	temp int
}

func readWeather(inp string, c chan weatherData) {
	f, err := os.Open(inp)
	check(err)
	defer f.Close()

	r := bufio.NewReader(f)

	for {
		line, err := r.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				close(c)
				break
			} else {
				panic(err)
			}
		}

		parsed := strings.Split(line, ";")
		name := parsed[0]
		fTemp, _ := strconv.ParseFloat(parsed[1][:len(parsed[1])-1], 64)
		temp := int(fTemp * 10)

		c <- weatherData{name, temp}
	}
}

func computeWeather(c chan weatherData, rc chan string) {
	// {"city": [min, sum, max, count]}
	resultMap := make(map[string][]int)
	for v := range c {
		if _, ok := resultMap[v.city]; ok {
			if v.temp < resultMap[v.city][0] {
				resultMap[v.city][0] = v.temp
			}
			if v.temp > resultMap[v.city][2] {
				resultMap[v.city][2] = v.temp
			}
			resultMap[v.city][1] += v.temp
			resultMap[v.city][3]++
		} else {
			resultMap[v.city] = []int{v.temp, v.temp, v.temp, 1}
		}
	}

	computedValues := make([]computedValue, 0, len(resultMap))
	computedCh := make(chan computedValue)
	count := 0

	for k, v := range resultMap {
		// computedValues = append(computedValues, computedValue{k, float64(v[0]) / 10, math.Round(float64(v[1])/float64(v[3])) / 10, float64(v[2]) / 10})
		go func(city string, val []int, c chan computedValue) {
			c <- computedValue{city, float64(val[0]) / 10, math.Round(float64(val[1])/float64(val[3])) / 10, float64(val[2]) / 10}
		}(k, v, computedCh)
		count++
	}

	// for v := range computedCh {
	// 	computedValues = append(computedValues, v)
	// }

	for count > 0 {
		v := <-computedCh
		computedValues = append(computedValues, v)
		count--
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

	rc <- strBuilder.String()[:strBuilder.Len()-2]
}

func evaluate(inp string) string {
	c := make(chan weatherData)
	rc := make(chan string)

	go readWeather(inp, c)
	go computeWeather(c, rc)
	res := <-rc

	return res
}
