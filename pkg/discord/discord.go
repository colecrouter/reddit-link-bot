package discord

import (
	"fmt"
	"time"

	"github.com/Mexican-Man/reddit-bot/pkg/media"
	"github.com/Mexican-Man/reddit-bot/pkg/scrape"
	"github.com/bwmarrin/discordgo"
)

func ToDiscordMessage(URL string) (msg discordgo.MessageSend, err error) {
	a, v, spoiler, err := scrape.Scrape(URL)
	if err != nil {
		return discordgo.MessageSend{}, err
	}

	// This is to help with rate limiting
	time.Sleep(time.Second)

	if a == "" {
		content := v
		if spoiler {
			content = fmt.Sprintf("|| %s ||", v)
		}
		msg = discordgo.MessageSend{
			Content: content,
		}
	} else {
		f2, err := media.Merge(a, v)
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

	return
}
