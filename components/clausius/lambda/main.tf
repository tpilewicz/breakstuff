data "aws_region" "current" {}
data "aws_caller_identity" "current" {}

## Lambda, its source, its role, its log group

data "archive_file" "this" {
  source_file = var.lambda_src
  output_path = "${var.lambda_src}.zip"
  type        = "zip"
}

resource "aws_lambda_function" "this" {
  function_name    = local.default_name
  filename         = data.archive_file.this.output_path
  handler          = "main"
  source_code_hash = data.archive_file.this.output_base64sha256
  role             = aws_iam_role.this.arn
  runtime          = "go1.x"
  memory_size      = 128
  timeout          = 6

  depends_on = [aws_cloudwatch_log_group.this]

  tags = local.tags

  environment {
    variables = {
      ENVIRONMENT = var.environment
      FUNES_TABLE = var.funes_table.name
      NB_ROWS     = var.nb_rows
      NB_COLS     = var.nb_cols

      SRC_HASH = filebase64sha256(var.lambda_src)
    }
  }
}

resource "aws_iam_role" "this" {
  name = local.default_name

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

resource "aws_cloudwatch_log_group" "this" {
  name              = "/aws/lambda/${local.default_name}"
  retention_in_days = "14"
  tags              = local.tags
}

resource "aws_iam_role_policy_attachment" "attach_execution_policy" {
  role       = aws_iam_role.this.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

## API

resource "aws_api_gateway_resource" "this" {
  rest_api_id = var.api.id
  parent_id   = var.api.root_resource_id
  path_part   = var.api_resource
}

resource "aws_api_gateway_method" "this" {
  rest_api_id      = var.api.id
  resource_id      = aws_api_gateway_resource.this.id
  http_method      = var.http_method
  api_key_required = false
  authorization    = "NONE"
}

resource "aws_api_gateway_integration" "this" {
  rest_api_id             = var.api.id
  resource_id             = aws_api_gateway_resource.this.id
  http_method             = aws_api_gateway_method.this.http_method
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.this.invoke_arn
}

resource "aws_lambda_permission" "apigw" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.this.function_name
  principal     = "apigateway.amazonaws.com"

  # More: http://docs.aws.amazon.com/apigateway/latest/developerguide/api-gateway-control-access-using-iam-policies-to-invoke-api.html
  source_arn = "arn:aws:execute-api:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:${var.api.id}/*/${aws_api_gateway_method.this.http_method}${aws_api_gateway_resource.this.path}"
}
