package indexer

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var ipLocationCache sync.Map

type IPLocation struct {
	City     string `json:"city"`
	Region   string `json:"region"`
	Country  string `json:"country"`
	Location string `json:"loc"`
	Timezone string `json:"timezone"`
}

func QueryLocation(ip string) (*IPLocation, error) {
	if loc, ok := ipLocationCache.Load(ip); ok {
		return loc.(*IPLocation), nil
	}

	url := fmt.Sprintf("http://ipinfo.io/%v/json", ip)
	resp, err := http.Get(url)
	if err != nil {
		return nil, errors.WithMessagef(err, "Failed to http GET from %v", url)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.WithMessage(err, "Failed to read http response body")
	}

	var location IPLocation
	if err = json.Unmarshal(body, &location); err != nil {
		return nil, errors.WithMessage(err, "Failed to unmarshal Location data")
	}

	ipLocationCache.Store(ip, &location)

	return &location, nil
}

func StartIPLocationCache(filename string, persistInterval time.Duration) {
	n, err := readIPLocationCache(filename)
	if err != nil {
		logrus.WithError(err).Warn("Failed to read IP location cache")
	} else {
		logrus.WithField("ips", n).Info("Succeeded to read IP location cache")
	}

	go writeIPLocationCache(filename, persistInterval)
}

func readIPLocationCache(filename string) (int, error) {
	data, err := os.ReadFile(filename)
	if os.IsNotExist(err) {
		return 0, nil
	}

	if err != nil {
		return 0, errors.WithMessagef(err, "Failed to read file %v", filename)
	}

	var ips map[string]*IPLocation
	if err = json.Unmarshal(data, &ips); err != nil {
		return 0, errors.WithMessage(err, "Failed to unmarshal data")
	}

	for ip, loc := range ips {
		ipLocationCache.Store(ip, loc)
	}

	return len(ips), nil
}

func writeIPLocationCache(filename string, persistInterval time.Duration) {
	ticker := time.NewTicker(persistInterval)
	defer ticker.Stop()

	for range ticker.C {
		ips := map[string]*IPLocation{}

		ipLocationCache.Range(func(key, value any) bool {
			ips[key.(string)] = value.(*IPLocation)
			return true
		})

		if len(ips) == 0 {
			continue
		}

		data, err := json.MarshalIndent(ips, "", "    ")
		if err != nil {
			logrus.WithError(err).Warn("Failed to marshal ip locations")
			continue
		}

		if err = os.WriteFile(filename, data, os.ModePerm); err != nil {
			logrus.WithError(err).Warn("Failed to write ip locations to cache file")
			continue
		}

		logrus.WithField("ips", len(ips)).Info("Succeeded to cache IP locations")
	}
}
