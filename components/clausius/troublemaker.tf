data "archive_file" "troublemaker" {
  source_file = local.troublemaker_src
  output_path = "${local.troublemaker_src}.zip"
  type        = "zip"
}

resource "aws_lambda_function" "troublemaker" {
  function_name    = local.troublemaker_name
  filename         = data.archive_file.troublemaker.output_path
  handler          = "main"
  source_code_hash = data.archive_file.troublemaker.output_base64sha256
  role             = aws_iam_role.troublemaker.arn
  runtime          = "go1.x"
  memory_size      = 128
  timeout          = 6

  depends_on = [aws_cloudwatch_log_group.troublemaker]

  tags = local.tags

  environment {
    variables = {
      ENVIRONMENT = var.environment
      FUNES_TABLE = var.funes_table.name
      NB_ROWS     = local.nb_rows
      NB_COLS     = local.nb_cols

      SRC_HASH = filebase64sha256(local.troublemaker_src)
    }
  }
}

resource "aws_iam_role" "troublemaker" {
  name = local.troublemaker_name

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow"
    }
  ]
}
EOF

  tags = local.tags
}

resource "aws_cloudwatch_log_group" "troublemaker" {
  name              = "/aws/lambda/${local.troublemaker_name}"
  retention_in_days = "14"
  tags              = local.tags
}

resource "aws_iam_role_policy_attachment" "troublemaker_execution" {
  role       = aws_iam_role.troublemaker.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

## Funes access

data "aws_iam_policy_document" "hit_funes" {
  statement {
    sid = "HitFunes"

    actions = [
      "dynamodb:GetItem",
      "dynamodb:UpdateItem"
    ]

    resources = [var.funes_table.arn]
  }
}

resource "aws_iam_role_policy" "troublemaker_funes" {
  name       = "TroubleMakerHitFunes"
  role       = aws_iam_role.troublemaker.name
  policy     = data.aws_iam_policy_document.hit_funes.json
}

## Trigger every week

resource "aws_cloudwatch_event_rule" "troublemaker" {
  name                = "${local.troublemaker_name}-one-week"
  schedule_expression = "rate(7 days)"
  description         = "Fires every week"

  tags = local.tags
}

resource "aws_cloudwatch_event_target" "troublemaker" {
  rule      = aws_cloudwatch_event_rule.troublemaker.name
  target_id = "${local.troublemaker_name}-target"
  arn       = aws_lambda_function.troublemaker.arn
}

resource "aws_lambda_permission" "permission" {
  statement_id  = "AllowExecutionFromCloudWatch"
  action        = "lambda:InvokeFunction"
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.troublemaker.arn
  function_name = aws_lambda_function.troublemaker.arn
}
