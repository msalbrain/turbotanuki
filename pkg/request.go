package pkg

import (
	"bytes"
	"fmt"

	"context"
	"math"
	"net/http"
	"sync"
	"time"

	"github.com/fatih/color"
)

var red = color.New(color.BgRed).Add(color.Underline)
var green = color.New(color.BgGreen).Add(color.Underline)
var yellow = color.New(color.BgHiYellow).Add(color.Underline)

type HttpRequest struct {
	Url    string
	Method string
	Header map[string]interface{}
	Body   []byte
	Timeout int64 // in millisecond
}


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


// make general http request
func MakePlainHttpCall(url string) error {
	_, err := http.Get(url)

	if err != nil {
		return err
	}

	return nil
}


// make http call with more parmeters (in body and header)
func MakeComplexHttpCall(reqData HttpRequest) MicroRequestStat {

	startTime := time.Now()
	
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(reqData.Timeout) * time.Millisecond)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, reqData.Method, reqData.Url, bytes.NewBuffer(reqData.Body))
	
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
	elapsedTime := time.Since(startTime)
	if err != nil {
		return MicroRequestStat{Failed: true, StatusCode: 500, TotalTime: float64(elapsedTime.Milliseconds())}
	}
	
	defer resp.Body.Close()

	headerResLen := 0
	for key, value := range resp.Header {
		headerResLen += len(key)
		for _, i := range value {
			headerResLen += len(i)
		}
	}

	
	
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


/*
	make number of request in `connReq` and returns arr of stat on each request
*/
func MakeConnReq(reqData HttpRequest, connReq int64) []MicroRequestStat {

	var wg sync.WaitGroup

	t := []MicroRequestStat{}
	rs := make(chan MicroRequestStat)
	
	for i := 0; int64(i) < connReq; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			rs <- MakeComplexHttpCall(reqData)
			}()
		}
		
		go func() { // no wait group cause all it does is listen on rs and append 
			for r := range rs {
				t = append(t, r)
			}		
		}()
				
		wg.Wait()
		return t
}



/*
	batches up request and make them
*/
func MultipleRequest(reqData HttpRequest, numReq int64, connReq int64, reqChan chan Uper) RequestStat {

	startTime := time.Now()

	tList := [][]MicroRequestStat{} // contains an array of an array of the stats from batch request

	TInter := int(numReq / connReq) 
	remind := numReq % connReq

	if remind > 0 {
		TInter += 1
	}
	if remind < 1 {
		remind = connReq
	}

	var u []float64

	for i := 0; i < TInter; i++ {
		
		if i == TInter-1 {
			microStat := MakeConnReq(reqData, remind)		
			u = append(u, 1/float64(TInter))
			reqChan <- Uper{By: 1/float64(TInter)} // makes update to progress view
			tList = append(tList, microStat)
			break
		}
		
		microStat := MakeConnReq(reqData, connReq)
		u = append(u, 1/float64(TInter))
		reqChan <- Uper{By: 1/float64(TInter)} // makes update to progress view

		tList = append(tList, microStat)
	}

	close(reqChan)

	fmt.Println("array of u", u)	

	elapsedTime := time.Since(startTime)
	rs := analysis(tList, elapsedTime.Seconds(), connReq)

	return rs
}


/*
	performs analysis on request stats and returns result
*/
func analysis(r [][]MicroRequestStat, totalTime float64, cunnReq int64) (reqstat RequestStat) {

	index := 0.0

	NumSuccess := 0
	NumFail := 0

	reqSecArr := []float64{}
	SecTime := []float64{}

	totalNumReq := 0
	
	for i := 0; i < len(r); i++ {
		numReq := 0
		timeMicro := 0.0
		
		for j := 0; j < len(r[i]); j++ {
			totalNumReq += 1
			
			if math.IsInf(r[i][j].TotalTime, 1) {
				NumFail += 1
				continue
			}
	
			if r[i][j].Failed {
				NumFail += 1
			} else {
				NumSuccess += 1
			}

			numReq += 1
			index += 1

			timeMicro += r[i][j].TotalTime

			SecTime = append(SecTime, r[i][j].TotalTime)
		}

		reqSecArr = append(reqSecArr, float64(numReq)/(timeMicro * 0.001))	
	}

	_, PeakHit := highestFloat(filterInfinity(reqSecArr))
	
	reqSec := findAverage(reqSecArr)

	reqstat = RequestStat{
		TotalTime: totalTime,
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
