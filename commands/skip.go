package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

func HandleSkip(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})

	go func() {
		StopAudio()
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Embeds: &[]*discordgo.MessageEmbed{
				{
					Title:       "Music Player",
					Description: fmt.Sprintf("Skipped the current song"),
					Color:       0xffff00,
				},
			},
		})
	}()
}
