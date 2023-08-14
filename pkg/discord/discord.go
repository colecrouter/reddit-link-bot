package discord

import (
	"fmt"
	"time"

	"github.com/Mexican-Man/reddit-bot/pkg/scrape"
	"github.com/Mexican-Man/reddit-bot/pkg/video"
	"github.com/bwmarrin/discordgo"
)

func ToDiscordMessage(URL string) (msg discordgo.MessageSend, err error) {
	media, spoiler, err := scrape.Scrape(URL)
	if err != nil {
		return discordgo.MessageSend{}, err
	}

	// This is to help with rate limiting
	time.Sleep(time.Second)

	for _, m := range media {
		if m.AudioURL == "" {
			// Newline for each non-first link
			if msg.Content != "" {
				msg.Content += "\n"
			}

			msg.Content += m.VideoURL
		} else {
			f2, err := video.Merge(m.AudioURL, m.VideoURL)
			if err != nil {
				return discordgo.MessageSend{}, err
			}

			filename := "video.mp4"
			if spoiler {
				filename = "SPOILER_video.mp4"
			}
			msg = discordgo.MessageSend{
				Files: []*discordgo.File{
					{Name: filename, Reader: *f2, ContentType: "video/mp4"},
				},
			}
		}
	}

	if spoiler {
		msg.Content = fmt.Sprintf("|| %s ||", msg.Content)
	}

	return
}
