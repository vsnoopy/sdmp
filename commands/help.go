package commands

import (
	"github.com/bwmarrin/discordgo"
)

func HandleHelp(s *discordgo.Session, i *discordgo.InteractionCreate) {
	commands := map[string]string{
		"/play [search query]": "Play a song from YouTube.",
		"/pause":               "Pause the current song.",
		"/resume":              "Resume the current song.",
		"/skip":                "Skip the current song.",
		"/stop":                "Stop the current song.",
		"/queue":               "Display the current queue.",
		"/ping":                "Check if the bot is online and response time.",
		"/invite":              "Get the invite link for the bot.",
		"/help":                "Display the help message.",
	}

	fields := make([]*discordgo.MessageEmbedField, 0, len(commands))
	for cmd, desc := range commands {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   cmd,
			Value:  desc,
			Inline: false,
		})
	}

	embed := &discordgo.MessageEmbed{
		Title:       "Help - List of Commands",
		Description: "Here are the available commands:",
		Color:       0x00FF00, // Green color
		Fields:      fields,
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})
}
