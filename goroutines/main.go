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
			mu.Lock()
			crc32Md5 := DataSignerCrc32(DataSignerMd5(strData))
			mu.Unlock()
			crc32Data := DataSignerCrc32(strData)
			result := crc32Data + "~" + crc32Md5
			out <- result
		}(data)
	}

	innerWg.Wait()
	close(out)
}

func MultiHash(in, out chan interface{}) {
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

func CombineResults(in, out chan interface{}) {
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
		defer wg.Done()
		defer close(chans[i+1])
		go j(chans[i], chans[i+1])
	}

	wg.Wait()
}

func main() {
	now := time.Now()
	defer func() {
		fmt.Println("SingleHash done in ", time.Since(now))
	}()
	inputData := []int{0, 1, 1, 2, 3, 5, 8}

	hashSignJobs := []job{
		func(in, out chan interface{}) {
			defer close(out)
			for _, fibNum := range inputData {
				in <- fibNum
				out <- fibNum
			}
		},
		SingleHash,
		MultiHash,
		CombineResults,
		func(in, out chan interface{}) {
			dataRaw := <-in
			defer close(out)
			_, ok := dataRaw.(string)
			fmt.Println("SingleHash done in ", time.Since(now))
			if !ok {
				fmt.Errorf("unable to convert data to string")
				return
			}
			fmt.Println(dataRaw)
		},
	}

	ExecutePipeline(hashSignJobs...)
}
