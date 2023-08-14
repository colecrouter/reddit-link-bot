package scrape

import (
	"fmt"
	"net/http"
	"net/url"
	"path"
	"regexp"

	"github.com/Mexican-Man/reddit-bot/pkg/fetch"
)

var regex = regexp.MustCompile(`(?:http|https):\/\/(?:www\.)?reddit\.com\/r\/([a-zA-Z0-9_]+)\/(comments|s)\/([a-zA-Z0-9_]+)?`)

func Scrape(URL string) (media []Media, spoiler bool, err error) {
	groups := regex.FindAllStringSubmatch(URL, -1)
	if len(groups[0]) == 0 {
		err = fmt.Errorf("URL is not a reddit post")
		return
	}

	// New Reddit app share links are weird. It redirects to a the proper URL, so we run a single Get,
	// then it will follow all the redirects. Then, we get *that* URL and run the normal code.
	if groups[0][2] == "s" {
		resp, _ := http.Get(URL)
		return Scrape(resp.Request.URL.String())
	}

	subreddit := groups[0][1]
	id := groups[0][3]

	u, _ := url.Parse("https://www.reddit.com/")
	u.Path = path.Join("r", subreddit, "comments", id, ".json")

	// Go to webpage, extract URL
	body, err := fetch.Fetch(u.String())
	if err != nil {
		return
	}

	media, spoiler, err = parse(*body)
	if err != nil {
		return
	}

	return
}
