package reddit

import "io"

// Media represents the media components of a Reddit post.
type Media struct {
	AudioURL string
	VideoURL string
	Audio    io.ReadCloser
	Video    io.ReadCloser
}

type listing struct {
	ID   string `json:"id"`
	Kind string `json:"kind"`
	Data struct {
		Before   string `json:"before"`
		After    string `json:"after"`
		Children []struct {
			Data struct {
				Subreddit           string         `json:"subreddit"`
				Gilded              int            `json:"gilded"`
				Clicked             bool           `json:"clicked"`
				Title               string         `json:"title"`
				Hidden              bool           `json:"hidden"`
				Downs               int            `json:"downs"`
				Ups                 int            `json:"ups"`
				TotalAwardsReceived int            `json:"total_awards_received"`
				Score               int            `json:"score"`
				Thumbnail           string         `json:"thumbnail"`
				Created             float64        `json:"created"`
				URLOverriddenByDest string         `json:"url_overridden_by_dest"`
				Archived            bool           `json:"archived"`
				Over18              bool           `json:"over_18"`
				Spoiler             bool           `json:"spoiler"`
				Locked              bool           `json:"locked"`
				SubredditID         string         `json:"subreddit_id"`
				ID                  string         `json:"id"`
				IsRobotIndexable    bool           `json:"is_robot_indexable"`
				Permalink           string         `json:"permalink"`
				URL                 string         `json:"url"`
				IsVideo             bool           `json:"is_video"`
				SecureMedia         *secureMedia   `json:"secure_media"`
				MediaMetadata       *mediaMetadata `json:"media_metadata"`
				GalleryData         *GalleryData   `json:"gallery_data,omitempty"`
			} `json:"data"`
		} `json:"children"`
	} `json:"data"`
}

// secureMedia holds the details for Reddit videos.
type secureMedia struct {
	RedditVideo struct {
		BitrateKbps       int    `json:"bitrate_kbps"`
		FallbackURL       string `json:"fallback_url"`
		Height            int    `json:"height"`
		Width             int    `json:"width"`
		ScrubberMediaURL  string `json:"scrubber_media_url"`
		DashURL           string `json:"dash_url"`
		Duration          int    `json:"duration"`
		HlsURL            string `json:"hls_url"`
		IsGif             bool   `json:"is_gif"`
		TranscodingStatus string `json:"transcoding_status"`
	} `json:"reddit_video"`
}

// mediaMetadata represents a mapping of media IDs to their metadata.
type mediaMetadata map[string]struct {
	Source struct {
		Width  int    `json:"x"`
		Height int    `json:"y"`
		URL    string `json:"u"`
	} `json:"s"`
}

// GalleryData holds the ordering information for gallery posts.
type GalleryData struct {
	Items []GalleryItem `json:"items"`
}

// GalleryItem defines a single gallery item with its media id and display order.
type GalleryItem struct {
	MediaID string `json:"media_id"`
	Order   int    `json:"order"`
}
