package test

import (
	"fmt"
	"testing"
	"github.com/msalbrain/turbotanuki/pkg"
)




func TestCaller(t *testing.T) {
	pkg.MakePlainHttpCall("https://google.com")

}

func TestComplexCaller(t *testing.T) {
	p := pkg.HttpRequest{

		Url: "https://google.com",
		
		Method: "GET",
		Body: []byte(`{"msg": "hello there"}`),
		Header: map[string]interface{}{"authorization": "Bearer df7yujhfg68uijhytyu"},
	}

	i := 0
	Elapsed := 0.0
	for {
		stat := pkg.MakeComplexHttpCall(p)
		Elapsed += stat.TotalTime

		if i > 10 {
			break
		}

		i += 1
	}

	fmt.Printf("Total time taken is %f avg time id %f \n\n", Elapsed,  Elapsed/10)

}

