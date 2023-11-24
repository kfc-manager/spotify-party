resource "aws_apigatewayv2_api" "main" {
  name          = "${var.project_tag}-api"
  description   = "${var.project} API"
  protocol_type = "HTTP"

  cors_configuration {
    allow_methods = ["GET", "POST"]
    allow_origins = ["https://kiliansqueue.com"]
    allow_headers = ["Content-Type", "Authorization"]
  }

  tags = {
    Project     = var.project
    Environment = var.env
    Description = "API Gateway that provides endpoints"
  }
}

resource "aws_apigatewayv2_integration" "main" {
  count              = length(var.lambda_routes)
  api_id             = aws_apigatewayv2_api.main.id
  integration_uri    = var.lambda_routes[count.index]["lambda_invoke_arn"]
  integration_type   = "AWS_PROXY"
  integration_method = "POST"
}

resource "aws_apigatewayv2_route" "main" {
  count     = length(var.lambda_routes)
  api_id    = aws_apigatewayv2_api.main.id
  route_key = "${var.lambda_routes[count.index]["method"]} ${var.lambda_routes[count.index]["route"]}"
  target    = "integrations/${aws_apigatewayv2_integration.main[count.index].id}"
}

resource "aws_lambda_permission" "main" {
  count         = length(var.lambda_routes)
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = var.lambda_routes[count.index]["lambda_arn"]
  principal     = "apigateway.amazonaws.com"
  source_arn    = "arn:aws:execute-api:${var.region}:${var.account_id}:${aws_apigatewayv2_api.main.id}/*/*${var.lambda_routes[count.index]["route"]}"
}

resource "aws_apigatewayv2_stage" "main" {
  api_id      = aws_apigatewayv2_api.main.id
  name        = "$default"
  auto_deploy = true

  dynamic "route_settings" {
    for_each = var.lambda_routes

    content {
    route_key              = "${route_settings.value["method"]} ${route_settings.value["route"]}"
    throttling_burst_limit = 1000
    throttling_rate_limit  = 0.01
    }
  }
}

resource "aws_apigatewayv2_api_mapping" "main" {
  api_id      = aws_apigatewayv2_api.main.id
  domain_name = var.api_domain_name_id
  stage       = aws_apigatewayv2_stage.main.id
}
