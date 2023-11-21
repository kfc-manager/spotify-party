data "aws_caller_identity" "current" {}
data "aws_region" "current" {}

locals {
  project        = "Spotify Party"
  project_tag    = "spotify-party"
  env            = "Production"
  domain_name    = "kiliansqueue.com"
  region         = data.aws_region.current.name
  account_id     = data.aws_caller_identity.current.account_id
  callback_route = "/callback"
  base_uri       = "https://${local.domain_name}"
  redirect_uri   = "${local.base_uri}${local.callback_route}"
}

module "domain" {
  source = "./modules/domain"

  domain_name = local.domain_name
}

module "secrets" {
  source = "./modules/secrets"

  project     = local.project
  env         = local.env
  project_tag = local.project_tag
  static_secrets = {
    client_id     = var.client_id
    client_secret = var.client_secret
  }
}

module "permissions" {
  source = "./modules/permissions"

  project            = local.project
  env                = local.env
  project_tag        = local.project_tag
  static_secrets_arn = module.secrets.static_arn
  token_secret_arn   = module.secrets.access_token_arn
}

module "callback_lambda" {
  source = "./modules/function"

  project        = local.project
  env            = local.env
  region         = local.region
  account_id     = local.account_id
  name           = "${local.project_tag}-api-callback"
  description    = "Callback for the Spotify API as part of the authorization flow"
  iam_role_arn   = module.permissions.write_token_secret_role_arn
  ecr_image_name = "${local.project_tag}-api-callback:latest"
  env_variables = {
    REGION            = local.region
    STATIC_SECRETS_ID = module.secrets.static_arn
    TOKEN_SECRET_ID   = module.secrets.access_token_arn
    REDIRECT_URI      = local.redirect_uri
    BASE_URI          = local.base_uri
  }
}

module "login_lambda" {
  source = "./modules/function"

  project        = local.project
  env            = local.env
  region         = local.region
  account_id     = local.account_id
  name           = "${local.project_tag}-api-login"
  description    = "Login to Spotify as part of the authorization flow of the Spotify API"
  iam_role_arn   = module.permissions.read_static_secrets_role_arn
  ecr_image_name = "${local.project_tag}-api-login:latest"
  env_variables = {
    REGION            = local.region
    STATIC_SECRETS_ID = module.secrets.static_arn
    REDIRECT_URI      = local.redirect_uri
  }
}

module "get_queue_lambda" {
  source = "./modules/function"

  project        = local.project
  env            = local.env
  region         = local.region
  account_id     = local.account_id
  name           = "${local.project_tag}-get-queue"
  description    = "Lambda Function to retrieve live queue of the player from the Spotify API"
  iam_role_arn   = module.permissions.read_token_secret_role_arn
  ecr_image_name = "${local.project_tag}-get-queue:latest"
  env_variables = {
    REGION          = local.region
    TOKEN_SECRET_ID = module.secrets.access_token_arn
  }
}

module "update_queue_lambda" {
  source = "./modules/function"

  project        = local.project
  env            = local.env
  region         = local.region
  account_id     = local.account_id
  name           = "${local.project_tag}-update-queue"
  description    = "Lambda Function to add song to Spotify queue with the Spotify API"
  iam_role_arn   = module.permissions.read_token_secret_role_arn
  ecr_image_name = "${local.project_tag}-update-queue:latest"
  env_variables = {
    REGION          = local.region
    TOKEN_SECRET_ID = module.secrets.access_token_arn
  }
}

module "search_track_lambda" {
  source = "./modules/function"

  project        = local.project
  env            = local.env
  region         = local.region
  account_id     = local.account_id
  name           = "${local.project_tag}-search-track"
  description    = "Lambda Function to query for a Spotify song in the Spotify API"
  iam_role_arn   = module.permissions.read_token_secret_role_arn
  ecr_image_name = "${local.project_tag}-search-track:latest"
  env_variables = {
    REGION          = local.region
    TOKEN_SECRET_ID = module.secrets.access_token_arn
  }
}

module "api" {
  source = "./modules/api"

  project                    = local.project
  env                        = local.env
  project_tag                = local.project_tag
  region                     = local.region
  account_id                 = local.account_id
  api_domain_name_id         = module.domain.api_domain_name_id
  callback_lambda_invoke_arn = module.callback_lambda.invoke_arn
  callback_lambda_arn        = module.callback_lambda.arn
  lambda_routes = [
    {
      lambda_invoke_arn = module.callback_lambda.invoke_arn
      lambda_arn        = module.callback_lambda.arn
      method            = "GET"
      route             = local.callback_route
    },
    {
      lambda_invoke_arn = module.login_lambda.invoke_arn
      lambda_arn        = module.login_lambda.arn
      method            = "GET"
      route             = "/login"
    },
    {
      lambda_invoke_arn = module.get_queue_lambda.invoke_arn
      lambda_arn        = module.get_queue_lambda.arn
      method            = "GET"
      route             = "/queue"
    },
    {
      lambda_invoke_arn = module.update_queue_lambda.invoke_arn
      lambda_arn        = module.update_queue_lambda.arn
      method            = "POST"
      route             = "/song"
    },
    {
      lambda_invoke_arn = module.search_track_lambda.invoke_arn
      lambda_arn        = module.search_track_lambda.arn
      method            = "GET"
      route             = "/search"
    },
  ]
}
