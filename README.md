# reddit-link-bot

A Discord bot that automatically replies to Reddit links with inline media (videos and images).

## Purpose

A while back, Reddit redid their platform's OpenGraph link previews/embeds. Previously, you would get single image previews (static images and GIFs, not videos) when sharing links to posts. Now, Reddit now only serves generated image with watermarks (een for GIFs).

This bot solves this issue by responding to messages with the embeded media from Reddit, so you don't have to visit Reddit to vew it/them.

## Features

- Static images & GIFs
- Videos w/ or w/o sound
- Galleries (posts with multiple images)
- Marks spoiler or NSFW content as a spoiler

## Usage

```bash
go run ./cmd/main.go
```

This will generate a `config.yml` file in the current directory. Insert your Discord bot token and restart the bot.

## Requirements

- Must have FFMPEG installed and exported/added to PATH. Does not require FFMPEG build with HTTPS support.
