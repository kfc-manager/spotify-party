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
	SearchResult SpotifyAPISearchResult `json:"tracks"`
}

type SpotifyAPISearchResult struct {
	Items []SpotifyAPISong `json:"items"`
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
	Tracks []Song `json:"tracks"`
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

func buildRequest(token string, query string) (*http.Request, error) {

	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/search", nil)
	if err != nil {
		return nil, err
	}

	params := req.URL.Query()
	params.Add("q", query)
	params.Add("type", "track")
	req.URL.RawQuery = params.Encode()

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)

	return req, nil
}

func transformSearchResult(spotifyRes SpotifyAPIResponse) []Song {

	result := []Song{}
	for _, elem := range spotifyRes.SearchResult.Items {
		result = append(result, *transformSong(elem))
	}

	return result
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

	// check if query is provided
	if event.QueryStringParameters["query"] == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Missing Parameter",
		}, nil
	}

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
	req, err := buildRequest(*token, event.QueryStringParameters["query"])
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
	tracksRes := &Response{Tracks: transformSearchResult(*body)}
	resBytes, err := json.Marshal(tracksRes)
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
