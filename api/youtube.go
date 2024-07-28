package api

import (
	"bytes"
	"context"
	"fmt"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	"os"
	"os/exec"
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

// DownloadAudio downloads the audio of a song using yt-dlp
func DownloadAudio(song *storage.Song) error {
	cmd := exec.Command("yt-dlp", "-f", "bestaudio", "-x", "--audio-format", "mp3", "-o", song.Id, song.URL)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error downloading audio: %w, %s", err, stderr.String())
	}
	return nil
}
