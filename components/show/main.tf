# OBJECTS
resource "aws_s3_bucket_object" "index" {
  bucket = aws_s3_bucket.main.bucket
  key    = local.index_key
  content = local.rendered_index
  # etag makes the file update when it changes
  etag   = md5(local.rendered_index)
  content_type = "text/html"
}

resource "aws_s3_bucket_object" "css" {
  bucket = aws_s3_bucket.main.bucket
  key    = local.css_file
  source = "../../components/show/assets/${local.css_file}"
  # etag makes the file update when it changes
  etag   = filemd5("../../components/show/assets/${local.css_file}")
  content_type = "text/css"
}

resource "aws_s3_bucket_object" "ok" {
  bucket = aws_s3_bucket.main.bucket
  key    = local.ok_file
  source = "../../components/show/assets/${local.ok_file}"
  # etag makes the file update when it changes
  etag   = filemd5("../../components/show/assets/${local.ok_file}")
  content_type = "image/png"
}

resource "aws_s3_bucket_object" "broken" {
  bucket = aws_s3_bucket.main.bucket
  key    = local.broken_file
  source = "../../components/show/assets/${local.broken_file}"
  # etag makes the file update when it changes
  etag   = filemd5("../../components/show/assets/${local.broken_file}")
  content_type = "image/png"
}

# BUCKET
resource "aws_s3_bucket" "main" {
  bucket = "tpilewicz-${local.default_name}-public"

  website {
    index_document = "index.html"
    # error_document = "error.html"
  }
}

resource "aws_s3_bucket_policy" "main" {
  bucket = aws_s3_bucket.main.id
  policy = data.aws_iam_policy_document.main.json
}

data "aws_iam_policy_document" "main" {
  version = "2012-10-17"
  statement {
    effect =  "Allow"
    principals {
        type = "*"
        identifiers = ["*"]
    }
    actions = ["s3:GetObject"]
    resources = ["${aws_s3_bucket.main.arn}/*"]
  }
}
