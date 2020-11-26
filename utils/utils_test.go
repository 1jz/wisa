package utils

import (
	"testing"
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
