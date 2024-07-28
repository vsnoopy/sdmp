package commands

import (
	"github.com/bwmarrin/discordgo"
)

// HandleInvite is a command that sends the bot invite link
func HandleInvite(s *discordgo.Session, i *discordgo.InteractionCreate) {
	inviteLink := "https://discord.com/oauth2/authorize?client_id=1265334719025774734&scope=bot&permissions=8"
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title: "Invite me to your server",
					Color: 0x0000ff, // Blue color
				},
			},
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label: "Invite",
							Style: discordgo.LinkButton,
							URL:   inviteLink,
						},
					},
				},
			},
		},
	})
}
