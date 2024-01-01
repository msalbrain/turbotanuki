package pkg

import (
	"bytes"
	"sync"
	"fmt"
	"net/http"
	"time"
	"github.com/fatih/color"
)


var red = color.New(color.BgRed).Add(color.Underline)
var green = color.New(color.BgGreen).Add(color.Underline)
var yellow = color.New(color.BgHiYellow).Add(color.Underline)

type HttpRequest struct {
	Url string

	Method string

	Header map[string]interface{}

	Body []byte
}


type RequestStat struct {
	TotalTime float64
	requestPerSec int64
	ConcurrentReq int64
	ErrorRate float32
	PeakLoadCapacity int
}


func MakePlainHttpCall(url string) error {
	resp, err := http.Get(url)

	if err != nil {
		return err
	}

	fmt.Println("status code: ", resp.StatusCode)
	return nil
}


func MakeComplexHttpCall(reqData HttpRequest) RequestStat {
	
	startTime := time.Now()

	req, err := http.NewRequest(reqData.Method, reqData.Url, bytes.NewBuffer(reqData.Body))
	
	for key, value := range reqData.Header{
		req.Header.Add(key, value.(string))
	}

	if err != nil {
		fmt.Println("Error creating request:", err)
		return RequestStat{}
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return RequestStat{}
	}
	defer resp.Body.Close()	

	if resp.StatusCode > 199 && resp.StatusCode < 300 {
		green.Println("congratulations it was a success with ", resp.StatusCode)		
	} else if resp.StatusCode > 299 && resp.StatusCode < 400 {
		yellow.Println("redirect with ", resp.StatusCode)
	} else {
		red.Println("failed with ", resp.StatusCode)
	} 
	
	elapsedTime := time.Since(startTime)

	// fmt.Println("Elapsed Time:", elapsedTime)
	client.CloseIdleConnections()
	return RequestStat{TotalTime: float64(elapsedTime.Milliseconds())}

}

func MakeConnReq(reqData HttpRequest, connReq int64) []RequestStat {
	
	var wg sync.WaitGroup
	var t []RequestStat

	// rs := make(chan RequestStat, connReq + 1)

	for i := 0; int64(i) < connReq; i++ {
		// fmt.Println("this is conn loop: ", i)
		wg.Add(1)
		func() {
			defer wg.Done()
			// fmt.Println("\nthis is func:", i, "start")
			t = append(t, MakeComplexHttpCall(reqData))
			// rs <- MakeComplexHttpCall(reqData)
			// fmt.Println("this is func:", i, "end\n\n.")		
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
	tList := [][]RequestStat{}	
	
	TInter := int(numReq / connReq)
	remind := numReq % connReq
	
	
	if remind > 0 {
		TInter += 1
	}

	for i := 0; i < TInter; i++ { 
		// fmt.Println("tinter", i)
		if i == TInter - 1 {
			tList = append(tList, MakeConnReq(reqData, remind))
			break
		}
		tList = append(tList, MakeConnReq(reqData, connReq))
	}

	return RequestStat{TotalTime: analysis(tList)}
}




func analysis(r [][]RequestStat) (tTime float64) {
	
	index := 0
	for i := 0; i < len(r); i++ {
        for j := 0; j < len(r[i]); j++ {
			index += 1
			tTime += r[i][j].TotalTime
		}
	}
	tTime = tTime/float64(index)

	return
}
