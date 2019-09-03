// the feed for data.json
package search

import (
	"encoding/json"
	"os"
)

const dataFile = "data/data.json"

// Feed contains information we need to process a feed
type Feed struct {
	Name string `json:"site"`
	URI  string `json:"link"`
	Type string `json:"type"`
}

// RetrieveFeeds reads and unmarshals the feed data file
func RetrieveFeeds() ([]*Feed, error) {
	// open the file
	file, err := os.Open(dataFile)
	if err != nil {
		return nil, err
	}

	// schedule the file to be closed once
	defer file.Close()

	// decode the file into a slice of pointers
	// to feed values
	var feeds []*Feed
	err = json.NewDecoder(file).Decode(&feeds)

	return feeds, err
}
