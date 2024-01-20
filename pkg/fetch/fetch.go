package fetch

import (
	"compress/gzip"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"
)

func Fetch(URL *url.URL) (*http.Response, error) {
	var resp *http.Response

	for {
		time.Sleep(time.Millisecond * 100)

		req, err := http.NewRequest("GET", URL.String(), nil)
		if err != nil {
			return nil, fmt.Errorf("unable to create GET request: %w", err)
		}

		req.Header.Set("authority", "www.reddit.com")
		req.Header.Set("method", "GET")
		req.Header.Set("path", URL.Path)

		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
		req.Header.Set("Accept-Encoding", "gzip, deflate, br")

		resp, err = http.DefaultClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("unable to GET URL: %w", err)
		}

		if resp.StatusCode != http.StatusTooManyRequests {
			break
		}

		fmt.Println("rate limited, waiting")
		<-time.After(time.Second * 30)
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("unable to GET URL: %w", os.ErrNotExist)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unable to GET URL: %s", resp.Status)
	}

	// Handle compression
	if resp.Header.Get("Content-Encoding") == "gzip" {
		var err error
		resp.Body, err = gzip.NewReader(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("unable to create gzip reader: %w", err)
		}
	}

	return resp, nil
}
