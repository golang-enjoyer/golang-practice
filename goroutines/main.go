// main.go

package main

import (
	"crypto/md5"
	"fmt"
	"hash/crc32"
	"sort"
	"strconv"
	"strings"
	"sync"
)

var mu sync.Mutex

func DataSigner(algo, data string) string {
	if algo == "md5" {
		return fmt.Sprintf("%x", md5.Sum([]byte(data)))
	} else if algo == "crc32" {
		return fmt.Sprintf("%v", crc32.ChecksumIEEE([]byte(data)))
	}
	return data
}

func SingleHash(in, out chan interface{}, wg *sync.WaitGroup) {
	defer wg.Done()

	var innerWg sync.WaitGroup

	for dataRaw := range in {
		innerWg.Add(1)
		data, ok := dataRaw.(int)
		if !ok {
			fmt.Errorf("unable to convert data to int")
			continue
		}

		go func(data int) {
			defer innerWg.Done()
			strData := strconv.Itoa(data)
			crc32Md5 := DataSignerCrc32(DataSignerMd5(strData))
			crc32Data := DataSignerCrc32(strData)
			result := crc32Data + "~" + crc32Md5
			out <- result
		}(data)
	}

	innerWg.Wait()
	close(out)
}

func MultiHash(in, out chan interface{}, wg *sync.WaitGroup) {
	defer wg.Done()

	var innerWg sync.WaitGroup

	for dataRaw := range in {
		innerWg.Add(1)
		data, ok := dataRaw.(string)
		if !ok {
			fmt.Errorf("unable to convert data to string")
			continue
		}

		go func(data string) {
			defer innerWg.Done()
			var result string
			for th := 0; th < 6; th++ {
				crc32Result := DataSignerCrc32(strconv.Itoa(th) + data)
				result += crc32Result
			}
			out <- result
		}(data)
	}

	innerWg.Wait()
	close(out)
}

func CombineResults(in, out chan interface{}, wg *sync.WaitGroup) {
	defer wg.Done()

	var results []string

	for dataRaw := range in {
		data, ok := dataRaw.(string)
		if !ok {
			fmt.Errorf("unable to convert data to string")
			continue
		}
		results = append(results, data)
	}

	sort.Strings(results)
	out <- fmt.Sprint(strings.Join(results, "_"))
	close(out)
}

func ExecutePipeline(jobs ...job) {
	var chans []chan interface{}
	for i := 0; i < len(jobs)+1; i++ {
		chans = append(chans, make(chan interface{}, MaxInputDataLen))
	}

	var wg sync.WaitGroup
	for i, j := range jobs {
		wg.Add(1)
		go j(chans[i], chans[i+1], &wg)
	}

	wg.Wait()
}

func main() {
	inputData := []int{0, 1, 1, 2, 3, 5, 8}

	hashSignJobs := []job{
		func(in, out chan interface{}, wg *sync.WaitGroup) {
			defer wg.Done()
			defer close(out)
			for _, fibNum := range inputData {
				in <- fibNum
				out <- fibNum
			}
		},
		SingleHash,
		MultiHash,
		CombineResults,
		func(in, out chan interface{}, wg *sync.WaitGroup) {
			defer wg.Done()
			dataRaw := <-in
			_, ok := dataRaw.(string)
			if !ok {
				fmt.Errorf("unable to convert data to string")
				return
			}
			fmt.Println(dataRaw)
		},
	}

	ExecutePipeline(hashSignJobs...)
}
