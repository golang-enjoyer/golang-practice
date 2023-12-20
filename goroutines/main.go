package main

import (
	"crypto/md5"
	"fmt"
	"hash/crc32"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var mu sync.Mutex

type job func(in, out chan interface{})

const (
	MaxInputDataLen = 100
)

var (
	dataSignerOverheat uint32 = 0
	DataSignerSalt            = ""
)

var OverheatLock = func() {
	for {
		if swapped := atomic.CompareAndSwapUint32(&dataSignerOverheat, 0, 1); !swapped {
			fmt.Println("OverheatLock happend")
			time.Sleep(time.Second)
		} else {
			break
		}
	}
}

var OverheatUnlock = func() {
	for {
		if swapped := atomic.CompareAndSwapUint32(&dataSignerOverheat, 1, 0); !swapped {
			fmt.Println("OverheatUnlock happend")
			time.Sleep(time.Second)
		} else {
			break
		}
	}
}

var DataSignerMd5 = func(data string) string {
	OverheatLock()
	defer OverheatUnlock()
	data += DataSignerSalt
	dataHash := fmt.Sprintf("%x", md5.Sum([]byte(data)))
	time.Sleep(10 * time.Millisecond)
	return dataHash
}

var DataSignerCrc32 = func(data string) string {
	data += DataSignerSalt
	crcH := crc32.ChecksumIEEE([]byte(data))
	dataHash := strconv.FormatUint(uint64(crcH), 10)
	time.Sleep(time.Second)
	return dataHash
}

func SingleHash(in, out chan interface{}) {
	now := time.Now()
	defer func() {
		fmt.Println("SingleHash done in ", time.Since(now))
	}()

	var innerWg sync.WaitGroup

	for data := range in {
		innerWg.Add(1)
		go func(data int) {
			strData := strconv.Itoa(data)
			mu.Lock()
			Md5 := DataSignerMd5(strData)
			mu.Unlock()
			crc32Md5 := DataSignerCrc32(Md5)
			crc32Data := DataSignerCrc32(strData)
			result := crc32Data + "~" + crc32Md5
			out <- result
			innerWg.Done()
		}(data.(int))
	}

	innerWg.Wait()
}

func MultiHash(in, out chan interface{}) {
	now := time.Now()
	defer func() {
		fmt.Println("MultiHash done in ", time.Since(now))
	}()

	var outerWg sync.WaitGroup

	for data := range in {
		outerWg.Add(1)

		go func(data string) {
			var result string
			var innerWg sync.WaitGroup
			for th := 0; th < 6; th++ {
				innerWg.Add(1)
				go func(th int) {
					crc32Result := DataSignerCrc32(strconv.Itoa(th) + data)
					result += crc32Result
					innerWg.Done()
				}(th)
			}
			innerWg.Wait()
			out <- result
			outerWg.Done()
		}(data.(string))
	}
	outerWg.Wait()
}

func CombineResults(in, out chan interface{}) {
	now := time.Now()
	defer func() {
		fmt.Println("Combine Results done in ", time.Since(now))
	}()

	var innerWg sync.WaitGroup
	var results []string

	appendResult := func(result string) {
		defer innerWg.Done()
		mu.Lock()
		results = append(results, result)
		mu.Unlock()
	}

	for data := range in {
		innerWg.Add(1)
		go appendResult(data.(string))
	}

	innerWg.Wait()

	sort.Strings(results)
	out <- fmt.Sprint(strings.Join(results, "_"))
}

func ExecutePipeline(jobs ...job) {
	var chans []chan interface{}
	var wg sync.WaitGroup

	for i := 0; i < len(jobs)+1; i++ {
		chans = append(chans, make(chan interface{}, MaxInputDataLen))
	}

	for i, j := range jobs {
		wg.Add(1)
		go func(i int, j job) {
			defer wg.Done()
			defer close(chans[i+1])
			j(chans[i], chans[i+1])
		}(i, j)
	}

	wg.Wait()
}

func main() {
	inputData := []int{0, 1, 1, 2, 3, 5, 8}

	hashSignJobs := []job{
		func(in, out chan interface{}) {
			for _, fibNum := range inputData {
				out <- fibNum
			}
		},
		SingleHash,
		MultiHash,
		CombineResults,
		func(in, out chan interface{}) {
			fmt.Println(<-in)
		},
	}

	ExecutePipeline(hashSignJobs...)
}
