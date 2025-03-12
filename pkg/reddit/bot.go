package reddit

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/Mexican-Man/reddit-bot/pkg/config"
	"github.com/Mexican-Man/reddit-bot/pkg/video"
	"github.com/vartanbeno/go-reddit/v2/reddit"
)

type RedditBot struct {
	config config.Config
	reddit *reddit.Client
}

func NewRedditBot(config config.Config) *RedditBot {
	bot := &RedditBot{
		config: config,
	}

	var err error
	if config.Reddit == nil {
		bot.reddit, err = reddit.NewReadonlyClient()
	} else {
		bot.reddit, err = reddit.NewClient(reddit.Credentials{
			ID:       config.Reddit.ClientID,
			Secret:   config.Reddit.ClientSecret,
			Username: config.Reddit.Username,
			Password: config.Reddit.Password,
		})
	}

	if err != nil {
		panic(err)
	}

	return bot
}

func (b *RedditBot) GetMedia(ctx context.Context, link string) (media []Media, nsfw bool, spoiler bool, err error) {
	postURL, err := url.Parse(link)
	if err != nil {
		return
	}

	// Mobile links have a different ID
	// Instead of https://www.reddit.com/r/.../comments/random_id, it's https://www.reddit.com/r/.../s/different_random_id
	mobileLinkRegex := regexp.MustCompile(`^/r/[a-zA-Z0-9_]+/s/`)
	if mobileLinkRegex.MatchString(postURL.Path) {
		// Request page, get the actual ID
		response, err := b.reddit.Do(ctx, &http.Request{URL: postURL}, nil)
		if err != nil {
			return nil, false, false, fmt.Errorf("failed to create request: %w", err)
		}

		postURL = response.Response.Request.URL
	}
	postURL.RawQuery = ""

	// Extract subreddit, post ID from path
	r := regexp.MustCompile(`^/r/([a-zA-Z0-9_]+)/[a-z]+/([a-zA-Z0-9]+)(?:/.*)?$`)
	matches := r.FindStringSubmatch(postURL.Path)
	if len(matches) != 3 {
		return nil, false, false, fmt.Errorf("invalid post URL: %s", link)
	}

	// subreddit := matches[1]
	postID := matches[2]

	_, _, err = b.reddit.Post.Get(ctx, postID)
	if err != nil {
		return nil, false, false, fmt.Errorf("failed to get post %s: %w", postID, err)
	}

	// Grab media from post
	listings := []listing{}
	newURLString := strings.Replace(postURL.String()+".json", "www", "oauth", 1)
	newURL, _ := url.Parse(newURLString)
	_, err = b.reddit.Do(ctx, &http.Request{
		URL: newURL,
	}, &listings)
	if err != nil {
		return nil, false, false, fmt.Errorf("failed to get JSON response: %w", err)
	}

	// Check if kind is listing
	if listings[0].Kind != "Listing" {
		return nil, false, false, fmt.Errorf("kind is %s, not \"Listing\"", listings[0].Kind)
	}

	// Post should always be first listing
	postData := &listings[0].Data.Children[0].Data

	// The only way I know is to figure out if a post has multiple media is by checking media_metadata
	numMedia := 1
	if postData.MediaMetadata != nil {
		numMedia = len(*postData.MediaMetadata)
	}

	media = make([]Media, numMedia)

	spoiler = postData.Spoiler
	nsfw = postData.Over18

	if numMedia == 1 {
		// Default behaviour
		if postData.IsVideo {
			if postData.SecureMedia.RedditVideo.IsGif {
				media[0].VideoURL = postData.SecureMedia.RedditVideo.FallbackURL
				return
			}

			var resp *reddit.Response
			u, _ := url.Parse(postData.SecureMedia.RedditVideo.DashURL)
			resp, err = b.reddit.Do(ctx, &http.Request{
				URL: u,
			}, nil)
			if err != nil {
				return
			}
			defer resp.Body.Close()

			bytes, _ := io.ReadAll(resp.Response.Body)

			mpd := new(video.MPD)
			err = xml.Unmarshal(bytes, mpd)
			if err != nil {
				return
			}

			audio, videoLink, err := mpd.GetMediaLinks()
			if err != nil {
				return nil, false, false, err
			}

			if audio != "" {
				media[0].AudioURL, _ = url.JoinPath(postData.URLOverriddenByDest, audio)
			}
			if videoLink != "" {
				media[0].VideoURL, _ = url.JoinPath(postData.URLOverriddenByDest, videoLink)
			}

		} else {
			media[0].VideoURL = postData.URLOverriddenByDest
		}
	} else {
		// It's a gallery; order images using GalleryData if available
		if postData.GalleryData != nil && postData.MediaMetadata != nil {
			orderedMedia := make([]Media, len(postData.GalleryData.Items))
			for i, item := range postData.GalleryData.Items {
				if meta, ok := (*postData.MediaMetadata)[item.MediaID]; ok {
					orderedMedia[i].VideoURL = strings.Replace(meta.Source.URL, "&amp;", "&", -1)
				}
			}
			media = orderedMedia
		} else if postData.MediaMetadata != nil {
			// Fallback to iterating over MediaMetadata map (order may be random)
			i := 0
			for _, v := range *postData.MediaMetadata {
				media[i].VideoURL = strings.Replace(v.Source.URL, "&amp;", "&", -1)
				i++
			}
		}
	}

	// Populate audio and video ReadCloser fields
	for i := range media {
		if media[i].AudioURL != "" {
			var resp *reddit.Response
			u, _ := url.Parse(media[i].AudioURL)
			resp, err = b.reddit.Do(ctx, &http.Request{
				URL: u,
			}, nil)
			if err != nil {
				return nil, false, false, err
			}
			defer resp.Body.Close()

			media[i].Audio = resp.Response.Body
		}
		if media[i].VideoURL != "" {
			var resp *reddit.Response
			u, _ := url.Parse(media[i].VideoURL)
			resp, err = b.reddit.Do(ctx, &http.Request{
				URL: u,
			}, nil)
			if err != nil {
				return nil, false, false, err
			}
			defer resp.Body.Close()

			media[i].Video = resp.Response.Body
		}
	}

	return
}
