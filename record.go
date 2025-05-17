package main

// Represents a record in the JSON file with test results.
type Record struct {
	URL         string `json:"url"`
	Status      string `json:"status"`
	RequestTime string `json:"requestTime"`
}

// Creates and returns a new record with test results.
func NewRecord(url string, status string, requestTime string) Record {
	return Record{
		URL:         url,
		Status:      status,
		RequestTime: requestTime,
	}
}
