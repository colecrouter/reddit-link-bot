//go:build !goexperiment.arenas

package scrape

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

func parse(b *[]byte) (audio string, video string, err error) {
	var response []thing
	err = json.Unmarshal(*b, &response)
	if err != nil {
		return
	}

	// Check if kind is listing
	if response[0].Kind != "Listing" {
		err = fmt.Errorf("kind is %s, not \"Listing\"", response[0].Kind)
		return
	}

	// Since we know it's a listing, we can unmarshal it again into a Listing struct
	var listings []listing
	err = json.Unmarshal(*b, &listings)
	if err != nil {
		return
	}

	if strings.HasPrefix(listings[0].Data.Children[0].Data.URLOverriddenByDest, "https://v.") { // v subdomain means video
		audio, _ = url.JoinPath(listings[0].Data.Children[0].Data.URLOverriddenByDest, "DASH_audio.mp4")
		video = listings[0].Data.Children[0].Data.SecureMedia.RedditVideo.FallbackURL
	} else { // Otherwise it's an image
		video = listings[0].Data.Children[0].Data.URLOverriddenByDest
	}

	return
}
