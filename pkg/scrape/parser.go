//go:build !goexperiment.arenas

package scrape

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"

	"github.com/unki2aut/go-mpd"
)

func parse(b *[]byte) (audio string, video string, spoiler bool, err error) {
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

	if listings[0].Data.Children[0].Data.IsVideo {
		var body *io.ReadCloser
		body, err = Fetch(listings[0].Data.Children[0].Data.SecureMedia.RedditVideo.DashURL)
		if err != nil {
			return
		}
		defer (*body).Close()

		var bytes []byte
		bytes, err = io.ReadAll(*body)
		if err != nil {
			return
		}

		mpd := new(mpd.MPD)
		err = mpd.Decode(bytes)
		if err != nil {
			return
		}

		for _, adaptationSet := range mpd.Period[0].AdaptationSets {
			if *adaptationSet.ContentType == "audio" {
				audio, err = url.JoinPath(listings[0].Data.Children[0].Data.URLOverriddenByDest, adaptationSet.Representations[len(adaptationSet.Representations)-1].BaseURL[0].Value)
			} else if *adaptationSet.ContentType == "video" {
				video, err = url.JoinPath(listings[0].Data.Children[0].Data.URLOverriddenByDest, adaptationSet.Representations[len(adaptationSet.Representations)-1].BaseURL[0].Value)
			}
			if err != nil {
				return
			}
		}
	} else {
		video = listings[0].Data.Children[0].Data.URLOverriddenByDest
	}

	spoiler = listings[0].Data.Children[0].Data.Spoiler || listings[0].Data.Children[0].Data.Over18

	return
}
