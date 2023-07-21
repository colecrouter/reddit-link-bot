package scrape

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"
)

const TEST_URL = "https://www.reddit.com/r/pics/comments/haucpf/ive_found_a_few_funny_memories_during_lockdown/.json"

var body *io.ReadCloser
var str *[]byte

func init() {
	var err error
	body, err = Fetch(TEST_URL)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	str, err = read(body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func BenchmarkParser(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _, _, err := parse(str)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkReader(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := read(body)
		if err != nil {
			b.Fatal(err)
		}

		*body = io.NopCloser(bytes.NewReader(*str))
	}

}
