package discord

import (
	"fmt"
	"time"

	"github.com/Mexican-Man/reddit-bot/pkg/scrape"
	"github.com/Mexican-Man/reddit-bot/pkg/video"
	"github.com/bwmarrin/discordgo"
)

func ToDiscordMessages(URL string) ([]discordgo.MessageSend, error) {
	media, spoiler, err := scrape.Scrape(URL)
	if err != nil {
		return nil, err
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
			if spoiler {
				urlContent = fmt.Sprintf("|| %s ||", urlContent)
			}
			currentMessage.Content += urlContent
			imageCount++
		} else {
			// Check if the currentMessage needs to be split due to video count
			if videoCount >= maxVideosPerMessage {
				addNewMessage()
			}

			f2, err := video.Merge(m.AudioURL, m.VideoURL)
			if err != nil {
				return nil, err
			}

			filename := "video.mp4"
			if spoiler {
				filename = "SPOILER_video.mp4"
			}
			currentMessage.Files = append(currentMessage.Files, &discordgo.File{
				Name: filename, Reader: *f2, ContentType: "video/mp4",
			})
			videoCount++
		}
	}

	// Add the last message, if any content or files are present
	addNewMessage()

	return messages, nil
}
