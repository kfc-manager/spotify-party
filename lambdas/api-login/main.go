package main

import (
	"encoding/json"
	"net/url"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

type Secrets struct {
	ClientID      string `json:"client_id"`
	ClientSecrets string `json:"client_secret"`
}

func getClientID(serviceClient *secretsmanager.SecretsManager) (*string, error) {

	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(os.Getenv("STATIC_SECRETS_ID")),
	}
	result, err := serviceClient.GetSecretValue(input)
	if err != nil {
		return nil, err
	}

	secrets := &Secrets{}
	if err := json.Unmarshal([]byte(*result.SecretString), secrets); err != nil {
		return nil, err
	}

	return &secrets.ClientID, nil
}

func buildAPIRedirect(clientID string) (*string, error) {

	params := url.Values{}
	params.Add("response_type", "code")
	params.Add("client_id", clientID)
	params.Add("redirect_uri", os.Getenv("REDIRECT_URI"))
	params.Add(
		"scope",
		"user-read-currently-playing user-read-playback-state user-modify-playback-state",
	)

	redirect := "https://accounts.spotify.com/authorize?" + params.Encode()

	return &redirect, nil
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

	// retrieve static secrets from the SecretsManager
	clientID, err := getClientID(serviceClient)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Internal Server Error",
		}, err
	}

	// build redirect for the Spotify API login
	redirectLogin, err := buildAPIRedirect(*clientID)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Internal Server Error",
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 301,
		Body:       "Redirecting...",
		Headers: map[string]string{
			"Location": *redirectLogin,
		},
		IsBase64Encoded: false,
	}, nil
}

func main() {
	lambda.Start(handler)
}
