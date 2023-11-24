variable "project_tag" {
  description = "Tag of project used as identifier for resources"
  type        = string
}

variable "project" {
  description = "Poject name"
  type        = string
}

variable "env" {
  description = "Environment"
  type        = string
}

variable "static_secrets_arn" {
  description = "ARN of the secret with the static secrets"
  type        = string
}

variable "token_secret_arn" {
  description = "ARN of the secret with the access token"
  type        = string
}
