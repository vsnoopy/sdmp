package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"os"

	"sdmp/commands"
)

// InitBot initializes the bot and registers the slash commands
func InitBot() (*discordgo.Session, error) {
	Token := os.Getenv("DISCORD_TOKEN")
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		return nil, fmt.Errorf("error creating Discord session: %w", err)
	}

	// Register the slash command handler
	dg.AddHandler(commands.SlashCommandHandler)

	// Open a websocket connection to Discord and begin listening
	err = dg.Open()
	if err != nil {
		return nil, fmt.Errorf("error opening connection: %w", err)
	}

	for _, cmd := range commands.GetCommands() {
		_, err = dg.ApplicationCommandCreate(dg.State.User.ID, "", cmd)
		if err != nil {
			return nil, fmt.Errorf("error creating command: %w", err)
		}
	}

	return dg, nil
}
