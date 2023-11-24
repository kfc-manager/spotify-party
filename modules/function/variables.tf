variable "project" {
  description = "Poject name"
  type        = string
}

variable "env" {
  description = "Environment"
  type        = string
}

variable "region" {
  description = "Region of the ECR repository which holds the image for this Lambda Function"
  type        = string
}

variable "account_id" {
  description = "ID of the account that deploys this structure"
  type        = string
}

variable "name" {
  description = "Unique name of the Lambda Function"
  type        = string
}

variable "iam_role_arn" {
  description = "ARN of IAM Role for the Lambda Function to assume"
  type        = string
}

variable "ecr_image_name" {
  description = "Name of the image in the ECR repository for the Lambda Function to execute"
  type        = string
}

variable "description" {
  description = "Description of the purpose of this Lambda Function"
  type        = string
}

variable "env_variables" {
  description = "Environment variables for in the Lambda Function"
  type        = map(string)
}

variable "memory_size" {
  description = "Size of allocated memory to Lambda Function"
  type        = number
  default     = 128
}

variable "timeout" {
  description = "Time after which the Lambda Function should be stopped"
  type        = number
  default     = 3
}

variable "vpc_config" {
  description = "Parameter to add the Lambda Function to a VPC"
  type        = map(list(string))
  default     = {}
}
