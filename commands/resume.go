package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

func HandleResume(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})

	go func() {
		ResumeAudio()
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Embeds: &[]*discordgo.MessageEmbed{
				{
					Title:       "Music Player",
					Description: fmt.Sprintf("Resumed the music"),
					Color:       0x00ff00,
				},
			},
		})
	}()
}

func ResumeAudio() {
	pauseSignal <- false
}
