package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/Mexican-Man/reddit-bot/pkg/config"
	"github.com/Mexican-Man/reddit-bot/pkg/discord"
	"github.com/bwmarrin/discordgo"
)

func main() {
	// Load config
	cfg := config.Config{}
	err := cfg.Load()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Start bot
	discord, err := discordgo.New("Bot " + cfg.DiscordToken)
	if err != nil {
		fmt.Println(err)
	}

	// Add handlers
	discord.AddHandler(messageCreate)
	discord.AddHandler(ready)

	// Set intents
	discord.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages)

	// Connect
	err = discord.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Stop from closing
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	discord.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	// Ignore messages from bots
	if m.Author.Bot {
		return
	}

	if len(m.Content) < 23 || m.Content[:23] != "https://www.reddit.com/" {
		return
	}

	s.ChannelTyping(m.ChannelID)

	newM, err := discord.ToDiscordMessage(m.Content)
	if err != nil {
		goto ERROR
	}

	_, err = s.ChannelMessageSendComplex(m.ChannelID, &newM)
	if err != nil {
		if strings.HasPrefix(err.Error(), "HTTP 413") {
			s.MessageReactionAdd(m.ChannelID, m.ID, "🥵")
			return
		}

		goto ERROR
	}

	return

ERROR:
	s.MessageReactionAdd(m.ChannelID, m.ID, "😵")
	fmt.Printf("%+v\n", err)
}

func ready(s *discordgo.Session, e *discordgo.Ready) {
	// Set status
	act := []*discordgo.Activity{
		{
			Name: "your browser history",
			Type: discordgo.ActivityTypeWatching,
		},
	}
	idle := 0
	s.UpdateStatusComplex(discordgo.UpdateStatusData{IdleSince: &idle, Activities: act, AFK: false})
}
