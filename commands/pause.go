package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

func HandlePause(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})

	go func() {
		PauseAudio()
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Embeds: &[]*discordgo.MessageEmbed{
				{
					Title:       "Music Player",
					Description: fmt.Sprintf("Paused the music"),
					Color:       0xffff00,
				},
			},
		})
	}()
}
