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

variable "domain_name" {
  description = "Name of base domain"
  type        = string
}

variable "public_host_zone_id" {
  description = "ID of public DNS host zone"
  type        = string
}

variable "acm_certificate_arn" {
  description = "ARN of ACM certificate for base domain"
  type        = string
}

variable "account_arn" {
  description = "ARN of account which deploys this infrastructure"
  type        = string
}
