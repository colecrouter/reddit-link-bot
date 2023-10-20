package scrape

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/url"
	"strings"

	"github.com/Mexican-Man/reddit-bot/pkg/fetch"
	"github.com/Mexican-Man/reddit-bot/pkg/video"
)

// parse attempts to parse the JSON response from reddit, and returns the Media objects, as well as whether it needs to be a spoiler
func parse(r io.Reader) (output []Media, spoiler bool, err error) {
	// Since we know it's a listing, we can unmarshal it again into a Listing struct
	var listings []listing
	err = json.NewDecoder(r).Decode(&listings)
	if err != nil && err != io.EOF {
		return
	}
	err = nil

	// Check if kind is listing
	if listings[0].Kind != "Listing" {
		return nil, false, fmt.Errorf("kind is %s, not \"Listing\"", listings[0].Kind)
	}

	// Post should always be first listing

	// The only way I know how figure out if a post has multiple media is by checking media_metadata, which may not exist
	// If it doesn't exist, then there's only one media
	numMedia := 1
	if listings[0].Data.Children[0].Data.MediaMetadata != nil {
		numMedia = len(*listings[0].Data.Children[0].Data.MediaMetadata)
	}
	output = make([]Media, numMedia)

	spoiler = listings[0].Data.Children[0].Data.Spoiler || listings[0].Data.Children[0].Data.Over18

	if numMedia == 1 {
		// Default behaviour
		if listings[0].Data.Children[0].Data.IsVideo {
			if listings[0].Data.Children[0].Data.SecureMedia.RedditVideo.IsGif {
				output[0].VideoURL = listings[0].Data.Children[0].Data.SecureMedia.RedditVideo.FallbackURL
				return
			}

			var body *io.ReadCloser
			body, err = fetch.Fetch(listings[0].Data.Children[0].Data.SecureMedia.RedditVideo.DashURL)
			if err != nil {
				return
			}
			defer (*body).Close()

			bytes, _ := io.ReadAll(*body)

			mpd := new(video.MPD)
			err = xml.Unmarshal(bytes, mpd)
			if err != nil {
				return
			}

			audio, video, err := mpd.GetMediaLinks()
			if err != nil {
				return nil, false, err
			}

			if audio != "" {
				output[0].AudioURL, _ = url.JoinPath(listings[0].Data.Children[0].Data.URLOverriddenByDest, audio)
			}
			if video != "" {
				output[0].VideoURL, _ = url.JoinPath(listings[0].Data.Children[0].Data.URLOverriddenByDest, video)
			}

		} else {
			output[0].VideoURL = listings[0].Data.Children[0].Data.URLOverriddenByDest
		}
	} else {
		// It's a gallery (only have to worry about images)
		i := 0
		for _, v := range *listings[0].Data.Children[0].Data.MediaMetadata {
			// These ones have ampersands which will have been escaped by json.Unmarshal
			output[i].VideoURL = strings.Replace(v.Source.URL, "&amp;", "&", -1)
			i++
		}
	}

	return
}
