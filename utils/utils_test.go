package utils

import (
	"regexp"
	"testing"

	flag "github.com/spf13/pflag"
)

func TestRemoveDuplicate(t *testing.T) {
	data := []string{"http://abc.com", "http://abc.com", "http://abc.ca"}
	expectedData := []string{"http://abc.com", "http://abc.ca"}
	duplicateRemovedData := RemoveDuplicate(data)

	if len(duplicateRemovedData) != 2 {
		t.Errorf("RemoveDuplicate was incorrect, received %d, want %d", len(duplicateRemovedData), 2)
	}

	for i, link := range duplicateRemovedData {
		if link != expectedData[i] {
			t.Errorf("RemoveDuplicate was incorrect, received %s, want %s", link, expectedData[i])
		}
	}
}

func TestCheckLink(t *testing.T) {

	JSONPtr = flag.BoolP("json", "j", false, "json output")
	VerbosePtr = flag.BoolP("version", "v", true, "verbose output")

	_, result1, _ := CheckLink("https://httpstat.us/200")
	_, result2, _ := CheckLink("https://httpstat.us/404")
	_, result3, _ := CheckLink("https://httpstat.us/403")
	_, result4, _ := CheckLink("---")

	if result1 != 1 {
		t.Errorf("CheckLink was incorrect, received: %d, want: %d.", result1, 1)
	}

	if result2 != 2 {
		t.Errorf("CheckLink was incorrect, received: %d, want: %d.", result2, 2)
	}

	if result3 != 3 {
		t.Errorf("CheckLink was incorrect, received: %d, want: %d.", result3, 3)
	}

	if result4 != -1 {
		t.Errorf("CheckLink was incorrect, received: %d, want: %d.", result4, -1)
	}
}

func TestIgnoreURL(t *testing.T) {
	data := []string{"http://abc.com", "http://abc.ca", "http://abc.net"}
	ignoreData := []string{"http://abc.com", "http://abc.ca"}
	expectedData := []string{"", "", "http://abc.net"}

	validUrls := IgnoreURL(data, ignoreData)

	for i, link := range validUrls {
		if link != expectedData[i] {
			t.Errorf("IgnoreURL was incorrect, received %s, want %s", link, expectedData[i])
		}
	}
}

func TestGetIgnorePatterns(t *testing.T) {
	urlRegex := regexp.MustCompile(`(?:(?:(?:https?|ftp):)\/\/)(?:\S+(?::\S*)?@)?(?:(x??!(?:10|127)(?:\.\d{1,3}){3})(x??!(?:169\.254|192\.168)(?:\.\d{1,3}){2})(x??!172\.(?:1[6-9]|2\d|3[0-1])(?:\.\d{1,3}){2})(?:[1-9]\d?|1\d\d|2[01]\d|22[0-3])(?:\.(?:1?\d{1,2}|2[0-4]\d|25[0-5])){2}(?:\.(?:[1-9]\d?|1\d\d|2[0-4]\d|25[0-4]))|(?:(?:[a-z0-9\x{00a1}-\x{ffff}][a-z0-9\x{00a1}-\x{ffff}_-]{0,62})?[a-z0-9\x{00a1}-\x{ffff}]\.)+(?:[a-z\x{00a1}-\x{ffff}]{2,}\.?))(?::\d{2,5})?([\w\-\.,@?^=%&amp;:/~\+#]*[\w\-\@?^=%&amp;/~\+#])?`)
	ipRegex := regexp.MustCompile(`(?:(?:(?:https?|ftp):)\/\/)\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}(:\d{2,5})?`)
	urls := GetIgnorePatterns("../test_ignore", urlRegex, ipRegex)

	expectedData := []string{"https://www.google.com"}

	for i, link := range urls {
		if link != expectedData[i] {
			t.Errorf("GetIgnorePatterns was incorrect, received %s, want %s", link, expectedData[i])
		}
	}
}
