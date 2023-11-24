terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "5.21.0"
    }
  }

  backend "s3" {
    bucket = "spotify-party-terraform-state"
    key    = "state"
    region = "eu-central-1"
  }
}
