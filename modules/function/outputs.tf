output "invoke_arn" {
  description = "ARN to invoke the Lambda Function"
  value       = length(var.vpc_config) == 0 ? aws_lambda_function.non_vpc.0.invoke_arn : aws_lambda_function.vpc.0.invoke_arn
}

output "arn" {
  description = "ARN of the Lambda Function"
  value       = length(var.vpc_config) == 0 ? aws_lambda_function.non_vpc.0.arn : aws_lambda_function.vpc.0.arn
}
