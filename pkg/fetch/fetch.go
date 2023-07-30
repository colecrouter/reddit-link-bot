package fetch

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func Fetch(u string) (b *io.ReadCloser, err error) {
	var resp *http.Response

	for {
		// Set user agent
		var req *http.Request
		req, err = http.NewRequest("GET", u, nil)
		if err != nil {
			err = fmt.Errorf("unable to create GET request: %w", err)
			return
		}

		req.Header.Set("User-Agent", "discord-bot 1.0.0")

		resp, err = http.DefaultClient.Do(req)
		if err != nil {
			err = fmt.Errorf("unable to GET URL: %w", err)
			return
		}

		if resp.StatusCode != http.StatusTooManyRequests {
			break
		}

		fmt.Println("rate limited, waiting")
		<-time.After(time.Second * 30)
	}

	if resp.StatusCode == http.StatusNotFound {
		err = fmt.Errorf("unable to GET URL: %w", os.ErrNotExist)
	}

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("unable to GET URL: %s", resp.Status)
		return
	}

	b = &resp.Body

	return
}
