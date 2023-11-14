variable "client_id" {
  desciption = "Client ID generated when creating an API project at Spotify developer dashboard"
  type       = string
}

variable "client_secret" {
  desciption = "Client Secret generated when creating an API project at Spotify developer dashboard"
  type       = string
}

variable "spotify_username" {
  desciption = "Username of a Spotify account to request access tokens for the Spotify API"
  type       = string
}

variable "spotify_password" {
  desciption = "Password of Spotify account with previous username"
  type       = string
}
