package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

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

func buildRequest(token string, id string) (*http.Request, error) {

	req, err := http.NewRequest("POST", "https://api.spotify.com/v1/me/player/queue", nil)
	if err != nil {
		return nil, err
	}

	params := req.URL.Query()
	params.Add("uri", "spotify:track:"+id)
	req.URL.RawQuery = params.Encode()

	req.Header.Add("Authorization", "Bearer "+token)

	return req, nil
}

func handler(
	event events.APIGatewayProxyRequest,
) (events.APIGatewayProxyResponse, error) {

	// check if id is provided
	if event.QueryStringParameters["id"] == "" {
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
	req, err := buildRequest(*token, event.QueryStringParameters["id"])
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

	// handle expired access token to Spofity API
	if res.StatusCode == 401 {
		return events.APIGatewayProxyResponse{
			StatusCode: 401,
			Body:       "Invalid Access Token",
		}, nil
	}

	// error handle Spotify API response with other bad status codes
	if res.StatusCode != 204 {
		return events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       "Internal Server Error",
			}, errors.New(
				"Spotify API responded with status code: " + fmt.Sprint(
					res.StatusCode,
				),
			)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "Success",
	}, nil
}

func main() {
	lambda.Start(handler)
}
