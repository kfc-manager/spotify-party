# gets triggered on each apply
resource "null_resource" "always_run" {
  triggers = {
    timestamp = "${timestamp()}"
  }
}

resource "aws_lambda_function" "non_vpc" {
  count         = length(var.vpc_config) == 0 ? 1 : 0 # this lambda resource is not added to vpc
  function_name = var.name
  package_type  = "Image"
  role          = var.iam_role_arn
  image_uri     = "${var.account_id}.dkr.ecr.${var.region}.amazonaws.com/${var.ecr_image_name}"
  memory_size   = var.memory_size
  timeout       = var.timeout

  environment {
    variables = var.env_variables
  }

  # always recreates the function so it pulls the latest image
  lifecycle {
    replace_triggered_by = [
      null_resource.always_run
    ]
  }

  tags = {
    Project     = var.project
    Environment = var.env
    Description = var.description
  }
}

resource "aws_lambda_function" "vpc" {
  count         = length(var.vpc_config) > 0 ? 1 : 0 # this lambda resource is added to vpc
  function_name = var.name
  package_type  = "Image"
  role          = var.iam_role_arn
  image_uri     = "${var.account_id}.dkr.ecr.${var.region}.amazonaws.com/${var.ecr_image_name}"
  memory_size   = var.memory_size
  timeout       = var.timeout

  environment {
    variables = var.env_variables
  }

  vpc_config {
    subnet_ids         = var.vpc_config["subnet_ids"]
    security_group_ids = var.vpc_config["security_group_ids"]
  }

  # always recreates the function so it pulls the latest image
  lifecycle {
    replace_triggered_by = [
      null_resource.always_run
    ]
  }

  tags = {
    Project     = var.project
    Environment = var.env
    Description = var.description
  }
}
