//go:build goexperiment.arenas

package scrape

import (
	"arena"
	"encoding/json"
	"fmt"
	"io"
	"net/url"

	"github.com/unki2aut/go-mpd"
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

	if listings[0].Data.Children[0].Data.IsVideo {
		var body *io.ReadCloser
		body, err = Fetch(listings[0].Data.Children[0].Data.SecureMedia.RedditVideo.DashURL)
		if err != nil {
			return
		}
		defer (*body).Close()

		var bytes = arena.MakeSlice[byte](a, 0, 10240)
		bytes, err = io.ReadAll(*body)
		if err != nil {
			return
		}

		// TODO idk if this is doing anything
		mpd := arena.New[mpd.MPD](a)
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
