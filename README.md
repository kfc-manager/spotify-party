# Spotify Party :musical_note:

Spotify Party is a web application to control the queue of my personal Spotify account. Initial motivation behind the project was to be able to share a queue of someone, so people at a party have the possibility to add a song to the queue on their device and not have to pass around the smartphone of the account owner. The application is fully deployed on AWS. The backend consists of an API Gateway and Lambda functions, which makes it fully serverless.

## How it works :thinking:

![alt text](https://github.com/kfc-manager/spotify-party/blob/main/assets/architecture_transparent.png?raw=true)

### &rarr; Frontend :desktop_computer:

- Domain is registered with `Route53` and points to `CloudFront`
- `CloudFront` fetches and caches the static files in order for the `React App` to work
- `React App` makes requests to the `API Gateway` to view and make changes to the Spotify queue

### &rarr; Backend :gear:

- `API Gateway` suppplies endpoints to invoke the `Lambda functions`
- `get-queue`, `search-tracks` and `update-queue` use the access token stored in the `SecretsManager` to access the Spotify API
- `api-login` and `api-callback` use the `client_id` and `client_secret` to fetch a token from the Spotify API

![alt text](https://github.com/kfc-manager/spotify-party/blob/main/assets/spotify_auth_code_flow.png?raw=true)

### &rarr; `api-login`

- When request of the `React App` returns an `UNAUTHORIZED` error user gets redirected to this endpoint
- Builds request with URL encoded parameters and redirects to the Spotify login

### &rarr; `api-callback`

- After successful Spotify login the user gets redirected to this endpoint
- Forms request with `code` provided from Spotify for an access token
- Saves fresh token from response in the `SecreatsManager` and redirects user back to the base domain

## Deployment :mechanic:

### &rarr; Lambdas :purple_circle:

- Create ECR repositories for the images of the Lambda functions, with name `spotify-party-{name of the Lambda}`
- Build the image of each Lambda except `token-caller`
- Push the images to it's corresponding repository in the registry

### &rarr; React App :large_blue_circle:

- Import required dependencies with `npm install`
- Build the static `js`, `css` and `html` files with `npm run build`

### &rarr; AWS infrastructure :orange_circle:

- Run `terraform init` to initialize required modules
- Then run `terraform apply` and supply the `client_id` and `client_secret` given by Spotify for the API registration

The Lambda functions will then automatically draw their images from the created registries and the static files for the React App will be uploaded as `S3 objects`.
