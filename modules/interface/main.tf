resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

locals {
  bucket_name = "${var.project_tag}-user-interface-${random_string.suffix.result}"
}

resource "aws_s3_bucket" "main" {
  bucket = local.bucket_name

  tags = {
    Project     = var.project
    Environment = var.env
    Description = "S3 bucket to serve static website files"
  }
}

resource "aws_s3_bucket_ownership_controls" "main" {
  bucket = aws_s3_bucket.main.id

  rule {
    object_ownership = "BucketOwnerPreferred"
  }
}


resource "aws_s3_bucket_public_access_block" "main" {
  bucket = aws_s3_bucket.main.id

  block_public_acls       = false
  block_public_policy     = false
  ignore_public_acls      = false
  restrict_public_buckets = false
}


resource "aws_s3_bucket_acl" "main" {
  bucket = aws_s3_bucket.main.id
  acl    = "public-read"

  depends_on = [aws_s3_bucket_ownership_controls.main, aws_s3_bucket_public_access_block.main]
}

resource "aws_cloudfront_distribution" "main" {
  enabled             = true
  default_root_object = "index.html"
  aliases             = [var.domain_name]

  origin {
    origin_id   = "${local.bucket_name}-origin"
    domain_name = aws_s3_bucket.main.bucket_regional_domain_name
    custom_origin_config {
      http_port              = 80
      https_port             = 443
      origin_protocol_policy = "http-only"
      origin_ssl_protocols   = ["TLSv1"]
    }
  }

  default_cache_behavior {

    target_origin_id = "${local.bucket_name}-origin"
    allowed_methods  = ["GET", "HEAD"]
    cached_methods   = ["GET", "HEAD"]

    forwarded_values {
      query_string = true

      cookies {
        forward = "all"
      }
    }

    viewer_protocol_policy = "redirect-to-https"
    min_ttl                = 0
    default_ttl            = 0
    max_ttl                = 0
  }

  restrictions {
    geo_restriction {
      restriction_type = "none"
    }
  }

  viewer_certificate {
    acm_certificate_arn = var.acm_certificate_arn
    ssl_support_method  = "sni-only"
  }

  price_class = "PriceClass_200"

  tags = {
    Project     = var.project
    Environment = var.env
    Description = "CloudFornt for static website content"
  }
}

resource "aws_route53_record" "main" {
  name    = var.domain_name
  type    = "A"
  zone_id = var.public_host_zone_id

  alias {
    name                   = aws_cloudfront_distribution.main.domain_name
    zone_id                = aws_cloudfront_distribution.main.hosted_zone_id
    evaluate_target_health = false
  }
}

data "aws_iam_policy_document" "s3_access" {
  statement {
    sid       = "AllowCloudFrontServicePrincipal"
    actions   = ["s3:GetObject"]
    resources = ["${aws_s3_bucket.main.arn}/*"]

    condition {
      test     = "StringEquals"
      variable = "AWS:SourceArn"
      values   = [aws_cloudfront_distribution.main.arn]
    }

    principals {
      type        = "Service"
      identifiers = ["cloudfront.amazonaws.com"]
    }
  }
}

data "aws_iam_policy_document" "s3_put" {
  statement {
    actions   = ["s3:PutObject"]
    resources = ["${aws_s3_bucket.main.arn}/*"]

    principals {
      type        = "AWS"
      identifiers = [var.account_arn]
    }
  }
}

resource "aws_s3_bucket_policy" "s3_access" {
  bucket = aws_s3_bucket.main.id
  policy = data.aws_iam_policy_document.s3_access.json
}

resource "aws_s3_bucket_policy" "s3_put" {
  bucket = aws_s3_bucket.main.id
  policy = data.aws_iam_policy_document.s3_put.json
}

# gets triggered on each apply
resource "null_resource" "always_run" {
  triggers = {
    timestamp = "${timestamp()}"
  }
}

resource "aws_s3_object" "index" {
  bucket       = aws_s3_bucket.main.id
  key          = "index.html"
  source       = "${abspath(path.root)}/frontend/dist/index.html"
  content_type = "text/html"
  acl          = "public-read"
  etag         = filemd5("${abspath(path.root)}/frontend/dist/index.html")

  # always recreates the function so it pulls the latest image
  lifecycle {
    replace_triggered_by = [
      null_resource.always_run
    ]
  }
}

resource "aws_s3_object" "css" {
  for_each = fileset("${abspath(path.root)}/frontend/dist/", "**/*.css")

  bucket       = aws_s3_bucket.main.id
  key          = each.value
  source       = "${abspath(path.root)}/frontend/dist/${each.value}"
  content_type = "text/css"
  acl          = "public-read"
  etag         = filemd5("${abspath(path.root)}/frontend/dist/${each.value}")

  # always recreates the function so it pulls the latest image
  lifecycle {
    replace_triggered_by = [
      null_resource.always_run
    ]
  }
}

resource "aws_s3_object" "js" {
  for_each = fileset("${abspath(path.root)}/frontend/dist/", "**/*.js")

  bucket       = aws_s3_bucket.main.id
  key          = each.value
  source       = "${abspath(path.root)}/frontend/dist/${each.value}"
  content_type = "application/javascript"
  acl          = "public-read"
  etag         = filemd5("${abspath(path.root)}/frontend/dist/${each.value}")

  # always recreates the function so it pulls the latest image
  lifecycle {
    replace_triggered_by = [
      null_resource.always_run
    ]
  }
}
