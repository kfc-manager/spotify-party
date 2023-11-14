resource "aws_secretsmanager_secret" "static" {
  name                    = "${var.project_tag}-static-secrets-v2"
  recovery_window_in_days = 0

  tags = {
    Project     = var.project
    Environment = var.env
    Description = "Static Secrets for requesting an access token from the Spotify API"
  }
}

resource "aws_secretsmanager_secret_version" "static" {
  secret_id     = aws_secretsmanager_secret.static.id
  secret_string = jsonencode(var.static_secrets)
}

resource "aws_secretsmanager_secret" "access_token" {
  name                    = "${var.project_tag}-access-token-v2"
  recovery_window_in_days = 0

  tags = {
    Project     = var.project
    Environment = var.env
    Description = "Secret of stored token to access the Spotify API"
  }
}

resource "aws_secretsmanager_secret_version" "access_token" {
  secret_id     = aws_secretsmanager_secret.access_token.id
  secret_string = "placeholder"
}

