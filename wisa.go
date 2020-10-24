package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"sync"

	"github.com/1jz/WISA/utils"
	flag "github.com/spf13/pflag"
)

func main() {

	// https://github.com/spf13/pflag
	utils.FilenamePtr = flag.StringP("file", "f", "", "filename input (required)") // filename input
	utils.VerbosePtr = flag.BoolP("version", "v", false, "verbose output")         // (error logs)
	utils.JSONPtr = flag.BoolP("json", "j", false, "json output")                  // turns off verbose output
	utils.IgnoreFilePtr = flag.BoolP("ignore", "i", false, "ignores certain urls based off a text file")

	flag.Parse()

	// regex for finding urls https://golang.org/pkg/regexp/#Regexp
	// thank you https://gist.github.com/dperini/729294 for the regex.
	// modified from ECMAscript regex, ported to work with Go regexp
	urlRegex := regexp.MustCompile(`(?:(?:(?:https?|ftp):)\/\/)(?:\S+(?::\S*)?@)?(?:(x??!(?:10|127)(?:\.\d{1,3}){3})(x??!(?:169\.254|192\.168)(?:\.\d{1,3}){2})(x??!172\.(?:1[6-9]|2\d|3[0-1])(?:\.\d{1,3}){2})(?:[1-9]\d?|1\d\d|2[01]\d|22[0-3])(?:\.(?:1?\d{1,2}|2[0-4]\d|25[0-5])){2}(?:\.(?:[1-9]\d?|1\d\d|2[0-4]\d|25[0-4]))|(?:(?:[a-z0-9\x{00a1}-\x{ffff}][a-z0-9\x{00a1}-\x{ffff}_-]{0,62})?[a-z0-9\x{00a1}-\x{ffff}]\.)+(?:[a-z\x{00a1}-\x{ffff}]{2,}\.?))(?::\d{2,5})?([\w\-\.,@?^=%&amp;:/~\+#]*[\w\-\@?^=%&amp;/~\+#])?`)
	// seperate regex for IP urls since the above did not work
	ipRegex := regexp.MustCompile(`(?:(?:(?:https?|ftp):)\/\/)\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}(:\d{2,5})?`)

	//String array holding the list of urls to ignore
	ignoreURLs := []string{}

	// check if arg is passed
	if len(os.Args) == 1 {
		fmt.Println("Usage: wisa -f [file]")
		os.Exit(1)
	} else {
		//Ignore
		if *utils.IgnoreFilePtr {
			//Get the file path
			ignoreFilePath := flag.Args()[0]
			ignoreURLs = utils.GetIgnorePatterns(ignoreFilePath, urlRegex, ipRegex)
			//Remove any dupes
			ignoreURLs = utils.RemoveDuplicate(ignoreURLs)

		}

		// notify if -v flag is passed
		if *utils.VerbosePtr && !*utils.JSONPtr {
			fmt.Println("verbose output enabled...")
		}

		// Open a file (read-only) https://golang.org/pkg/os/#Open
		if !*utils.JSONPtr {
			fmt.Printf("Reading %s...\n", *utils.FilenamePtr)
		}

		file, err := os.Open(*utils.FilenamePtr)

		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}

		// read all data in to byte slice
		text, err := ioutil.ReadAll(file)

		// find all regex matches and in converted byte data and concat both string slices into single slice
		textUrls := urlRegex.FindAllString(string(text), -1)
		ipUrls := ipRegex.FindAllString(string(text), -1)

		urls := append(textUrls, ipUrls...)

		// stop reading file
		file.Close()

		urls = utils.RemoveDuplicate(urls)

		//Call the ignoreURLs func to remove any urls that match the list
		if *utils.IgnoreFilePtr == true {
			urls = utils.IgnoreURL(urls, ignoreURLs)
		}

		// create workgroup to ensure all routines finish https://golang.org/pkg/sync/#WaitGroup
		var wg sync.WaitGroup

		// json output stuff
		var jsonSlice []utils.RequestResult
		var mut sync.Mutex
		finalExit := 0

		// check if urls found are alive
		for _, url := range urls {
			wg.Add(1)
			go func(url string) {
				res, status, err := utils.CheckLink(&wg, url)

				if err == nil {
					mut.Lock()
					jsonSlice = append(jsonSlice, res)
					mut.Unlock()
				}

				if status == 3 {
					mut.Lock()
					finalExit = status
					mut.Unlock()
				}
			}(url)
		}

		// wait for go routines to finish
		wg.Wait()

		if *utils.JSONPtr {
			urlsJ, err := json.Marshal(jsonSlice)

			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(string(urlsJ))
			}

		}
		os.Exit(finalExit)
	}
}
