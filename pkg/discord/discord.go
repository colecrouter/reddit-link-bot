package discord

import (
	"fmt"

	"github.com/Mexican-Man/reddit-bot/pkg/scrape"
	"github.com/Mexican-Man/reddit-bot/pkg/video"
	"github.com/bwmarrin/discordgo"
)

func ToDiscordMessage(URL string) (msg discordgo.MessageSend, err error) {
	a, v, err := scrape.Scrape(URL)
	if err != nil {
		fmt.Println(err)
		return discordgo.MessageSend{}, err
	}

	if a == "" {
		msg = discordgo.MessageSend{
			Content: v,
		}
	} else {
		f2, err := video.Merge(a, v)
		if err != nil {
			fmt.Println(err)
			return discordgo.MessageSend{}, err
		}

		msg = discordgo.MessageSend{
			Files: []*discordgo.File{
				{Name: "video.mp4", Reader: *f2, ContentType: "video/mp4"},
			},
		}
	}

	return
}
