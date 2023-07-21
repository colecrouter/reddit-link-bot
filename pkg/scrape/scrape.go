package scrape

import (
	"fmt"
	"net/url"
	"regexp"
)

var regex = regexp.MustCompile(`(http|https)(:\/\/)(www\.?)(reddit\.com\/r\/)([a-zA-Z0-9_]+)(\/comments\/)([a-zA-Z0-9_]+)(\/)([a-zA-Z0-9_]+)(\/)`)

func Scrape(URL string) (audioURL string, videoURL string, err error) {
	groups := regex.FindAllStringSubmatch(URL, -1)
	if len(groups[0]) != 11 {
		err = fmt.Errorf("URL is not a reddit post")
		return
	}

	u, err := url.Parse(URL)
	if err != nil {
		fmt.Printf("unable to parse URL: %v", err)
		return
	}

	u.RawQuery = ""
	u.Path += ".json"

	// Go to webpage, extract URL
	body, err := Fetch(u.String())
	if err != nil {
		return
	}

	str, err := read(body)
	if err != nil {
		return
	}

	audioURL, videoURL, err = parse(str)
	if err != nil {
		return
	}

	return
}
