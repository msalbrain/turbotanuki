// TODO: timeout

/*
Copyright © 2023 NAME HERE <salbiz2021@gmail.com>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"sync"
	"github.com/msalbrain/turbotanuki/pkg"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"log"
	"math"
	"net/url"
	"os"
	"strconv"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

var URL string

var numRequest int64
var connRequest int64

var method string
var header string
var body string

var file string

func GenDrawable() {

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	sinData := func() [][]float64 {
		n := 220
		data := make([][]float64, 2)
		data[0] = make([]float64, n)
		data[1] = make([]float64, n)
		for i := 0; i < n; i++ {
			data[0][i] = 1 + math.Sin(float64(i)/5)
			data[1][i] = 1 + math.Cos(float64(i)/5)
		}
		return data
	}()

	rect := []int{0, 0, 100, 40}
	
	p2 := widgets.NewPlot()

	
	p2.Title = "request perfomance"
	p2.Marker = widgets.MarkerDot
	p2.Data = make([][]float64, 2)
	p2.Data[0] = []float64{1, 2, 3, 4, 5}
	p2.Data[1] = sinData[1][4:]
	p2.Data[1] = []float64{}
	p2.SetRect(rect[0], rect[1], rect[2], rect[3])
	p2.AxesColor = ui.ColorWhite
	p2.LineColors[0] = ui.ColorRed
	p2.PlotType = widgets.ScatterPlot

	ui.Render(p2)

	return
	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		case "e":
			rect := []int{rect[0], rect[1] + 1, rect[2], rect[3] + 1}
			p2.SetRect(rect[0], rect[1], rect[2], rect[3])
			ui.Render(p2)
		case "w":
			p2.Data[0] = append(p2.Data[0], p2.Data[0][len(p2.Data[0]) -1 ] + 2)
			ui.Render(p2)
		}
	}

}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "turbotanuki",
	Short: "\nyour all in one http load tester",

	// don't touch the below for any reason
	Long: `
	___________          ___.            ___________                    __   .__ 
	\__    ___/_ ________\_ |__   ____   \__    ___/____    ____  __ __|  | _|__|
	  |    | |  |  \_  __ \ __ \ /  _ \    |    |  \__  \  /    \|  |  \  |/ /  |
	  |    | |  |  /|  | \/ \_\ (  <_> )   |    |   / __ \|   |  \  |  /    <|  |
	  |____| |____/ |__|  |___  /\____/    |____|  (____  /___|  /____/|__|_ \__|
				  \/                        \/     \/           \/   	
	`,

	Run: func(cmd *cobra.Command, args []string) {


		if numRequest < connRequest {
			fmt.Println("\nyour maths is out of order,\nthere are far more number of conncurrent request than number of request to be made\n ")
			return
		}

		if len(args) < 1 && URL == "" {
			fmt.Println("no url provided. check `turbotanuki --help`")
		} else {

			for i, maybeUrl := range args {
				_, err := url.Parse(maybeUrl)
				if err != nil {
					continue
				}

				URL = args[i]
				break
			}

			if URL == "" {
				fmt.Println("invalid url provided")
				return
			}
		} 

		if URL != "" {
			fmt.Println("url: ", URL, " number of request: ", numRequest, " conncurrent request: ", connRequest,
				" method: ", method, " header: ", header, " body: ", body, " file: ", file)

			var dataHeader map[string]interface{}

			if header != "" {
				err := json.Unmarshal([]byte(header), &dataHeader)

				if err != nil {
					fmt.Println("wrong parameters passed", err)
					return
				}
			}

			p := pkg.HttpRequest{
				Url:    URL,
				Method: method,
				Body:   []byte(body),
				Header: dataHeader,
			}

			var wg sync.WaitGroup

			rs := make(chan pkg.RequestStat)

			wg.Add(1)
			
			go func () {
				rs <- pkg.MultipleRequest(p, numRequest, connRequest)
				// res = pkg.MultipleRequest(p, numRequest, connRequest)
				wg.Done()
			}()

			res := <- rs

			metrics := [][]string{
				{"Total time taken (S)", strconv.FormatFloat(res.TotalTime, 'f', -1, 64)},
				{"Total Number of request", strconv.FormatInt(res.TotalNumReq, 10)},
				{"Number of concurrent reques", strconv.FormatInt(res.ConcurrentReq, 10)},
				{"AVG Requests Per Second (RPS)", strconv.FormatFloat(float64(res.AvgRequestPerSec), 'f', -1, 64)},
				{"Number of success", strconv.FormatInt(res.NumSuccess, 10)},
				{"Number of faliure", strconv.FormatInt(res.NumFailed, 10)},
				{"average Latency (ms)", strconv.FormatFloat(res.Latency, 'f', 5, 64)},
			}

			table := tablewriter.NewWriter(os.Stdout)

			table.SetHeader([]string{"Metric", "Value"})

			for _, v := range metrics {
				table.Append(v)
			}

			// wg.Add(1)
			// go func() {
			// 	GenDrawable()
			// 	wg.Done()
			// }()

			wg.Wait()
			table.Render()
		}

	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

	rootCmd.Flags().StringVarP(&URL, "url", "u", "", "this is the url for the deed")
	rootCmd.Flags().Int64VarP(&numRequest, "numreq", "n", 1, "this is the total number of request to be made")
	rootCmd.Flags().Int64VarP(&connRequest, "cunnreq", "c", 1, "this is the number of concurrent request to be made at a time")

	rootCmd.Flags().StringVarP(&file, "file", "f", "", "this is to be followed by a file location that contains a tanuki directives(commands) , it allows for more complex request")

	rootCmd.Flags().StringVarP(&method, "method", "m", "", "this specifies the method to be used for the http request")
	rootCmd.Flags().StringVarP(&header, "header", "d", ``, "this specifies the header to be used for the http request")
	rootCmd.Flags().StringVarP(&body, "body", "b", "", "this specifies the body to be used for the http request")
}
