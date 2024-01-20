package scrape

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/Mexican-Man/reddit-bot/pkg/fetch"
)

const TEST_URL = "https://www.reddit.com/r/pics/comments/haucpf/ive_found_a_few_funny_memories_during_lockdown/.json"

var body io.Reader
var str []byte

func init() {
	u, _ := url.Parse(TEST_URL)

	resp, err := fetch.Fetch(u)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	str, err = io.ReadAll(io.TeeReader(resp.Body, os.Stdout))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	body = strings.NewReader(string(str))
}

func BenchmarkParser(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _, _, err := parse(body)
		if err != nil {
			b.Fatal(err)
		}

		// New reader from str
		body = strings.NewReader(string(str))
	}
}
