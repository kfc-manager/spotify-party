output "api_domain_name_id" {
  description = "ID of the API domain name"
  value       = aws_apigatewayv2_domain_name.main.id
}

output "acm_certificate_arn" {
  description = "ARN of ACM certificate for base domain"
  value       = aws_acm_certificate.main.arn
}

output "public_host_zone_id" {
  description = "ID of public DNS host zone"
  value       = data.aws_route53_zone.public.id
}
