package discord

import (
	"fmt"
	"strings"

	"github.com/Mexican-Man/reddit-bot/pkg/config"
	"github.com/Mexican-Man/reddit-bot/pkg/reddit"
	"github.com/bwmarrin/discordgo"
)

type DiscordBot struct {
	config  config.Config
	discord *discordgo.Session
	reddit  *reddit.RedditBot
}

func NewDiscordBot(cfg config.Config) *DiscordBot {
	redditBot := reddit.NewRedditBot(cfg)

	// Start bot
	discord, err := discordgo.New("Bot " + cfg.Discord.Token)
	if err != nil {
		panic(err)
	}

	b := &DiscordBot{
		config:  cfg,
		discord: discord,
		reddit:  redditBot,
	}

	// Add handlers
	discord.AddHandler(b.messageCreate)
	discord.AddHandler(b.ready)

	// Set intents
	discord.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages)

	return &DiscordBot{
		config:  cfg,
		discord: discord,
		reddit:  redditBot,
	}
}

func (d *DiscordBot) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
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
	if len(d.config.Discord.Roles) > 0 {
		member, err := s.GuildMember(m.GuildID, m.Author.ID)
		if err != nil {
			goto ERROR
		}

		hasRole := false
		for _, role := range d.config.Discord.Roles {
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
	if len(d.config.Discord.Channels) > 0 {
		hasChannel := false
		for _, channel := range d.config.Discord.Channels {
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

	messages, nsfw, err = d.toDiscordMessages(m.Content)
	if err != nil {
		goto ERROR
	}

	// Check if the message is NSFW
	if d.config.Discord.NoNSFW && nsfw {
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

func (d *DiscordBot) ready(s *discordgo.Session, e *discordgo.Ready) {
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

func (d *DiscordBot) Close() {
	d.discord.Close()
}

func (d *DiscordBot) Start() error {
	return d.discord.Open()
}
