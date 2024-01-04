// TODO: timeout

package pkg

import (
	"bytes"
	
	"math"
	"github.com/fatih/color"
	"net/http"
	"sync"
	"time"
)

var red = color.New(color.BgRed).Add(color.Underline)
var green = color.New(color.BgGreen).Add(color.Underline)
var yellow = color.New(color.BgHiYellow).Add(color.Underline)

type HttpRequest struct {
	Url    string
	Method string
	Header map[string]interface{}
	Body   []byte
}

// peak hits per sec
// peak transfer speed
// latency
// point where

type RequestStat struct {
	TotalTime        float64
	TotalNumReq      int64
	ConcurrentReq    int64
	NumSuccess       int64
	NumFailed        int64
	AvgRequestPerSec int64
	PeakLoadCapacity int
	Latency          float64
}

type MicroRequestStat struct {
	TotalTime        float64
	ByteSentSize     int64
	ByteRecievedSize int64
	StatusCode       int
	Failed           bool
}

func MakePlainHttpCall(url string) error {
	_, err := http.Get(url)

	if err != nil {
		return err
	}

	return nil
}

func MakeComplexHttpCall(reqData HttpRequest) MicroRequestStat {

	startTime := time.Now()

	req, err := http.NewRequest(reqData.Method, reqData.Url, bytes.NewBuffer(reqData.Body))


	headerLen := 0
	for key, value := range reqData.Header {
		headerLen += len(key)
		headerLen += len(value.(string))
		req.Header.Add(key, value.(string))
	}

	if err != nil {
		return MicroRequestStat{}
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return MicroRequestStat{}
	}
	defer resp.Body.Close()

	headerResLen := 0
	for key, value := range resp.Header {
		headerResLen += len(key)
		for _, i := range value {
			headerResLen += len(i)
		}
	}

	elapsedTime := time.Since(startTime)
	
	Ms := MicroRequestStat{TotalTime: float64(elapsedTime.Milliseconds()),
		StatusCode: resp.StatusCode,
		ByteSentSize: req.ContentLength + int64(headerLen),
		ByteRecievedSize: resp.ContentLength + int64(headerResLen),
	}
	
	if resp.StatusCode > 199 && resp.StatusCode < 300 {
		// green.Println("congratulations it was a success with ", resp.StatusCode)
		Ms.Failed = false
	} else if resp.StatusCode > 299 && resp.StatusCode < 400 {
		// yellow.Println("redirect with ", resp.StatusCode)
		Ms.Failed = false
	} else {
		Ms.Failed = true
		// red.Println("failed with ", resp.StatusCode)
	}

	return Ms

}

func MakeConnReq(reqData HttpRequest, connReq int64) []MicroRequestStat {

	var wg sync.WaitGroup
	var t []MicroRequestStat

	// var Store []float64

	// rs := make(chan RequestStat, connReq + 1)

	for i := 0; int64(i) < connReq; i++ {
		
		wg.Add(1)
		func() {
			defer wg.Done()
			t = append(t, MakeComplexHttpCall(reqData))
			// rs <- MakeComplexHttpCall(reqData)
		}()
	}
	wg.Wait()

	// for r := range rs {
	// 	fmt.Println("i recieved one")
	// 	t = append(t, r)
	// }

	// close(rs)

	return t
}

func MultipleRequest(reqData HttpRequest, numReq int64, connReq int64) RequestStat {
	tList := [][]MicroRequestStat{}

	TInter := int(numReq / connReq)
	remind := numReq % connReq

	if remind > 0 {
		TInter += 1
	}
	if remind < 1 {
		remind = connReq
	}

	for i := 0; i < TInter; i++ {
		if i == TInter-1 {
			MCR := MakeConnReq(reqData, remind)
			tList = append(tList, MCR)
			break
		}
		MCR := MakeConnReq(reqData, connReq)
		tList = append(tList, MCR)
	}

	_, _, rs := analysis(tList)

	return rs
}

func analysis(r [][]MicroRequestStat) (tTime float64, reqSec float64, reqstat RequestStat) {

	// num of succesfull req
	// num of failed req
	// peak hits per sec
	// peak transfer speed
	
	index := 0.0


	NumSuccess := 0
	NumFail := 0

	reqSecArr := []float64{}
	SecTime := []float64{}

	totalNumReq := 0
	cunnReq := 0

	for i := 0; i < len(r); i++ {
		numReq := 0
		timeMicro := 0.0
		cunnReq = 0
		for j := 0; j < len(r[i]); j++ {
			cunnReq += 1
			
			if math.IsInf(r[i][j].TotalTime, 1) {
				continue
			}
	
			if r[i][j].StatusCode < 300 && r[i][j].StatusCode > 199 {
				NumSuccess += 1
			} else {
				NumFail += 1
			}

			totalNumReq += 1
			numReq += 1
			index += 1
			tTime += r[i][j].TotalTime
			timeMicro += r[i][j].TotalTime
			SecTime = append(SecTime, r[i][j].TotalTime)
		}

		reqSecArr = append(reqSecArr, float64(numReq)/(timeMicro * 0.001))	
	}

	_, PeakHit := highestFloat(reqSecArr)	

	tTime = tTime * 0.001
	reqSec = findAverage(reqSecArr)

	reqstat = RequestStat{
		TotalTime: tTime,
		TotalNumReq: int64(totalNumReq),
		ConcurrentReq: int64(cunnReq), 
		NumSuccess: int64(NumSuccess),
		NumFailed: int64(NumFail),
		AvgRequestPerSec: int64(reqSec),
		PeakLoadCapacity: int(PeakHit),
		Latency: findAverage(SecTime),
	}

	return
}
