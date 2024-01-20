package scrape

import (
	"fmt"
	"net/url"
	"path"
	"regexp"
	"time"

	"github.com/Mexican-Man/reddit-bot/pkg/fetch"
)

var regex = regexp.MustCompile(`(?:http|https):\/\/(?:www\.)?reddit\.com\/r\/([a-zA-Z0-9_]+)\/(comments|s)\/([a-zA-Z0-9_]+)?`)

func Scrape(URL string) (media []Media, spoiler bool, nsfw bool, err error) {
	var groups []string

	for i := 0; ; i++ {
		if i == 5 {
			err = fmt.Errorf("unable to get to reddit post")
			return
		}

		groups = regex.FindAllStringSubmatch(URL, -1)[0]
		if len(groups) == 0 {
			err = fmt.Errorf("URL is not a reddit post")
			return
		}

		// New Reddit app share links are weird. It redirects to a the proper URL, so we run a single Get,
		// then it will follow all the redirects. Then, we get *that* URL and run the normal code.
		if groups[2] != "s" {
			break
		}

		u, _ := url.Parse(URL)
		resp, _ := fetch.Fetch(u)
		URL = resp.Request.URL.String()

		time.Sleep(1 * time.Second)
	}

	subreddit := groups[1]
	id := groups[3]

	u, _ := url.Parse("https://www.reddit.com/")
	u.Path = path.Join("r", subreddit, "comments", id, ".json")

	// Go to webpage, extract URL
	resp, err := fetch.Fetch(u)
	if err != nil {
		return
	}

	media, spoiler, nsfw, err = parse(resp.Body)
	if err != nil {
		return
	}

	return
}
