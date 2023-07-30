package scrape

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
		resp, err = http.Get(u)
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
