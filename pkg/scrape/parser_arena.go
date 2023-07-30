//go:build goexperiment.arenas

package scrape

import (
	"arena"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/url"

	"github.com/Mexican-Man/reddit-bot/pkg/fetch"
	"github.com/Mexican-Man/reddit-bot/pkg/media"
)

func parse(b []byte) (audio string, video string, spoiler bool, err error) {
	a := arena.NewArena()
	defer a.Free()

	response := arena.MakeSlice[thing](a, 1, 1)
	err = json.Unmarshal(b, &response)
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
	err = json.Unmarshal(b, &listings)
	if err != nil {
		return
	}

	if listings[0].Data.Children[0].Data.IsVideo {
		var body *io.ReadCloser
		body, err = fetch.Fetch(listings[0].Data.Children[0].Data.SecureMedia.RedditVideo.DashURL)
		if err != nil {
			return
		}
		defer (*body).Close()

		bytes, _ := io.ReadAll(*body)

		mpd := arena.New[media.MPD](a)
		err = xml.Unmarshal(bytes, mpd)
		if err != nil {
			return
		}

		audio, video, err = mpd.GetMediaLinks()
		if err != nil {
			return
		}

		audio, _ = url.JoinPath(listings[0].Data.Children[0].Data.URLOverriddenByDest, audio)
		video, _ = url.JoinPath(listings[0].Data.Children[0].Data.URLOverriddenByDest, video)

	} else {
		video = listings[0].Data.Children[0].Data.URLOverriddenByDest
	}

	spoiler = listings[0].Data.Children[0].Data.Spoiler || listings[0].Data.Children[0].Data.Over18

	return
}
