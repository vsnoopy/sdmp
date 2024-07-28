package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"sdmp/storage"
)

func HandleQueue(s *discordgo.Session, i *discordgo.InteractionCreate) {
	songQueue := storage.GetSongQueue()
	if len(songQueue.Songs) == 0 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "The queue is currently empty.",
			},
		})
		return
	}

	fields := make([]*discordgo.MessageEmbedField, 0, len(songQueue.Songs))
	for idx, song := range songQueue.Songs {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   fmt.Sprintf("%d. %s", idx+1, song.Title),
			Value:  song.URL,
			Inline: false,
		})
	}

	embed := &discordgo.MessageEmbed{
		Title:       "Current Song Queue",
		Description: "Here are the songs currently in the queue:",
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
