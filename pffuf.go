package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
)

/*
this is a simple tool for parsing the output data of ffuf !!!!
*/
type ffufResult struct {
	CommandLine string   `json:"commandline"`
	Time        string   `json:"time"`
	Results     []Result `json:"results"`
	Config      Config   `json:"config"`
}

type Result struct {
	Input            Input    `json:"input"`
	Position         int      `json:"position"`
	Status           int      `json:"status"`
	Length           int      `json:"length"`
	Words            int      `json:"words"`
	Lines            int      `json:"lines"`
	ContentType      string   `json:"content-type"`
	RedirectLocation string   `json:"redirectlocation"`
	Scraper          struct{} `json:"scraper"`
	Duration         int64    `json:"duration"`
	ResultFile       string   `json:"resultfile"`
	URL              string   `json:"url"`
	Host             string   `json:"host"`
}

type Input struct {
	FFUFHash string `json:"FFUFHASH"`
	FUZZ     string `json:"FUZZ"`
}

type Config struct {
	Autocalibration         bool     `json:"autocalibration"`
	AutocalibrationKeyword  string   `json:"autocalibration_keyword"`
	AutocalibrationPerHost  bool     `json:"autocalibration_perhost"`
	AutocalibrationStrategy string   `json:"autocalibration_strategy"`
	AutocalibrationStrings  []string `json:"autocalibration_strings"`
	Colors                  bool     `json:"colors"`
	Cmdline                 string   `json:"cmdline"`
	ConfigFile              string   `json:"configfile"`
	PostData                string   `json:"postdata"`
	DebugLog                string   `json:"debuglog"`
	Delay                   struct {
		Value string `json:"value"`
	} `json:"delay"`
	DirSearchCompatibility bool              `json:"dirsearch_compatibility"`
	Extensions             []string          `json:"extensions"`
	FMode                  string            `json:"fmode"`
	FollowRedirects        bool              `json:"follow_redirects"`
	Headers                map[string]string `json:"headers"`
	IgnoreBody             bool              `json:"ignorebody"`
	IgnoreWordlistComments bool              `json:"ignore_wordlist_comments"`
	InputMode              string            `json:"inputmode"`
	CmdInputNum            int               `json:"cmd_inputnum"`
	InputProviders         []struct {
		Name     string `json:"name"`
		Keyword  string `json:"keyword"`
		Value    string `json:"value"`
		Template string `json:"template"`
	} `json:"inputproviders"`
	InputShell string `json:"inputshell"`
	JSON       bool   `json:"json"`
	Matchers   struct {
		IsCalibrated bool     `json:"IsCalibrated"`
		Mutex        struct{} `json:"Mutex"`
		Matchers     struct {
			Status struct {
				Value string `json:"value"`
			} `json:"status"`
		} `json:"Matchers"`
		Filters struct {
			Status struct {
				Value string `json:"value"`
			} `json:"status"`
		} `json:"Filters"`
		PerDomainFilters struct{} `json:"PerDomainFilters"`
	} `json:"matchers"`
	MMode               string   `json:"mmode"`
	MaxTime             int      `json:"maxtime"`
	MaxTimeJob          int      `json:"maxtime_job"`
	Method              string   `json:"method"`
	NonInteractive      bool     `json:"noninteractive"`
	OutputDirectory     string   `json:"outputdirectory"`
	OutputFile          string   `json:"outputfile"`
	OutputFormat        string   `json:"outputformat"`
	OutputSkipEmptyFile bool     `json:"OutputSkipEmptyFile"`
	ProxyURL            string   `json:"proxyurl"`
	Quiet               bool     `json:"quiet"`
	Rate                int      `json:"rate"`
	Recursion           bool     `json:"recursion"`
	RecursionDepth      int      `json:"recursion_depth"`
	RecursionStrategy   string   `json:"recursion_strategy"`
	ReplayProxyURL      string   `json:"replayproxyurl"`
	RequestFile         string   `json:"requestfile"`
	RequestProto        string   `json:"requestproto"`
	ScraperFile         string   `json:"scraperfile"`
	Scrapers            string   `json:"scrapers"`
	SNI                 string   `json:"sni"`
	Stop403             bool     `json:"stop_403"`
	StopAll             bool     `json:"stop_all"`
	StopErrors          bool     `json:"stop_errors"`
	Threads             int      `json:"threads"`
	Timeout             int      `json:"timeout"`
	URL                 string   `json:"url"`
	Verbose             bool     `json:"verbose"`
	Wordlists           []string `json:"wordlists"`
	HTTP2               bool     `json:"http2"`
}

func checkErr(e error) {
	if e != nil {
		log.Fatal(e)
	}

}

// checks if a slice is empty or not
func empty(slice []string) bool {
	if len(slice) > 0 {
		return false
	}
	return true
}

// gets a path and creates a file and add the elements of the slice into it
func writeInto(path string, slice []string) {

	f, err := os.Create(path)
	checkErr(err)

	defer f.Close()

	for _, e := range slice {
		_, err := f.WriteString(e + "\n")
		checkErr(err)
	}

}

// catching error in case some directories exist
func RecoverFromPanic() {
	if r := recover(); r != nil {
		fmt.Println("some directories exist!")
	}
}

func createDir() {
	defer RecoverFromPanic()

	comm := exec.Command("mkdir", "302", "403", "500")
	err := comm.Run()
	if err != nil {
		panic("something went wrong")
	}

}

func main() {
	var file_name string
	flag.StringVar(&file_name, "f", "", "path of the ffuf result")
	flag.Parse()

	data, err := os.ReadFile(file_name)
	checkErr(err)

	var ffuf ffufResult
	err = json.Unmarshal(data, &ffuf)
	checkErr(err)

	status := make([]int, len(ffuf.Results))
	urls := make([]string, len(ffuf.Results))

	// three slices for handling endpoints with different status code
	redirects := make([]string, 0)
	not_allowed := make([]string, 0)
	server_error := make([]string, 0)

	// create 3 directory in the path where the program is running
	createDir()

	// gets the endpoints and it's status code
	for i, result := range ffuf.Results {
		urls[i] = result.URL
		status[i] = result.Status
	}

	for i, u := range urls {

		if status[i] == 301 || status[i] == 302 {
			redirects = append(redirects, u)
		} else if status[i] == 401 || status[i] == 403 {
			not_allowed = append(not_allowed, u)
		} else if status[i] == 501 || status[i] == 502 || status[i] == 500 {
			server_error = append(server_error, u)
		}
	}
	if !empty(redirects) {
		writeInto("./302/endpoints.txt", redirects)
	}
	if !empty(not_allowed) {
		writeInto("./403/endpoints.txt", not_allowed)
	}
	if !empty(server_error) {
		writeInto("./500/endpoints.txt", server_error)
	}
}
