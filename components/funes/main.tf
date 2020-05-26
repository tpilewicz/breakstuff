locals {
  component    = "funes"
  default_name = "${var.environment}-${local.component}"
}

variable environment {}
variable table_name {}

resource "aws_dynamodb_table" "main" {
  name           = var.table_name
  billing_mode   = "PAY_PER_REQUEST"
  hash_key       = "K"

  attribute {
    name = "K"
    type = "S"
  }

  tags = {
    Name        = local.default_name
    Environment = var.environment
  }
}

output "table" {
    value = aws_dynamodb_table.main
}
