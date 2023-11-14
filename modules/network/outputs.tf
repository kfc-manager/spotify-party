output "security_group_id" {
  description = "ID of the security group assigned to Lambda Function to access NAT Gateway"
  value       = aws_security_group.main.id
}

output "private_subnet_id" {
  description = "ID of the private subnet where Lambda Function lives in"
  value       = aws_subnet.private.id
}
