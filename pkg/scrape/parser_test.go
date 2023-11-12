package scrape

import (
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/Mexican-Man/reddit-bot/pkg/fetch"
)

const TEST_URL = "https://www.reddit.com/r/pics/comments/haucpf/ive_found_a_few_funny_memories_during_lockdown/.json"

var body *io.ReadCloser
var str []byte

func init() {
	var err error
	body, err = fetch.Fetch(TEST_URL)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	str, err = io.ReadAll(io.TeeReader(*body, os.Stdout))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func BenchmarkParser(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _, _, err := parse(*body)
		if err != nil {
			b.Fatal(err)
		}
	}
}
