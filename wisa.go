package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"sync"

	"github.com/gookit/color"
)

// flag globals
var (
	filenamePtr *string
	verbosePtr  *bool
)

func checkLink(wg *sync.WaitGroup, url string) {

	// defered function is run when surrounding functions are completed
	defer wg.Done()

	// use HEAD request https://golang.org/pkg/net/http/
	resp, err := http.Head(url)

	if err != nil {
		if *verbosePtr {
			fmt.Println(err)
		}
	} else {
		// Status codes https://developer.mozilla.org/en-US/docs/Web/HTTP/Status
		if resp.StatusCode == 200 {
			color.New(color.FgGreen).Printf("[GOOD] [%d] %s\n", resp.StatusCode, url)
		} else if resp.StatusCode == 400 || resp.StatusCode == 404 {
			color.New(color.FgRed).Printf("[BAD] [%d] %s\n", resp.StatusCode, url)
		} else {
			color.New(color.FgGray).Printf("[UNKNOWN] [%d] %s\n", resp.StatusCode, url)
		}
	}
}

func main() {

	// https://golang.org/pkg/flag/
	filenamePtr = flag.String("f", "", "filename input (required)") // filename input
	verbosePtr = flag.Bool("v", false, "verbose output")            // (error logs)

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

	// notift if -v flag is passed
	if *verbosePtr {
		fmt.Println("verbose output enabled...")
	}

	// Open a file (read-only) https://golang.org/pkg/os/#Open
	fmt.Printf("Reading %s...\n", *filenamePtr)
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

	// create workgroup to ensure all routines finish https://golang.org/pkg/sync/#WaitGroup
	var wg sync.WaitGroup

	// check if urls found are alive
	for _, url := range urls {
		wg.Add(1)
		go checkLink(&wg, url)
	}

	// wait for go routines to finish
	wg.Wait()
}
