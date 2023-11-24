package main

import (
	"encoding/base64"
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

type Secrets struct {
	ClientID      string `json:"client_id"`
	ClientSecrets string `json:"client_secret"`
}

type SpotifyAPIResponse struct {
	Token            string `json:"access_token"`
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

func getStaticSecrets(serviceClient *secretsmanager.SecretsManager) (*Secrets, error) {

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

	return secrets, nil
}

func setTokenSecret(serviceClient *secretsmanager.SecretsManager, token string) error {

	input := &secretsmanager.UpdateSecretInput{
		SecretId:     aws.String(os.Getenv("TOKEN_SECRET_ID")),
		SecretString: aws.String(token),
	}
	_, err := serviceClient.UpdateSecret(input)
	if err != nil {
		return err
	}

	return nil
}

func buildRequest(code string, secrets *Secrets) (*http.Request, error) {

	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", nil)
	if err != nil {
		return nil, err
	}
	params := req.URL.Query()
	params.Add("code", code)
	params.Add("redirect_uri", os.Getenv("REDIRECT_URI"))
	params.Add("grant_type", "authorization_code")
	req.URL.RawQuery = params.Encode()
	req.Header.Add(
		"Authorization",
		"Basic "+base64.StdEncoding.EncodeToString(
			[]byte(secrets.ClientID+":"+secrets.ClientSecrets),
		),
	)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	return req, nil
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
	secrets, err := getStaticSecrets(serviceClient)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Internal Server Error",
		}, err
	}

	// build request for the Spotify API
	req, err := buildRequest(event.QueryStringParameters["code"], secrets)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Internal Server Error",
		}, nil
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

	// error handle Spotify API response with bad status code
	if res.StatusCode != 200 {
		return events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       "Internal Server Error",
			}, errors.New(
				"Spotify API responded with status code: " + fmt.Sprint(
					res.StatusCode,
				) + ", error: " + body.Error + ", and message: " + body.ErrorDescription,
			)
	}

	// put token received from the Spotify API into SecretsManager token secret
	if err = setTokenSecret(serviceClient, body.Token); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Internal Server Error",
		}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 301,
		Body:       "Redirecting...",
		Headers: map[string]string{
			"Location": os.Getenv("BASE_URI"),
		},
		IsBase64Encoded: false,
	}, nil
}

func main() {
	lambda.Start(handler)
}
