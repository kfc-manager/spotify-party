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

variable "region" {
  description = "Region of the API Gateway"
  type        = string
}

variable "account_id" {
  description = "ID of the account that deploys this structure"
  type        = string
}

variable "base_uri" {
  description = "Base URI or domain of the application for allowing CORS"
  type        = string
}

variable "api_domain_name_id" {
  description = "ID of the by domain created API domain name"
  type        = string
}

variable "callback_lambda_invoke_arn" {
  description = "Invoke ARN of the Lambda Function for the Spotify API callback"
  type        = string
}

variable "callback_lambda_arn" {
  description = "ARN of the Lambda Function for the Spotify API callback"
  type        = string
}

variable "lambda_routes" {
  description = "Information about the Lambda Function that needs to be mapped to a route of the API Gateway"
  type        = list(map(string))
}
