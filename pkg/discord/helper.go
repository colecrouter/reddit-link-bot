package discord

import (
	"context"
	"fmt"
	"time"

	"github.com/Mexican-Man/reddit-bot/pkg/video"
	"github.com/bwmarrin/discordgo"
)

func (b *DiscordBot) toDiscordMessages(URL string) (msgs []discordgo.MessageSend, nsfw bool, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	media, spoiler, nsfw, err := b.reddit.GetMedia(ctx, URL)
	if err != nil {
		return nil, false, err
	}

	const maxImagesPerMessage = 5
	const maxVideosPerMessage = 10

	var messages []discordgo.MessageSend
	var currentMessage discordgo.MessageSend

	imageCount, videoCount := 0, 0

	// Helper function to add a new message to messages slice
	addNewMessage := func() {
		if currentMessage.Content != "" || len(currentMessage.Files) > 0 {
			messages = append(messages, currentMessage)
		}
		currentMessage = discordgo.MessageSend{}
		imageCount, videoCount = 0, 0
	}

	for _, m := range media {
		// This is to help with rate limiting
		time.Sleep(time.Second)

		if m.AudioURL == "" {
			// Check if the currentMessage needs to be split due to image count
			if imageCount >= maxImagesPerMessage {
				addNewMessage()
			}

			// Newline for each non-first link or if the content is not empty
			if currentMessage.Content != "" {
				currentMessage.Content += "\n"
			}
			urlContent := m.VideoURL
			if spoiler || nsfw {
				urlContent = fmt.Sprintf("|| %s ||", urlContent)
			}
			currentMessage.Content += urlContent
			imageCount++
		} else {
			// Check if the currentMessage needs to be split due to video count
			if videoCount >= maxVideosPerMessage {
				addNewMessage()
			}

			// We can actually pass URLs directly to ffmpeg, but that requires a special
			// build of ffmpeg with HTTPS enabled. Instead, we'll download the files manually

			f2, err := video.Merge(m.Audio, m.Video)
			if err != nil {
				return nil, false, err
			}

			filename := "video.mp4"
			if spoiler {
				filename = "SPOILER_video.mp4"
			}
			currentMessage.Files = append(currentMessage.Files, &discordgo.File{
				Name: filename, Reader: f2, ContentType: "video/mp4",
			})
			videoCount++
		}
	}

	// Add the last message, if any content or files are present
	addNewMessage()

	return messages, nsfw, nil
}
