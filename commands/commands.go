package commands

import "github.com/bwmarrin/discordgo"

// GetCommands returns a list of all application commands
func GetCommands() []*discordgo.ApplicationCommand {
	return []*discordgo.ApplicationCommand{
		{
			Name:        "ping",
			Description: "Replies with pong and the response time in ms",
		},
		{
			Name:        "invite",
			Description: "Sends the bot invite link",
		},
		{
			Name:        "play",
			Description: "Plays a song",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "query",
					Description: "The song to play",
					Type:        discordgo.ApplicationCommandOptionString,
					Required:    true,
				},
			},
		},
		{
			Name:        "skip",
			Description: "Skips the current song",
		},
		{
			Name:        "stop",
			Description: "Stops the current song",
		},
		{
			Name:        "pause",
			Description: "Pauses the current song",
		},
		{
			Name:        "resume",
			Description: "Resumes the current song",
		},
		{
			Name:        "help",
			Description: "Displays the help message",
		},
		{
			Name:        "queue",
			Description: "Displays the current queue",
		},
	}
}

func SlashCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.ApplicationCommandData().Name {
	case "ping":
		HandlePing(s, i)
	case "invite":
		HandleInvite(s, i)
	case "play":
		HandlePlay(s, i)
	case "skip":
		HandleSkip(s, i)
	case "stop":
		HandleStop(s, i)
	case "pause":
		HandlePause(s, i)
	case "resume":
		HandleResume(s, i)
	case "help":
		HandleHelp(s, i)
	case "queue":
		HandleQueue(s, i)
	}
}
