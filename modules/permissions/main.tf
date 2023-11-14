data "aws_iam_policy_document" "assume_role" {
  statement {
    effect = "Allow"

    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }

    actions = ["sts:AssumeRole"]
  }
}

data "aws_iam_policy_document" "read_static_secrets" {
  statement {
    actions = ["secretsmanager:GetSecretValue"]

    resources = [var.static_secrets_arn]
  }
}

data "aws_iam_policy_document" "read_token_secret" {
  statement {
    actions = ["secretsmanager:GetSecretValue"]

    resources = [var.token_secret_arn]
  }
}

data "aws_iam_policy_document" "write_token_secret" {
  statement {
    actions = ["secretsmanager:UpdateSecret"]

    resources = [var.token_secret_arn]
  }
}

data "aws_iam_policy_document" "vpc_assignment" {
  statement {
    effect = "Allow"

    actions = [
      "ec2:DescribeInstances",
      "ec2:CreateNetworkInterface",
      "ec2:AttachNetworkInterface",
      "ec2:DescribeNetworkInterfaces",
      "autoscaling:CompleteLifecycleAction",
      "ec2:DeleteNetworkInterface"
    ]

    resources = ["*"]
  }
}

resource "aws_iam_role" "read_static_secrets" {
  name               = "${var.project_tag}-lambda-read-static-secrets"
  assume_role_policy = data.aws_iam_policy_document.assume_role.json

  inline_policy {
    name   = "${var.project_tag}-read-static-secrets"
    policy = data.aws_iam_policy_document.read_static_secrets.json
  }

  inline_policy {
    name   = "${var.project_tag}-vpc_assignment"
    policy = data.aws_iam_policy_document.vpc_assignment.json
  }

  tags = {
    Project     = var.project
    Environment = var.env
    Description = "IAM Role for a Lambda Function to read the static secrets from the Secrets Manager"
  }
}

resource "aws_iam_role" "read_token_secret" {
  name               = "${var.project_tag}-lambda-read-token-secret"
  assume_role_policy = data.aws_iam_policy_document.assume_role.json

  inline_policy {
    name   = "${var.project_tag}-read-token-secret"
    policy = data.aws_iam_policy_document.read_token_secret.json
  }

  tags = {
    Project     = var.project
    Environment = var.env
    Description = "IAM Role for a Lambda Function to read the token secret from the Secrets Manager"
  }
}

resource "aws_iam_role" "write_token_secret" {
  name               = "${var.project_tag}-lambda-write-token-secret"
  assume_role_policy = data.aws_iam_policy_document.assume_role.json

  inline_policy {
    name   = "${var.project_tag}-read-static-secrets"
    policy = data.aws_iam_policy_document.read_static_secrets.json
  }

  inline_policy {
    name   = "${var.project_tag}-write-token-secret"
    policy = data.aws_iam_policy_document.write_token_secret.json
  }

  tags = {
    Project     = var.project
    Environment = var.env
    Description = "IAM Role for a Lambda Function to read the static secrets and write the token secret from the Secrets Manager"
  }
}
