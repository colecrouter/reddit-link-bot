package scrape

import (
	"fmt"
	"io"
)

func readRedditResponse(r *io.ReadCloser) (b *[]byte, err error) {
	defer (*r).Close()

	// Read into byte slice
	bytes := make([]byte, 0, 102400) // 100KB
	buf := make([]byte, 32768)       // This seems to be how big the chunks from Reddit are
	total := 0
	for {
		n, err := (*r).Read(buf)
		total += n
		if err != nil && err != io.EOF {
			break
		}
		if n == 0 {
			err = nil
			break
		}
		bytes = append(bytes, buf[:n]...)
	}
	if err != nil {
		err = fmt.Errorf("unable to read response body: %w", err)
		return
	}

	trimmed := bytes[:total]
	b = &trimmed

	return
}
