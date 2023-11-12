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

var cfg config.Config

func main() {

	// Load config
	cfg = config.Config{}
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
		fmt.Printf("error opening connection: %+v\n", err)
		return
	}

	defer discord.Close()

	// Stop from closing
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	var messages []discordgo.MessageSend
	var nsfw bool
	var err error

	// Ignore messages from bots
	if m.Author.Bot {
		return
	}

	// Check roles
	if len(cfg.Roles) > 0 {
		member, err := s.GuildMember(m.GuildID, m.Author.ID)
		if err != nil {
			goto ERROR
		}

		hasRole := false
		for _, role := range cfg.Roles {
			for _, memberRole := range member.Roles {
				if role == memberRole {
					hasRole = true
					break
				}
			}
		}

		if !hasRole {
			return
		}
	}

	// Check channels
	if len(cfg.Channels) > 0 {
		hasChannel := false
		for _, channel := range cfg.Channels {
			if channel == m.ChannelID {
				hasChannel = true
				break
			}
		}

		if !hasChannel {
			return
		}
	}

	if (len(m.Content) < 23 || m.Content[:23] != "https://www.reddit.com/") && (len(m.Content) < 19 || m.Content[:19] != "https://reddit.com/") {
		return
	}

	s.ChannelTyping(m.ChannelID)

	messages, nsfw, err = discord.ToDiscordMessages(m.Content)
	if err != nil {
		goto ERROR
	}

	// Check if the message is NSFW
	if cfg.NoNSFW && nsfw {
		s.MessageReactionAdd(m.ChannelID, m.ID, "🔞")
		return
	}

	for _, msg := range messages {
		_, err = s.ChannelMessageSendComplex(m.ChannelID, &msg)
		if err != nil {
			if strings.HasPrefix(err.Error(), "HTTP 413") {
				s.MessageReactionAdd(m.ChannelID, m.ID, "🥵")
				return
			}

			goto ERROR
		}
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
