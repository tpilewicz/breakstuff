data "archive_file" "troublemaker" {
  source_file = local.troublemaker_src
  output_path = "${local.troublemaker}.zip"
  type        = "zip"
}

resource "aws_lambda_function" "troublemaker" {
  function_name    = "${local.default_name}-troublemaker"
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
      NB_ROWS     = local.nb_rows
      NB_COLS     = local.nb_cols

      SRC_HASH = filebase64sha256(local.troublemaker)
    }
  }
}
