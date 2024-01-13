/*
Copyright © 2023 NAME HERE <salbiz2021@gmail.com>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	

	"github.com/msalbrain/turbotanuki/pkg"
	// "github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	"net/url"
	"os"
	"strconv"
	"sync"
)

var URL string

var numRequest int64
var connRequest int64
var timeout int64

var method string
var header string
var body string

var file string





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
			
			var dataHeader map[string]interface{}
			
			var wg sync.WaitGroup

			model := pkg.NewProgressModel()
			letTalk := make(chan pkg.Uper)
			model.Channel = letTalk
			model.Info = fmt.Sprintln("url: ", URL, " number of request: ", numRequest, " conncurrent request: ", connRequest,
			" method: ", method, " header: ", header, " body: ", body, " file: ", file)
			model.EndReqMaker = func() {
				wg.Done()
			}

			// var wh sync.WaitGroup

			if header != "" {
				
				err := json.Unmarshal([]byte(header), &dataHeader)


				if err != nil {
					fmt.Println("wrong parameters passed", err)
					return
				}
			}

			p := pkg.HttpRequest{
				Url:     URL,
				Method:  method,
				Body:    []byte(body),
				Header:  dataHeader,
				Timeout: timeout,
			}

			var res pkg.RequestStat

			go func() {
				// This displays the progress of the request in the terminal
				model.Run()
			}()

			wg.Add(1)
			go func() {
				// This makes the request
				res = pkg.MultipleRequest(p, numRequest, connRequest, letTalk)
				wg.Done()
			}()

			wg.Wait()

			if res.TotalTime == 0.0 { // check if process was ended abruptly 
					fmt.Println("process stoped abruptly")
			} else {

			metrics := map[string]string{
				"Total time taken (S)": strconv.FormatFloat(res.TotalTime, 'f', -1, 64),
				"Total Number of request": strconv.FormatInt(res.TotalNumReq, 10),
				"Number of concurrent request": strconv.FormatInt(res.ConcurrentReq, 10),
				"AVG Requests Per Second (RPS)": strconv.FormatFloat(float64(res.AvgRequestPerSec), 'f', -1, 64),
				"Peak Requests Per Second (RPS)": strconv.FormatFloat(float64(res.PeakLoadCapacity), 'f', -1, 64),
				"Number of success": strconv.FormatInt(res.NumSuccess, 10),
				"Number of faliure": strconv.FormatInt(res.NumFailed, 10),
				"average Latency (ms)": strconv.FormatFloat(res.Latency, 'f', 5, 64),
			}

			p := pkg.PrintModel{Text: pkg.TableGen(map[string]string{"Metric": "Value"}, metrics)}

			p.Run()			
			}
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
	rootCmd.Flags().Int64VarP(&timeout, "timeout", "t", 100, "this is used to specify the timeout per request(in millisecond) default is 100")

	rootCmd.Flags().StringVarP(&file, "file", "f", "", "this is to be followed by a file location that contains a tanuki directives(commands) , it allows for more complex request")

	rootCmd.Flags().StringVarP(&method, "method", "m", "", "this specifies the method to be used for the http request")
	rootCmd.Flags().StringVarP(&header, "header", "d", ``, "this specifies the header to be used for the http request")
	rootCmd.Flags().StringVarP(&body, "body", "b", "", "this specifies the body to be used for the http request")
}
