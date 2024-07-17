package indexer

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

type IPLocation struct {
	City     string `json:"city"`
	Region   string `json:"region"`
	Country  string `json:"country"`
	Location string `json:"loc"`
	Timezone string `json:"timezone"`
}

func QueryLocation(ip string) (IPLocation, error) {
	url := fmt.Sprintf("http://ipinfo.io/%v/json", ip)
	resp, err := http.Get(url)
	if err != nil {
		return IPLocation{}, errors.WithMessagef(err, "Failed to http GET from %v", url)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return IPLocation{}, errors.WithMessage(err, "Failed to read http response body")
	}

	var location IPLocation
	if err = json.Unmarshal(body, &location); err != nil {
		return IPLocation{}, errors.WithMessage(err, "Failed to unmarshal Location data")
	}

	return location, nil
}
