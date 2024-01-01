/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"strconv"
	"github.com/msalbrain/turbotanuki/pkg"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	// "log"
	"net/url"
	"os"

	// ui "github.com/gizak/termui/v3"
	// "github.com/gizak/termui/v3/widgets"
)

var URL string

var numRequest int64
var connRequest int64

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
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {

		// for e := range ui.PollEvents() {
		// 	if e.Type == ui.KeyboardEvent {
		// 		break
		// 	}
		// }

		// table.Render()

		if (numRequest < connRequest) {
			fmt.Println("\nyour maths is out of order,\nthere are far more number of conncurrent request than number of request to be made\n ")
			return
		}

		if len(args) < 1 && URL == "" {
			fmt.Println("no url provided. check `turbotanuki --help`")
		} else if len(args) == 1 {
			_, err := url.Parse(args[0])
			if err != nil {
				fmt.Println("Invalid URL")
				return
			}
			err = pkg.MakePlainHttpCall(args[0])
			if err != nil {
				fmt.Println("network problem")
				return
			}

		} else {
			if URL != "" {
				// var URL string
				// var numRequest int64
				// var connRequest int64
				// var method string
				// var header string
				// var body string
				// var file string

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
					Url: URL,
					Method: method,
					Body: []byte(body),
					Header: dataHeader,
				}

				res := pkg.MultipleRequest(p, numRequest, connRequest)
				
				metrics := [][]string{
					{"Total time taken (S)", strconv.FormatFloat(res.TotalTime, 'f', -1, 64)},
					{"Requests Per Second (RPS)", strconv.FormatFloat(float64(res.RequestPerSec), 'f', -1, 64)},
				}

				table := tablewriter.NewWriter(os.Stdout)

				table.SetHeader([]string{"Metric", "Value"})

				for _, v := range metrics {
					table.Append(v)
				}

				table.Render()

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

	rootCmd.Flags().StringVarP(&file, "file", "f", "", "this is to be followed by a file location that contains a tanuki directives(commands) , it allows for more complex request")

	// var method string
	// var header string
	// var body string

	rootCmd.Flags().StringVarP(&method, "method", "m", "", "this specifies the method to be used for the http request")
	rootCmd.Flags().StringVarP(&header, "header", "d", ``, "this specifies the header to be used for the http request")
	rootCmd.Flags().StringVarP(&body, "body", "b", "", "this specifies the body to be used for the http request")

	// rootCmd.MarkFlagRequired("name")
}
