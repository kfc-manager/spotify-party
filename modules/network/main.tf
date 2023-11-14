locals {
  public_subnet_cidr = "10.0.1.0/24"
}

resource "aws_vpc" "main" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Project     = var.project
    Environment = var.env
    Description = "VPC for Lambda Function"
  }
}

resource "aws_subnet" "public" {
  cidr_block        = local.public_subnet_cidr
  vpc_id            = aws_vpc.main.id
  availability_zone = var.availability_zone

  tags = {
    Project     = var.project
    Environment = var.env
    Description = "Public subnet for the NAT Gateway"
  }
}

resource "aws_internet_gateway" "main" {
  vpc_id = aws_vpc.main.id

  tags = {
    Project     = var.project
    Environment = var.env
    Description = "Internet Gateway for public subnet internet access"
  }
}

resource "aws_security_group" "main" {
  name        = "${var.project_tag}-lambda-function"
  description = "Security group to allow Lambda Function outbound traffic"
  vpc_id      = aws_vpc.main.id

  tags = {
    Project     = var.project
    Environment = var.env
    Description = "Security group to allow Lambda Function outbound traffic"
  }
}

resource "aws_vpc_security_group_egress_rule" "main" {
  security_group_id = aws_security_group.main.id
  from_port         = 443
  to_port           = 443
  ip_protocol       = "tcp"
  cidr_ipv4         = local.public_subnet_cidr
}

resource "aws_eip" "main" {
  domain = "vpc"
}

resource "aws_nat_gateway" "main" {
  allocation_id = aws_eip.main.id
  subnet_id     = aws_subnet.public.id

  depends_on = [aws_internet_gateway.main]

  tags = {
    Project     = var.project
    Environment = var.env
    Description = "NAT Gateway for outbound traffic of Lambda Function from the private subnet"
  }
}

resource "aws_route_table" "public" {
  vpc_id = aws_vpc.main.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.main.id
  }

  tags = {
    Project     = var.project
    Environment = var.env
    Description = "Route table to route public subnet traffic to Internet Gateway"
  }
}

resource "aws_route_table_association" "public" {
  subnet_id      = aws_subnet.public.id
  route_table_id = aws_route_table.public.id
}

resource "aws_subnet" "private" {
  cidr_block        = "10.0.2.0/24"
  vpc_id            = aws_vpc.main.id
  availability_zone = var.availability_zone

  tags = {
    Project     = var.project
    Environment = var.env
    Description = "Private subnet to isolate the Lambda Function form inbound traffic"
  }
}

resource "aws_route_table" "private" {
  vpc_id = aws_vpc.main.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_nat_gateway.main.id
  }

  tags = {
    Project     = var.project
    Environment = var.env
    Description = "Route table to route public subnet traffic to Internet Gateway"
  }
}

resource "aws_route_table_association" "private" {
  subnet_id      = aws_subnet.private.id
  route_table_id = aws_route_table.private.id
}
