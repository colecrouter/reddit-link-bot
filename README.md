# reddit-link-bot

A Discord bot that automatically replies to Reddit links with inline media (videos and images).

## Usage

```bash
go run ./cmd/main.go
```

\- or (recommended) -

```bash
go run -tags=goexperiment.arenas ./cmd/main.go
```

This will generate a `config.yml` file in the current directory. Insert your Discord bot token and restart the bot.

## Requirements

- Must have FFMPEG installed and exported/added to PATH