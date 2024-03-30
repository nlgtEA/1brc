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
	"unsafe"

	"github.com/dolthub/swiss"
)

const (
	READ_BUFFER_SIZE = 1024 * 1024 * 20
	CONCURENT_GRADE  = 6
)

var file_path = flag.String("file", "test_cases/measurements-10.txt", "path to the file")
var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")

func byteSlice2String(bs []byte) string {
	return *(*string)(unsafe.Pointer(&bs))
}

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

func processReadBuffer(chunk_chans chan []byte, resultChan chan *swiss.Map[string, []int]) {
	resultMap := swiss.NewMap[string, []int](42)

	for validChunk := range chunk_chans {
		prevIdx := 0
		name := ""
		temp := 0

		for idx, b := range validChunk {
			if b == ';' {
				name = byteSlice2String(validChunk[prevIdx:idx])
				prevIdx = idx + 1
			} else if b == '\n' {
				temp = parseTempToInt(validChunk[prevIdx:idx])
				prevIdx = idx + 1

				if v, ok := resultMap.Get(name); ok {
					if temp < v[0] {
						v[0] = temp
					}
					if temp > v[2] {
						v[2] = temp
					}
					v[1] += temp
					v[3]++
					resultMap.Put(name, v)
				} else {
					resultMap.Put(name, []int{temp, temp, temp, 1})
				}

			}
		}
	}

	resultChan <- resultMap
}

func readChunk(f *os.File, chunkChan chan []byte, wg *sync.WaitGroup, resultChan chan *swiss.Map[string, []int]) {
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
		validChunk = append(validChunk[:size], readBuffer[:lastNewlineIdx+1]...)

		to_send := make([]byte, size+lastNewlineIdx+1)
		copy(to_send, validChunk)
		chunkChan <- to_send

		leftOverSize = copy(leftOver, readBuffer[lastNewlineIdx+1:n])
	}

	close(chunkChan)

	wg.Wait()
	close(resultChan)
}

func evaluate(inp string) string {
	chunksChan := make(chan []byte, 10)
	resultChan := make(chan *swiss.Map[string, []int], 10)

	// {"city": [min, sum, max, count]}
	resultMap := swiss.NewMap[string, []int](42)

	f, err := os.Open(inp)
	check(err)
	defer f.Close()

	var wg sync.WaitGroup
	for i := 0; i < CONCURENT_GRADE; i++ {
		wg.Add(1)
		go func() {
			processReadBuffer(chunksChan, resultChan)
			wg.Done()
		}()
	}

	go readChunk(f, chunksChan, &wg, resultChan)

	for r := range resultChan {
		r.Iter(func(k string, v []int) (stop bool) {
			if val, ok := resultMap.Get(k); ok {
				val[0] = min(val[0], v[0])
				val[1] = val[1] + v[1]
				val[2] = max(val[2], v[2])
				val[3] = val[3] + v[3]
				resultMap.Put(k, val)
			} else {
				resultMap.Put(k, v)
			}

			return false
		})

	}

	computedValues := make([]computedValue, resultMap.Count())

	count := 0
	resultMap.Iter(func(k string, v []int) (stop bool) {
		computedValues[count] = computedValue{k, float64(v[0]) / 10, math.Round(float64(v[1])/float64(v[3])) / 10, float64(v[2]) / 10}
		count++
		return false
	})

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
