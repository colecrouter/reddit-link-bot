package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Mexican-Man/reddit-bot/pkg/config"
	"github.com/Mexican-Man/reddit-bot/pkg/discord"
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

	// Start custom bot
	discord := discord.NewDiscordBot(cfg)
	err = discord.Start()
	if err != nil {
		fmt.Println(err)
		return
	}

	defer discord.Close()

	// Stop from closing
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
