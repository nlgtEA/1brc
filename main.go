package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"runtime/pprof"
	"sort"
	"strings"
	"sync"
)

const (
	READ_BUFFER_SIZE = 1024 * 1024 * 20
	CONCURENT_GRADE  = 5
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

func processReadBuffer(chunk_chans chan []byte, resultChan chan map[string][]int) {
	resultMap := make(map[string][]int)

	for validChunk := range chunk_chans {
		lines := bytes.Split(validChunk, []byte{'\n'})

		for _, line := range lines {
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
	}

	resultChan <- resultMap
}

func readChunk(f *os.File, chunkChan chan []byte) {
	readBuffer := make([]byte, READ_BUFFER_SIZE)
	leftOver := make([]byte, READ_BUFFER_SIZE)
	validChunk := make([]byte, READ_BUFFER_SIZE*2)

	leftOverSize := 0

	for {
		n, err := f.Read(readBuffer)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				panic(err)
			}
		}

		lastNewlineIdx := bytes.LastIndex(readBuffer[:n], []byte{'\n'})

		size := copy(validChunk, leftOver[:leftOverSize])
		validChunk = append(validChunk[:size], readBuffer[:lastNewlineIdx]...)

		to_send := make([]byte, size+lastNewlineIdx)
		copy(to_send, validChunk)
		chunkChan <- to_send

		leftOverSize = copy(leftOver, readBuffer[lastNewlineIdx+1:n])
	}

	close(chunkChan)
}

func evaluate(inp string) string {
	chunksChan := make(chan []byte, 1000000)
	resultChan := make(chan map[string][]int)

	// {"city": [min, sum, max, count]}
	resultMap := make(map[string][]int)

	f, err := os.Open(inp)
	check(err)
	defer f.Close()

	go readChunk(f, chunksChan)

	var wg sync.WaitGroup
	// go processReadBuffer(chunksChan, resultChan)
	for i := 0; i < CONCURENT_GRADE; i++ {
		wg.Add(1)
		go func() {
			processReadBuffer(chunksChan, resultChan)
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	for r := range resultChan {
		for k, v := range r {
			if _, ok := resultMap[k]; ok {
				resultMap[k][0] = min(resultMap[k][0], v[0])
				resultMap[k][1] = resultMap[k][1] + v[1]
				resultMap[k][2] = max(resultMap[k][2], v[2])
				resultMap[k][3] = resultMap[k][3] + v[3]
			} else {
				resultMap[k] = v
			}
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
