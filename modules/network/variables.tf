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

variable "availability_zone" {
  description = "Availability zone of the VPC"
  type        = string
}
