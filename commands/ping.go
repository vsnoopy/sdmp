package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"time"
)

// HandlePing is a command that replies with pong and the response time in ms
func HandlePing(s *discordgo.Session, i *discordgo.InteractionCreate) {
	start := time.Now()
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "Pong!",
					Description: "Calculating response time...",
					Color:       0x00ff00, // Green color
				},
			},
		},
	})
	elapsed := time.Since(start).Milliseconds()
	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{
			{
				Title:       "Pong!",
				Description: fmt.Sprintf("Response time: %d ms", elapsed),
				Color:       0x00ff00, // Green color
			},
		},
	})
}
