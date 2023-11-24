variable "project_tag" {
  description = "Tag of project used as identifier for resources"
  type        = string
}

variable "static_secrets" {
  description = "Static secrets for Spofity API login"
  type        = map(string)
}

variable "project" {
  description = "Poject name"
  type        = string
}

variable "env" {
  description = "Environment"
  type        = string
}

