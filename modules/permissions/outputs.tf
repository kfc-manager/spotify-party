output "read_static_secrets_role_arn" {
  description = "ARN of the IAM Role to read static secrets"
  value       = aws_iam_role.read_static_secrets.arn
}

output "read_token_secret_role_arn" {
  description = "ARN of the IAM Role to read token secret"
  value       = aws_iam_role.read_token_secret.arn
}

output "write_token_secret_role_arn" {
  description = "ARN of the IAM Role to read static secrets and write token secret"
  value       = aws_iam_role.write_token_secret.arn
}
