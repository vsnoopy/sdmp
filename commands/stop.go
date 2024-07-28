package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"sdmp/storage"
)

func HandleStop(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})

	go func() {
		StopAudio()
		storage.GetSongQueue().Clear()
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Embeds: &[]*discordgo.MessageEmbed{
				{
					Title:       "Music Player",
					Description: fmt.Sprintf("Stopped music and cleared the queue"),
					Color:       0xff0000,
				},
			},
		})
	}()
}
