package utils

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/gookit/color"
)

// RequestResult is a struct for storing urls and status codes
type RequestResult struct {
	URL    string `json:"url"`
	Status int    `json:"status"`
}

// flag globals
var (
	FilenamePtr   *string
	VerbosePtr    *bool
	JSONPtr       *bool
	IgnoreFilePtr *bool
)

// RemoveDuplicate removes duplicate strings from a slice of strings
func RemoveDuplicate(urls []string) []string {
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

// CheckLink makes a request to a URL passed and returns it's status
func CheckLink(wg *sync.WaitGroup, url string) (RequestResult, int, error) {
	var r RequestResult
	var reqErr error = nil
	var status int

	// defered function is run when surrounding functions are completed
	defer wg.Done()

	// use HEAD request https://golang.org/pkg/net/http/
	// resp, err := http.Head(url)
	client := http.Client{
		Timeout: 3 * time.Second,
	}
	resp, err := client.Head(url)

	if err != nil {
		if *VerbosePtr && !*JSONPtr {
			fmt.Println(err)
		}
		reqErr = errors.New("request error")
	} else {
		r = RequestResult{url, resp.StatusCode}
		if !*JSONPtr {
			if resp.StatusCode == 200 {
				status = 0
				color.New(color.FgGreen).Printf("[GOOD] [%d] %s\n", resp.StatusCode, url)
			} else if resp.StatusCode == 400 || resp.StatusCode == 404 {
				status = 3
				color.New(color.FgRed).Printf("[BAD] [%d] %s\n", resp.StatusCode, url)
			} else {
				status = 3
				color.New(color.FgGray).Printf("[UNKNOWN] [%d] %s\n", resp.StatusCode, url)
			}
		}
	}
	return r, status, reqErr
}

// IgnoreURL goes through the url list and removes any url that matches with the ignore url list
func IgnoreURL(urls []string, ignoreList []string) []string {
	//Loop through the url list
	for i := 0; i < len(urls); i++ {
		//Loop through the ignore url list
		for k := 0; k < len(ignoreList); k++ {
			//If the beginning of the url matches the ignoreList value, set the url value to ""
			if strings.HasPrefix(urls[i], ignoreList[k]) == true {
				urls[i] = ""
			}
		}
	}
	return urls
}

// GetIgnorePatterns creates a slice of URLs to ignore from the file string passed
func GetIgnorePatterns(ignoreFilePath string, urlRegex *regexp.Regexp, ipRegex *regexp.Regexp) []string {
	ignoreURLs := []string{}

	fmt.Printf("Reading %s...\n", ignoreFilePath)
	//Open the file
	fileIgn, err := os.Open(ignoreFilePath)
	//In case of panic
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	//Close the file
	defer fileIgn.Close()
	//Create new scanner
	scanner := bufio.NewScanner(fileIgn)
	//Loop through each line of the file
	for scanner.Scan() {
		line := scanner.Text()
		//If the starting character is not "#"
		startChar := string(line[0])
		if startChar != "#" {
			//Match agansit regrex patterns
			UrlsFromLine := urlRegex.FindAllString(string(line), -1)
			ipUrlsFromLine := ipRegex.FindAllString(string(line), -1)
			UrlsFromLine = append(UrlsFromLine, ipUrlsFromLine...)
			ignoreURLs = append(ignoreURLs, UrlsFromLine...)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return ignoreURLs
}
