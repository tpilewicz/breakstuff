# OBJECTS
resource "aws_s3_bucket_object" "index" {
  bucket = aws_s3_bucket.subdomain.bucket
  key    = local.index_key
  content = local.rendered_index
  # etag makes the file update when it changes
  etag   = md5(local.rendered_index)
  content_type = "text/html"
}

resource "aws_s3_bucket_object" "css" {
  bucket = aws_s3_bucket.subdomain.bucket
  key    = local.css_file
  source = "../../components/show/assets/${local.css_file}"
  # etag makes the file update when it changes
  etag   = filemd5("../../components/show/assets/${local.css_file}")
  content_type = "text/css"
}

resource "aws_s3_bucket_object" "ok" {
  bucket = aws_s3_bucket.subdomain.bucket
  key    = local.ok_file
  source = "../../components/show/assets/${local.ok_file}"
  # etag makes the file update when it changes
  etag   = filemd5("../../components/show/assets/${local.ok_file}")
  content_type = "image/png"
}

resource "aws_s3_bucket_object" "broken" {
  bucket = aws_s3_bucket.subdomain.bucket
  key    = local.broken_file
  source = "../../components/show/assets/${local.broken_file}"
  # etag makes the file update when it changes
  etag   = filemd5("../../components/show/assets/${local.broken_file}")
  content_type = "image/png"
}

# BUCKETS
resource "aws_s3_bucket" "main_domain" {
  bucket = var.domain_name

  website {
    redirect_all_requests_to = local.subdomain_name
  }
}

resource "aws_s3_bucket" "subdomain" {
  bucket = local.subdomain_name

  website {
    index_document = "index.html"
    # error_document = "error.html"
  }
}

resource "aws_s3_bucket_policy" "subdomain" {
  bucket = aws_s3_bucket.subdomain.id
  policy = data.aws_iam_policy_document.subdomain.json
}

data "aws_iam_policy_document" "subdomain" {
  version = "2012-10-17"
  statement {
    effect =  "Allow"
    principals {
        type = "*"
        identifiers = ["*"]
    }
    actions = ["s3:GetObject"]
    resources = ["${aws_s3_bucket.subdomain.arn}/*"]
  }
}

# ROUTE53

resource "aws_route53_zone" "main" {
  name = var.domain_name
}

resource "aws_route53_record" "main_domain" {
  zone_id = "${aws_route53_zone.main.zone_id}"
  name = var.domain_name
  type = "A"

  alias {
    name = aws_s3_bucket.main_domain.website_domain
    zone_id = aws_s3_bucket.main_domain.hosted_zone_id
    evaluate_target_health = false
  }
}

resource "aws_route53_record" "subdomain" {
  zone_id = "${aws_route53_zone.main.zone_id}"
  name = local.subdomain_name
  type = "A"

  alias {
    name = aws_s3_bucket.subdomain.website_domain
    zone_id = aws_s3_bucket.subdomain.hosted_zone_id
    evaluate_target_health = false
  }
}
