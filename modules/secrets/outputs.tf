output "static_arn" {
  description = "ARN of the secret with the static secrets"
  value       = aws_secretsmanager_secret.static.arn
}

output "access_token_arn" {
  description = "ARN of the secret with the access token"
  value       = aws_secretsmanager_secret.access_token.arn
}
