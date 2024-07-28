# sdmp - Discord Music Bot

A simple discord music bot that can play music from youtube.

## Installation

### Docker
The easiest way to run the bot is to use docker. First, you need to have docker installed on your system. You can install docker by following the instructions [here](https://docs.docker.com/get-docker/).

Next, make sure you have a `.env` file in the root directory of the repo with your tokens as described in the [Usage](#usage) section.

To build and run the docker image, run the following command:
```bash
make docker
```


### Local Installation
First, you need to have the following dependencies installed:
- [ffmpeg](https://ffmpeg.org/)
- [yt-dlp](https://github.com/yt-dlp/yt-dlp)

These can be installed easily on most linux distributions using the package manager. For example, on Fedora:
```bash
sudo dnf install ffmpeg yt-dlp
```

To install & compile, clone the repository and run the following command:
```bash
make build
```

## Usage
To run the bot, you need to have a discord bot token and a YouTube API token. Once you have these, create a `.env` file in the root directory of the repository and add the following lines:
```bash
DISCORD_TOKEN=[YOUR_DISCORD_BOT_TOKEN]
YOUTUBE_API_TOKEN=[YOUR_YOUTUBE_API_TOKEN]
```

To run the bot, use the following command:
```bash
make run
```
The bot runs on port 8080 by default.

## Commands
- /play [search query]: Play a song from YouTube.
- /pause: Pause the current song.
- /resume: Resume the current song.
- /skip: Skip the current song.
- /stop: Stop the current song.
- /queue: Display the current queue.
- /ping: Check if the bot is online and response time.
- /invite: Get the invite link for the bot.
- /help: Display the help message.
