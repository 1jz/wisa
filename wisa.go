package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"sync"
	"time"

	"github.com/gookit/color"
	flag "github.com/spf13/pflag"
)

// flag globals
var (
	filenamePtr *string
	verbosePtr  *bool
	jsonPtr     *bool
)

// RequestResult is a struct for storing urls and status codes
type RequestResult struct {
	URL    string `json:"url"`
	Status int    `json:"status"`
}

// remove duplicate strings from a slice of strings
func removeDuplicate(urls []string) []string {
	result := make([]string, 0, len(urls))
	temp := map[string]struct{}{}
	for _, item := range urls {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

func checkLink(wg *sync.WaitGroup, url string) (RequestResult, error) {
	var r RequestResult
	var reqErr error = nil

	// defered function is run when surrounding functions are completed
	defer wg.Done()

	// use HEAD request https://golang.org/pkg/net/http/
	// resp, err := http.Head(url)
	client := http.Client{
		Timeout: 3 * time.Second,
	}
	resp, err := client.Head(url)

	if err != nil {
		if *verbosePtr && !*jsonPtr {
			fmt.Println(err)
		}
		reqErr = errors.New("request error")
	} else {
		r = RequestResult{url, resp.StatusCode}
		if !*jsonPtr {
			if resp.StatusCode == 200 {
				color.New(color.FgGreen).Printf("[GOOD] [%d] %s\n", resp.StatusCode, url)
			} else if resp.StatusCode == 400 || resp.StatusCode == 404 {
				color.New(color.FgRed).Printf("[BAD] [%d] %s\n", resp.StatusCode, url)
			} else {
				color.New(color.FgGray).Printf("[UNKNOWN] [%d] %s\n", resp.StatusCode, url)
			}
		}
	}
	return r, reqErr
}

func main() {

	// https://github.com/spf13/pflag
	filenamePtr = flag.StringP("file", "f", "", "filename input (required)") // filename input
	verbosePtr = flag.BoolP("version", "v", false, "verbose output")         // (error logs)
	jsonPtr = flag.BoolP("json", "j", false, "json output")                  // turns off verbose output

	flag.Parse()

	// regex for finding urls https://golang.org/pkg/regexp/#Regexp
	// thank you https://gist.github.com/dperini/729294 for the regex.
	// modified from ECMAscript regex, ported to work with Go regexp
	r := regexp.MustCompile(`(?:(?:(?:https?|ftp):)\/\/)(?:\S+(?::\S*)?@)?(?:(x??!(?:10|127)(?:\.\d{1,3}){3})(x??!(?:169\.254|192\.168)(?:\.\d{1,3}){2})(x??!172\.(?:1[6-9]|2\d|3[0-1])(?:\.\d{1,3}){2})(?:[1-9]\d?|1\d\d|2[01]\d|22[0-3])(?:\.(?:1?\d{1,2}|2[0-4]\d|25[0-5])){2}(?:\.(?:[1-9]\d?|1\d\d|2[0-4]\d|25[0-4]))|(?:(?:[a-z0-9\x{00a1}-\x{ffff}][a-z0-9\x{00a1}-\x{ffff}_-]{0,62})?[a-z0-9\x{00a1}-\x{ffff}]\.)+(?:[a-z\x{00a1}-\x{ffff}]{2,}\.?))(?::\d{2,5})?([\w\-\.,@?^=%&amp;:/~\+#]*[\w\-\@?^=%&amp;/~\+#])?`)
	// seperate regex for IP urls since the above did not work
	rIP := regexp.MustCompile(`(?:(?:(?:https?|ftp):)\/\/)\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}(:\d{2,5})?`)

	// check if arg is passed
	if len(*filenamePtr) == 0 {
		fmt.Println("Usage: wisa -f [file]")
		os.Exit(1)
	}

	// notify if -v flag is passed
	if *verbosePtr && !*jsonPtr {
		fmt.Println("verbose output enabled...")
	}

	// Open a file (read-only) https://golang.org/pkg/os/#Open
	if !*jsonPtr {
		fmt.Printf("Reading %s...\n", *filenamePtr)
	}

	file, err := os.Open(*filenamePtr)

	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	// read all data in to byte slice
	text, err := ioutil.ReadAll(file)

	// find all regex matches and in converted byte data and concat both string slices into single slice
	textUrls := r.FindAllString(string(text), -1)
	ipUrls := rIP.FindAllString(string(text), -1)

	urls := append(textUrls, ipUrls...)

	// stop reading file
	file.Close()

	urls = removeDuplicate(urls)

	// create workgroup to ensure all routines finish https://golang.org/pkg/sync/#WaitGroup
	var wg sync.WaitGroup

	// json output stuff
	var mut sync.Mutex
	var jsonSlice []RequestResult

	// check if urls found are alive
	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			res, err := checkLink(&wg, url)

			if err == nil {
				mut.Lock()
				jsonSlice = append(jsonSlice, res)
				mut.Unlock()
			}
		}(url)
	}

	// wait for go routines to finish
	wg.Wait()

	if *jsonPtr {
		urlsJ, err := json.Marshal(jsonSlice)

		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(string(urlsJ))
		}

	}
}
