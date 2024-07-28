package api

import (
	"context"
	"fmt"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	"os"
	"os/exec"
	"path/filepath"
	"sdmp/storage"
)

// SearchYouTube searches YouTube for a video matching the query
func SearchYouTube(query string) (*storage.Song, error) {
	apiKey := os.Getenv("YOUTUBE_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("YOUTUBE_API_KEY environment variable not set")
	}

	ctx := context.Background()
	service, err := youtube.NewService(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("error creating YouTube service: %w", err)
	}

	call := service.Search.List([]string{"id", "snippet"}).Q(query).MaxResults(1).Type("video")
	response, err := call.Do()
	if err != nil {
		return nil, fmt.Errorf("error making YouTube search API call: %w", err)
	}

	if len(response.Items) == 0 {
		return nil, fmt.Errorf("no results found for query: %s", query)
	}

	videoID := response.Items[0].Id.VideoId
	videoURL := fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoID)
	title := response.Items[0].Snippet.Title

	song := &storage.Song{
		URL:   videoURL,
		Id:    videoID,
		Title: title,
	}

	return song, nil
}

// DownloadAudio downloads the audio of a song
func DownloadAudio(song *storage.Song) error {
	audioDir := "/app/audio-files"
	audioPath := filepath.Join(audioDir, song.Id+".mp3")

	// Create the audio directory if it doesn't exist
	if err := os.MkdirAll(audioDir, os.ModePerm); err != nil {
		return fmt.Errorf("error creating audio directory: %w", err)
	}

	// Use yt-dlp to download the audio
	cmd := exec.Command("yt-dlp", "-x", "--audio-format", "mp3", "-o", audioPath, song.URL)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error downloading audio: %w", err)
	}

	return nil
}
