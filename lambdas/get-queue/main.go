package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

type SpotifyAPIResponse struct {
	CurrentlyPlaying *SpotifyAPISong  `json:"currently_playing"`
	Queue            []SpotifyAPISong `json:"queue"`
}

type SpotifyAPISong struct {
	ID       string             `json:"id"`
	Name     string             `json:"name"`
	Artists  []SpotifyAPIArtist `json:"artists"`
	Album    SpotifyAPIAlbum    `json:"album"`
	Duration int                `json:"duration_ms"`
}

type SpotifyAPIArtist struct {
	Name string `json:"name"`
}

type SpotifyAPIAlbum struct {
	Images []SpotifyAPIImages `json:"images"`
}

type SpotifyAPIImages struct {
	URL string `json:"url"`
}

type Response struct {
	Queue []Song `json:"queue"`
}

type Song struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Image    string   `json:"image_url"`
	Artists  []string `json:"artists"`
	Duration int      `json:"duration_ms"`
}

func getTokenSecret(serviceClient *secretsmanager.SecretsManager) (*string, error) {

	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(os.Getenv("TOKEN_SECRET_ID")),
	}
	result, err := serviceClient.GetSecretValue(input)
	if err != nil {
		return nil, err
	}
	token := result.SecretString

	return token, nil
}

func buildRequest(token string) (*http.Request, error) {

	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/me/player/queue", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+token)

	return req, nil
}

func transformQueue(spotifyRes SpotifyAPIResponse) []Song {

	queue := []Song{}
	if spotifyRes.CurrentlyPlaying == nil {
		return queue
	}
	queue = append(queue, *transformSong(*spotifyRes.CurrentlyPlaying))
	for _, elem := range spotifyRes.Queue {
		queue = append(queue, *transformSong(elem))
	}

	return queue
}

func transformSong(spotifySong SpotifyAPISong) *Song {

	song := &Song{}
	song.Name = spotifySong.Name
	song.ID = spotifySong.ID
	if len(spotifySong.Album.Images) > 0 {
		song.Image = spotifySong.Album.Images[0].URL
	}
	song.Duration = spotifySong.Duration
	for _, elem := range spotifySong.Artists {
		song.Artists = append(song.Artists, elem.Name)
	}

	return song
}

func handler(
	event events.APIGatewayProxyRequest,
) (events.APIGatewayProxyResponse, error) {

	// create AWS SecretsManager session
	awsSession, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("REGION")),
	})
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Internal Server Error",
		}, err
	}
	serviceClient := secretsmanager.New(awsSession)

	// retrieve token secret from the SecretsManager
	token, err := getTokenSecret(serviceClient)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Internal Server Error",
		}, err
	}

	// build request for the Spotify API
	req, err := buildRequest(*token)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Internal Server Error",
		}, err
	}

	// send request to the Spotify API
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Internal Server Error",
		}, err
	}
	defer res.Body.Close()

	// read response from the Spotify API
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Internal Server Error",
		}, err
	}
	body := &SpotifyAPIResponse{}
	if err = json.Unmarshal(bodyBytes, body); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Internal Server Error",
		}, errors.New(err.Error() + "; response body: " + string(bodyBytes))
	}

	// handle expired access token to Spofity API
	if res.StatusCode == 401 {
		return events.APIGatewayProxyResponse{
			StatusCode: 401,
			Body:       "Invalid Access Token",
		}, nil
	}

	// error handle Spotify API response with other bad status codes
	if res.StatusCode != 200 {
		return events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       "Internal Server Error",
			}, errors.New(
				"Spotify API responded with status code: " + fmt.Sprint(
					res.StatusCode,
				),
			)
	}

	// build response
	queueRes := &Response{Queue: transformQueue(*body)}
	resBytes, err := json.Marshal(queueRes)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Internal Server Error",
		}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(resBytes),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}

func main() {
	lambda.Start(handler)
}
