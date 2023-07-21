package discord

import (
	"fmt"

	"github.com/Mexican-Man/reddit-bot/pkg/scrape"
	"github.com/Mexican-Man/reddit-bot/pkg/video"
	"github.com/bwmarrin/discordgo"
)

func ToDiscordMessage(URL string) (msg discordgo.MessageSend, err error) {
	a, v, spoiler, err := scrape.Scrape(URL)
	if err != nil {
		fmt.Println(err)
		return discordgo.MessageSend{}, err
	}

	if a == "" {
		content := v
		if spoiler {
			content = fmt.Sprintf("|| %s ||", v)
		}
		msg = discordgo.MessageSend{
			Content: content,
		}
	} else {
		f2, err := video.Merge(a, v)
		if err != nil {
			fmt.Println(err)
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
