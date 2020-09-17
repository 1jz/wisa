package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"sync"
)

func checkLink(wg *sync.WaitGroup, url string) {

	defer wg.Done()

	// use HEAD request https://golang.org/pkg/net/http/
	resp, err := http.Head(url)

	if err != nil {
		fmt.Println(err)
	} else {
		// Status codes https://developer.mozilla.org/en-US/docs/Web/HTTP/Status
		if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
			fmt.Printf("[PASS] [%d] %s\n", resp.StatusCode, url)
		} else if resp.StatusCode == 403 {
			fmt.Printf("[WARN] [%d] %s\n", resp.StatusCode, url)
		} else {
			fmt.Printf("[DEAD] [%d] %s\n", resp.StatusCode, url)
		}
	}
}

func main() {
	// regex for finding urls
	// modification of this ECMAscript regex ported to work with go regexp
	r := regexp.MustCompile(`(?:(?:(?:https?|ftp):)\/\/)(?:\S+(?::\S*)?@)?(?:(x??!(?:10|127)(?:\.\d{1,3}){3})(x??!(?:169\.254|192\.168)(?:\.\d{1,3}){2})(x??!172\.(?:1[6-9]|2\d|3[0-1])(?:\.\d{1,3}){2})(?:[1-9]\d?|1\d\d|2[01]\d|22[0-3])(?:\.(?:1?\d{1,2}|2[0-4]\d|25[0-5])){2}(?:\.(?:[1-9]\d?|1\d\d|2[0-4]\d|25[0-4]))|(?:(?:[a-z0-9\x{00a1}-\x{ffff}][a-z0-9\x{00a1}-\x{ffff}_-]{0,62})?[a-z0-9\x{00a1}-\x{ffff}]\.)+(?:[a-z\x{00a1}-\x{ffff}]{2,}\.?))(?::\d{2,5})?([\w\-\.,@?^=%&amp;:/~\+#]*[\w\-\@?^=%&amp;/~\+#])?`)
	// seperate regex for IP urls since the above did not work
	r_ip := regexp.MustCompile(`(?:(?:(?:https?|ftp):)\/\/)\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}(:\d{2,5})?`)

	// check if arg is passed
	if len(os.Args) < 2 {
		fmt.Println("Usage: wisa [file]")
		os.Exit(1)
	} else if len(os.Args) > 2 {
		fmt.Println("Please enter a single filename.")
		os.Exit(2)
	}

	// Open a file (read-only) https://golang.org/pkg/os/#Open
	file, err := os.Open(os.Args[1])

	if err != nil {
		fmt.Println(err)
		os.Exit(3)
	}

	// read all data in to byte slice
	text, err := ioutil.ReadAll(file)

	// find all regex matches and in converted byte data and concat both string slices into single slice
	text_urls := r.FindAllString(string(text), -1)
	ip_urls := r_ip.FindAllString(string(text), -1)

	urls := append(text_urls, ip_urls...)

	// stop reading file
	file.Close()

	// create workgroup to ensure all routines finish https://golang.org/pkg/sync/#WaitGroup
	var wg sync.WaitGroup

	// check if urls found are alive
	for _, url := range urls {
		wg.Add(1)
		go checkLink(&wg, url)
	}

	wg.Wait()
}
