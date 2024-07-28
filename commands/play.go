package commands

import (
	"encoding/binary"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io"
	"layeh.com/gopus"
	"os"
	"os/exec"
	"sdmp/api"
	"sdmp/storage"
)

const (
	channels  int = 2                   // 1 for mono, 2 for stereo
	frameRate int = 48000               // audio sampling rate
	frameSize int = 960                 // uint16 size of each audio frame
	maxBytes  int = (frameSize * 2) * 2 // max size of opus data
)

var stopSignal = make(chan bool)
var pauseSignal = make(chan bool)

// HandlePlay handles the play command
func HandlePlay(s *discordgo.Session, i *discordgo.InteractionCreate) {
	query := i.ApplicationCommandData().Options[0].StringValue()

	// Acknowledge the interaction immediately
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})

	go func() {
		// Get the voice state of the user who invoked the command
		voiceState, err := GetVoiceState(s, i.GuildID, i.Member.User.ID)
		if err != nil {
			s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Embeds: &[]*discordgo.MessageEmbed{
					{
						Title:       "Error",
						Description: fmt.Sprintf("Error getting voice state: %v", err),
						Color:       0xff0000,
					},
				},
			})
			return
		}

		// Check if the user is in a voice channel
		if voiceState == nil || voiceState.ChannelID == "" {
			s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Embeds: &[]*discordgo.MessageEmbed{
					{
						Title:       "Error",
						Description: fmt.Sprintf("You need to join a voice channel first!"),
						Color:       0xff0000,
					},
				},
			})
			return
		}

		// Search YouTube for the query
		song, err := api.SearchYouTube(query)
		if err != nil {
			s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Embeds: &[]*discordgo.MessageEmbed{
					{
						Title:       "Error",
						Description: fmt.Sprintf("Error searching YouTube: %v", err),
						Color:       0xff0000,
					},
				},
			})
			return
		}

		// Download the audio
		err = api.DownloadAudio(song)
		if err != nil {
			s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Embeds: &[]*discordgo.MessageEmbed{
					{
						Title:       "Error",
						Description: fmt.Sprintf("Error downloading audio: %v", err),
						Color:       0xff0000,
					},
				},
			})
			return
		}

		// Add the song to the queue
		songQueue := storage.GetSongQueue()
		songQueue.Add(*song)
		if len(songQueue.Songs) == 1 {
			PlayNextSong(s, i, voiceState.ChannelID)
		} else {
			s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Embeds: &[]*discordgo.MessageEmbed{
					{
						Title: fmt.Sprintf("Added to queue: %s", song.Title),
						Color: 0x0000ff,
					},
				},
			})
		}
	}()
}

// PlayAudio plays the audio of a song
func PlayAudio(s *discordgo.Session, guildID, channelID string, song *storage.Song) error {
	// Join vc
	vc, err := s.ChannelVoiceJoin(guildID, channelID, false, true)
	if err != nil {
		return fmt.Errorf("error joining voice channel: %w", err)
	}
	defer vc.Disconnect()

	// Send speaking packet
	vc.Speaking(true)
	defer vc.Speaking(false)

	// Start ffmpeg
	ffmpeg := exec.Command("ffmpeg", "-i", song.Id+".mp3", "-f", "s16le", "-ar", "48000", "-ac", "2", "pipe:1")
	ffmpegbuf, err := ffmpeg.StdoutPipe()
	if err != nil {
		return fmt.Errorf("error creating ffmpeg stdout pipe: %w", err)
	}

	if err := ffmpeg.Start(); err != nil {
		return fmt.Errorf("error starting ffmpeg: %w", err)
	}

	defer ffmpeg.Process.Kill()

	err = vc.Speaking(true)
	if err != nil {
		return fmt.Errorf("error sending speaking packet: %w", err)
	}

	defer func() {
		err := vc.Speaking(false)
		if err != nil {
			fmt.Println("Error stopping speaking:", err)
		}
	}()
	send := make(chan []int16, 2)
	defer close(send)

	close := make(chan bool)
	go func() {
		SendPCM(vc, send)
		close <- true
	}()

	// Send audio
	for {
		select {
		case <-stopSignal:
			return nil
		case <-pauseSignal:
			<-pauseSignal // Wait for resume signal
		default:
			audiobuf := make([]int16, frameSize*channels)
			err = binary.Read(ffmpegbuf, binary.LittleEndian, &audiobuf)
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				return nil
			}
			if err != nil {
				return fmt.Errorf("error reading from ffmpeg stdout: %s", err)
			}

			select {
			case send <- audiobuf:
			case <-close:
				return nil
			}
		}
	}
}

// StopAudio stops the audio
func StopAudio() {
	stopSignal <- true
}

// PauseAudio pauses the audio
func PauseAudio() {
	pauseSignal <- true
}

// PlayNextSong plays the next song in the queue
func PlayNextSong(s *discordgo.Session, i *discordgo.InteractionCreate, channelID string) {
	songQueue := storage.GetSongQueue()
	song := songQueue.Peek()
	if song == nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "Music Player",
						Description: fmt.Sprintf("Queue is empty"),
						Color:       0x0000ff,
					},
				},
			},
		})
		return
	}

	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{
			{
				Title:       "Music Player",
				Description: fmt.Sprintf("Now Playing: %s \nSource: %s", song.Title, song.URL),
				Color:       0x00ff00,
			},
		},
	})

	err := PlayAudio(s, i.GuildID, channelID, song)
	if err != nil {
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Embeds: &[]*discordgo.MessageEmbed{
				{
					Title:       "Error",
					Description: fmt.Sprintf("Error playing audio: %v", err),
					Color:       0xff0000,
				},
			},
		})
		return
	}

	// Delete the mp3 file after playing
	err = os.Remove(song.Id + ".mp3")
	if err != nil {
		fmt.Println("Error deleting file:", err)
	}

	select {
	case <-stopSignal:
		songQueue.Remove()
		return
	default:
		songQueue.Remove()
		PlayNextSong(s, i, channelID)
	}
}

// GetVoiceState returns the voice state of a user in a guild
func GetVoiceState(s *discordgo.Session, guildID, userID string) (*discordgo.VoiceState, error) {
	guild, err := s.State.Guild(guildID)
	if err != nil {
		return nil, err
	}

	for _, vs := range guild.VoiceStates {
		if vs.UserID == userID {
			return vs, nil
		}
	}

	return nil, nil
}

// SendPCM sends PCM data to the provided voice connection
func SendPCM(v *discordgo.VoiceConnection, pcm <-chan []int16) {
	if pcm == nil {
		return
	}

	var err error

	opusEncoder, err := gopus.NewEncoder(frameRate, channels, gopus.Audio)
	if err != nil {
		fmt.Println("NewEncoder Error:", err)
		return
	}

	for {
		select {
		case <-stopSignal:
			return
		default:
			recv, ok := <-pcm
			if !ok {
				fmt.Println("PCM Channel closed")
				return
			}

			opus, err := opusEncoder.Encode(recv, frameSize, maxBytes)
			if err != nil {
				fmt.Println("Encoding Error:", err)
				return
			}

			if v.Ready == false || v.OpusSend == nil {
				return
			}

			v.OpusSend <- opus
		}
	}
}
