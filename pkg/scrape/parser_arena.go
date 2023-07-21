//go:build goexperiment.arenas

package scrape

import (
	"arena"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

func parse(b *[]byte) (audio string, video string, spoiler bool, err error) {
	a := arena.NewArena()
	defer a.Free()

	response := arena.MakeSlice[thing](a, 1, 1)
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
	listings := arena.MakeSlice[listing](a, 1, 1)
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

	spoiler = listings[0].Data.Children[0].Data.Spoiler || listings[0].Data.Children[0].Data.Over18

	return
}
